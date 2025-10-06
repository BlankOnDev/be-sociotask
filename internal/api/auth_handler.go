package api

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/harundarat/be-socialtask/internal/auth"
	gAuth "github.com/harundarat/be-socialtask/internal/auth/google"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
	"golang.org/x/oauth2"
)

type AndroidLoginRequest struct {
	TokenID string `json:"token_id"`
}

type AuthHandler struct {
	logger      *log.Logger
	userStore   store.UserStore
	oauthConf   *oauth2.Config
	oauthGoogle *oauth2.Config
}

func NewAuthHandler(logger *log.Logger, userStore store.UserStore, oauthGoogle, oauthConf *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		userStore:   userStore,
		logger:      logger,
		oauthConf:   oauthConf,
		oauthGoogle: oauthGoogle,
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

func (h *AuthHandler) CallbackAuthenticationGooogle(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		h.logger.Println("error invalid oauth google state")
		http.Redirect(w, r, "/failed?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.oauthGoogle.Exchange(r.Context(), code)
	if err != nil {
		h.logger.Printf("error, failed exchange code : %v", err)
		http.Redirect(w, r, "/failed?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// get client dari Google akun
	client := h.oauthGoogle.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		h.logger.Printf("error, failed get user info data : %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, []string{"failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		h.logger.Printf("error, unmarshal data : %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, []string{"failed to parse user info"})
		return
	}

	user, err := h.userStore.FindEmailForGoogle(userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		h.logger.Printf("error method FindEmailForGoogle, failed get find email : %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, []string{"database operation failed"})
		return
	}

	jwtToken, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, "thisissecret")
	if err != nil {
		h.logger.Printf("ERROR: generating token: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, []string{"failed to generate token"})
		return
	}
	// redirectURL := fmt.Sprintf("/success?token=%s", jwtToken) // cek lokal jangan dipush!!!
	// http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageOAuthSuccess, http.StatusOK, utils.Envelope{"token": jwtToken, "user_id": user.ID}, nil)
}

func (h *AuthHandler) LoginAuthenticationGooogle(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32) // esih, rung paham
	_, err := rand.Read(b)
	if err != nil {
		h.logger.Printf("Error Failed to generate state: %v", err)
		http.Redirect(w, r, "/failed?error=state_generation_failed", http.StatusTemporaryRedirect)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Path:     "/",
	})

	url := h.oauthGoogle.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGoogleLoginAndroid(w http.ResponseWriter, r *http.Request) {
	var req AndroidLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, []string{"request not valid"})
		return
	}

	payload, err := gAuth.GoogleVerifytokenID(req.TokenID)
	if err != nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageUnauthorized, http.StatusUnauthorized, nil, []string{"token signature is unknown"})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	googleUserID := payload.Claims["sub"].(string)

	user, err := h.userStore.FindEmailForGoogle(googleUserID, email, name)
	if err != nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, []string{"database operation failed"})
		return
	}

	jwtToken, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, "thisissecret")
	if err != nil {
		h.logger.Printf("ERROR: generating token: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, []string{"failed to generate token"})
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageOAuthSuccess, http.StatusOK, utils.Envelope{"token": jwtToken, "user_id": user.ID}, nil)
}
