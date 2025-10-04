package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/harundarat/be-socialtask/internal/auth"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
	"golang.org/x/oauth2"
)

type registerUserRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type googleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserHandler struct {
	userStore store.UserStore
	gl        *oauth2.Config
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, gl *oauth2.Config, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		gl:        gl,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Fullname == "" {
		return errors.New("fullname is required")
	}
	if len(req.Fullname) > 255 {
		return errors.New("fullname must be less than 255 characters")
	}
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New(("username must be less than 50 character"))
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (h *UserHandler) validateLoginRequest(req *loginUserRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("error decoding request body: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	err = uh.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageValidationFailed, http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Fullname: sql.NullString{
			String: req.Fullname,
			Valid:  true,
		},
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	user, err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageRegisterFailed, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageRegisterSuccess, http.StatusCreated, utils.Envelope{"user_id": user.ID}, nil)
}

func (uh *UserHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("error decoding request body: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}
	err = uh.validateLoginRequest(&req)
	if err != nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageValidationFailed, http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

	user, err := uh.userStore.GetUserByEmail(req.Email)
	if err != nil {
		uh.logger.Printf("ERROR: get user by email: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}
	if user == nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageNotFound, http.StatusNotFound, nil, nil)
		return
	}

	isMatches, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: matching password: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}
	if !isMatches {
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidCredentials, http.StatusUnauthorized, nil, nil)
		return
	}

	token, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, utils.GetEnv("JWT_SECRET"))
	if err != nil {
		uh.logger.Printf("ERROR: generating token: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageLoginSuccess, http.StatusOK, utils.Envelope{"token": token}, nil)
}

func (uh *UserHandler) HandleGetUserTasks(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("ERROR: reading id param: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	tasks, err := uh.userStore.GetUserTasks(id)
	if err != nil {
		uh.logger.Printf("ERROR: getting user tasks: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTasksFetched, http.StatusOK, utils.Envelope{"tasks": tasks}, nil)
}

func (uh *UserHandler) CallbackAuthenticationGooogle(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		uh.logger.Println("error invalid oauth google state")
		http.Redirect(w, r, "/failed?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := uh.gl.Exchange(r.Context(), code)
	if err != nil {
		uh.logger.Printf("error, failed exchange code : %v", err)
		http.Redirect(w, r, "/failed?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// get client dari Google akun
	client := uh.gl.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		uh.logger.Printf("error, failed get user info data : %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, []string{"failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		uh.logger.Printf("error, unmarshal data : %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageOAuthFailed, http.StatusInternalServerError, nil, []string{"failed to parse user info"})
		return
	}

	user, err := uh.userStore.FindEmailForGoogle(userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		uh.logger.Printf("error method FindEmailForGoogle, failed get find email : %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, []string{"database operation failed"})
		return
	}

	jwtToken, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, utils.GetEnv("JWT_SECRET"))
	if err != nil {
		uh.logger.Printf("ERROR: generating token: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, []string{"failed to generate token"})
		return
	}
	redirectURL := fmt.Sprintf("/success?token=%s", jwtToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (uh *UserHandler) LoginAuthenticationGooogle(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32) // ra paham
	_, err := rand.Read(b)
	if err != nil {
		uh.logger.Printf("Error Failed to generate state: %v", err)
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

	url := uh.gl.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
