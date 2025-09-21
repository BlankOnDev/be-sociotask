package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/harundarat/be-socialtask/internal/auth"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	logger    *log.Logger
	userStore store.UserStore
	oauthConf *oauth2.Config
}

func NewAuthHandler(logger *log.Logger, userStore store.UserStore, oauthConf *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
		logger:    logger,
		oauthConf: oauthConf,
	}
}

func (h *AuthHandler) HandleTwitterLogin(w http.ResponseWriter, r *http.Request) {
	// generate random state string
	state := utils.GenerateRandomString(32)

	// generate pkce verifier
	verifier := oauth2.GenerateVerifier()

	// store state in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth2_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth2_verifier",
		Value:    verifier,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	url := h.oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleTwitterCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth2_state")
	if err != nil {
		h.logger.Println("ERROR: invalid state parameter")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "missing oauth2 state cookie"})
		return
	}
	if r.URL.Query().Get("state") != stateCookie.Value {
		h.logger.Println("ERROR: invalid state parameter")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid state parameter"})
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Println("ERROR: No authorization code received")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "no authorization code received"})
		return
	}

	h.logger.Printf("Received authorization code: %s", code[:10]+"...") // Log first 10 chars for debugging

	// ambil verifier dari cookie
	verifierCookie, err := r.Cookie("oauth2_verifier")
	if err != nil {
		h.logger.Println("ERROR: Missing oauth2_verifier cookie:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "missing oauth2 verifier cookie"})
		return
	}

	h.logger.Printf("Using verifier: %s", verifierCookie.Value[:10]+"...") // Log first 10 chars

	// exchange authorization code with access token
	token, err := h.oauthConf.Exchange(
		context.Background(),
		code,
		oauth2.VerifierOption(verifierCookie.Value),
	)
	if err != nil {
		h.logger.Println("ERROR: Failed to exchange token:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to exchange token"})
		return
	}

	// use token to get user data from X
	client := h.oauthConf.Client(context.Background(), token)
	response, err := client.Get("https://api.twitter.com/2/users/me?user.fields=id,name,username,profile_image_url")
	if err != nil {
		h.logger.Println("ERROR: Failed to get user from X:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to get user info"})
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		h.logger.Println("ERROR: Failed to read response body:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to read user data"})
		return
	}

	h.logger.Printf("Twitter API response: %s", string(body))

	var twitterUser struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &twitterUser); err != nil {
		h.logger.Println("ERROR: Failed to unmarshal twitter user data:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to process user data"})
		return
	}

	if len(twitterUser.Errors) > 0 {
		h.logger.Printf("ERROR: Twitter API errors: %+v", twitterUser.Errors)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to get user info from Twitter"})
		return
	}

	if twitterUser.Data.ID == "" {
		h.logger.Println("ERROR: No user data received from Twitter")
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "no user data received"})
		return
	}

	// If the user's X account has no email, create a placeholder.
	// database requires a unique email for each user.
	if twitterUser.Data.Email == "" {
		twitterUser.Data.Email = twitterUser.Data.Username + "@twitter.user"
	}

	// check if user already exists
	h.logger.Printf("Checking if user exists with email: %s", twitterUser.Data.Email)
	user, err := h.userStore.GetUserByEmail(twitterUser.Data.Email)
	if err != nil && err != sql.ErrNoRows {
		h.logger.Println("ERROR: checking existing user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "database error"})
		return
	}

	// if user doesn't exists, create new user
	if err == sql.ErrNoRows {
		h.logger.Printf("User not found, creating new user with username: %s, email: %s", twitterUser.Data.Username, twitterUser.Data.Email)
		newUser := &store.User{
			Username: twitterUser.Data.Username,
			Email:    twitterUser.Data.Email,
			Bio:      "Twitter user",
		}

		// user login via Oauth don't have a password in the system.
		// generate a random password to satisfy the database.
		// user will never need to know or use this password.
		randomPassword := utils.GenerateRandomString(16)
		err := newUser.PasswordHash.Set(randomPassword)
		if err != nil {
			h.logger.Println("ERROR: hashing password for new oauth user:", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
			return
		}

		createdUser, err := h.userStore.CreateUser(newUser)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				h.logger.Printf("ERROR: unique constraint violation: %v", err)
				utils.WriteJSON(w, http.StatusConflict, utils.Envelope{"error": "a user with this username or email already exists"})
				return
			}
			h.logger.Printf("ERROR: creating new user: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create user"})
			return
		}
		h.logger.Printf("Successfully created user with ID: %d", createdUser.ID)
		user = createdUser
	} else {
		h.logger.Printf("Found existing user with ID: %d", user.ID)
	}

	// Ensure user is not nil before generating JWT
	if user == nil {
		h.logger.Println("ERROR: user is nil after OAuth process")
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "user creation failed"})
		return
	}

	// Generate a JWT for the user
	jwtToken, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, "thisissecret")
	if err != nil {
		h.logger.Printf("ERROR: generating JWT token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": jwtToken, "user_id": user.ID})
}
