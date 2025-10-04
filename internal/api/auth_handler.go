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
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth2_verifier",
		Value:    verifier,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	url := h.oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleTwitterCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth2_state")
	if err != nil {
		h.logger.Println("ERROR: invalid state parameter")
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusBadRequest, nil, nil)
		return
	}
	if r.URL.Query().Get("state") != stateCookie.Value {
		h.logger.Println("ERROR: invalid state parameter")
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusBadRequest, nil, nil)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Println("ERROR: No authorization code received")
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusBadRequest, nil, nil)
		return
	}

	// ambil verifier dari cookie
	verifierCookie, err := r.Cookie("oauth2_verifier")
	if err != nil {
		h.logger.Println("ERROR: Missing oauth2_verifier cookie:", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusBadRequest, nil, nil)
		return
	}

	// exchange authorization code with access token
	token, err := h.oauthConf.Exchange(
		context.Background(),
		code,
		oauth2.VerifierOption(verifierCookie.Value),
	)
	if err != nil {
		h.logger.Println("ERROR: Failed to exchange token:", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}

	// use token to get user data from X
	client := h.oauthConf.Client(context.Background(), token)
	response, err := client.Get("https://api.twitter.com/2/users/me?user.fields=id,name,username,profile_image_url")
	if err != nil {
		h.logger.Println("ERROR: Failed to get user from X:", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		h.logger.Println("ERROR: Failed to read response body:", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}

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
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}

	if len(twitterUser.Errors) > 0 {
		h.logger.Printf("ERROR: Twitter API errors: %+v", twitterUser.Errors)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}

	if twitterUser.Data.ID == "" {
		h.logger.Println("ERROR: No user data received from Twitter")
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}

	// If the user's X account has no email, create a placeholder.
	// database requires a unique email for each user.
	if twitterUser.Data.Email == "" {
		twitterUser.Data.Email = twitterUser.Data.Username + "@twitter.user"
	}

	// check if user already exists
	user, err := h.userStore.GetUserByEmail(twitterUser.Data.Email)
	if err != nil && err != sql.ErrNoRows {
		h.logger.Println("ERROR: checking existing user:", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
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
			utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
			return
		}

		createdUser, err := h.userStore.CreateUser(newUser)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				h.logger.Printf("ERROR: unique constraint violation: %v", err)
				utils.WriteJSON(w, utils.StatusError, utils.MessageRegisterFailed, http.StatusConflict, nil, nil)
				return
			}
			h.logger.Printf("ERROR: creating new user: %v", err)
			utils.WriteJSON(w, utils.StatusError, utils.MessageRegisterFailed, http.StatusInternalServerError, nil, nil)
			return
		}
		h.logger.Printf("Successfully created user with ID: %d", createdUser.ID)
		user = createdUser
	} else {
	}

	// Ensure user is not nil before generating JWT
	if user == nil {
		h.logger.Println("ERROR: user is nil after OAuth process")
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, nil)
		return
	}

	// Generate a JWT for the user
	jwtToken, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, utils.GetEnv("JWT_SECRET"))
	if err != nil {
		h.logger.Printf("ERROR: generating JWT token: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageOAuthSuccess, http.StatusOK, utils.Envelope{"token": jwtToken, "user_id": user.ID}, nil)
}
