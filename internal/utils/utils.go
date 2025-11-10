package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any
type Status string
type Message string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)
const (
	MessageLoginSuccess       Message = "login successful"
	MessageLoginFailed        Message = "login failed"
	MessageRegisterSuccess    Message = "registration successful"
	MessageRegisterFailed     Message = "registration failed"
	MessageTaskCreated        Message = "task created successfully"
	MessageTaskRetrieved      Message = "task retrieved successfully"
	MessageTasksFetched       Message = "tasks fetched successfully"
	MessageTasksUpdated       Message = "tasks updated successfully"
	MessageTasksDeleted       Message = "tasks deleted successfully"
	MessageActionInvalidType  Message = "invalid action type"
	MessageActionCreated      Message = "action created successfully"
	MessageActionRetrieved    Message = "action retrieved successfully"
	MessageActionsFetched     Message = "actions fetched successfully"
	MessageActionsDelete      Message = "actions deleted successfully"
	MessageActionsUpdated     Message = "actions updated successfully"
	MessageRewardCreated      Message = "reward created successfully"
	MessageRewardRetrieved    Message = "reward retrieved successfully"
	MessageRewardsFetched     Message = "rewards fetched successfully"
	MessageRewardsDelete      Message = "rewards deleted successfully"
	MessageRewardsUpdated     Message = "rewards updated successfully"
	MessageInvalidRequest     Message = "invalid request"
	MessageInternalError      Message = "internal server error"
	MessageUnauthorized       Message = "unauthorized access"
	MessageNotFound           Message = "resource not found"
	MessageInvalidCredentials Message = "invalid credentials"
	MessageTokenGenerated     Message = "token generated successfully"
	MessageValidationFailed   Message = "validation failed"
	MessageOAuthFailed        Message = "oauth authentication failed"
	MessageOAuthSuccess       Message = "oauth authentication successful"
	MessageBadRequest         Message = "bad request"
	MessageUserRetrieved      Message = "user retrieved successfully"
)

func WriteJSON(w http.ResponseWriter, status Status, message Message, statusCode int, data Envelope, errorsList []string) error {
	var response Envelope

	switch status {
	case StatusSuccess:
		response = Envelope{
			"status":  status,
			"message": message,
			"data":    data,
		}
	case StatusError:
		response = Envelope{
			"status":  status,
			"message": message,
			"errors":  errorsList,
		}
	}

	js, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id paramater type")
	}

	return id, nil
}

func ReadPageParam(r *http.Request) int64 {
	pageParam := chi.URLParam(r, "page")
	if pageParam == "" {
		return 1
	}

	page, _ := strconv.ParseInt(pageParam, 10, 64)

	return page
}

// GenerateRandomString creates a cryptographically secure random string of a given length.
func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// This panic is safe because rand.Int should not fail in a normal environment
			panic(err)
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func GenerateSecureRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("add environment variable first!")
	}
	return value
}
