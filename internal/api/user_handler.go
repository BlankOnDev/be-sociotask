package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/harundarat/be-socialtask/internal/auth"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
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
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
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

	token, err := auth.GenerateJWTToken(user.ID, auth.RoleUser, "thisissecret")
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
