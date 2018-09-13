package middlewares

import (
	"context"
	"fmt"
	"go_rest_pg_starter/config"
	"go_rest_pg_starter/models"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type User struct {
	models.UserService
}

func getSignKey() string {
	cfg := config.GetConfig()
	return cfg.SigningKey
}

type UserWithToken struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"not null; unique_index"`
	UserEmail string `gorm:"not null; unique_index"`
}

func LookUpUserFromContext(ctx context.Context) *UserWithToken {
	temp := ctx.Value("logged_in_user")
	if temp != nil {
		user, ok := temp.(*UserWithToken)
		if ok {
			return user
		}
	}
	return nil
}

func PassSignKey(h http.Handler) http.Handler {
	signingKey := getSignKey()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "signingKey", signingKey)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

// Middleware to check JWT
func (us *User) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	signingKey := getSignKey()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(signingKey), nil
			})

		if err == nil {
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok || !token.Valid {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprint(w, "Token is not valid")
				}

				// We want to pass username extracted from JWT token to next handler:
				// Take the context out from the request
				ctx := r.Context()
				uid := claims["logged_in_user_id"]

				user, err := us.UserService.GetById(uint(uid.(float64)))
				if err != nil {
					next(w, r)
					return
				}

				mappedUser := &UserWithToken{
					ID:        user.ID,
					UserEmail: user.Email,
					Username:  user.Username,
				}

				ctx = context.WithValue(ctx, "logged_in_user", mappedUser)

				// Get new http.Request with the new context
				r = r.WithContext(ctx)

				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Token is not valid")
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized access to this resource")
		}
	})
}
