package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/harundarat/be-socialtask/internal/auth"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
)

type UserMiddleware struct {
	userStore store.UserStore
	jwtSecret string
}

func NewUserMiddleware(userStore store.UserStore, jwtSecret string) *UserMiddleware {
	return &UserMiddleware{
		userStore: userStore,
		jwtSecret: jwtSecret,
	}
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request") // bad actor call
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		claims, err := auth.ParseJWTToken(token, um.jwtSecret)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired token"})
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid user ID in token"})
			return
		}

		user, err := um.userStore.GetUserByID(int64(userID))
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "user not found"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "user not found"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
