package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/oauth"
	"github.com/evermos/boilerplate-go/transport/http/response"
)

type Authentication struct {
	db *infras.MySQLConn
}

const (
	HeaderAuthorization = "Authorization"
)

func ProvideAuthentication(db *infras.MySQLConn) *Authentication {
	return &Authentication{
		db: db,
	}
}

func (a *Authentication) ClientCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		token := oauth.New(a.db.Read, oauth.Config{})

		parseToken, err := token.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) ClientCredentialWithQueryParameter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		token := params.Get("token")
		tokenType := params.Get("token_type")
		accessToken := tokenType + " " + token

		auth := oauth.New(a.db.Read, oauth.Config{})
		parseToken, err := auth.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) Password(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		token := oauth.New(a.db.Read, oauth.Config{})

		parseToken, err := token.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyUserLoggedIn() {
			response.WithMessage(w, http.StatusUnauthorized, oauth.ErrorInvalidPassword)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware to verify the JWT token
func (a *Authentication) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		authHeader := r.Header.Get(HeaderAuthorization)
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse and validate the JWT token
		config := configs.Get()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if token.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.App.JWTAccessKey), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			http.Error(w, "Token invalid", http.StatusUnauthorized)
			return
		}

		// Extract user information from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "No user information from JWT", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "JWT user_id is not found", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "JWT username is not found", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "JWT role is not found", http.StatusUnauthorized)
			return
		}

		// Add user information to the request context
		// Dont Hardcode use shared const
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "username", username)
		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "token", tokenString)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok {
			http.Error(w, "Role not found in context", http.StatusUnauthorized)
			return
		}

		// check role if admin or not (case insencitive)
		if strings.EqualFold(role, "admin") {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
