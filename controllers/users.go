package controllers

import (
	"encoding/json"
	"net/http"

	"go_rest_pg_starter/auth"
	"go_rest_pg_starter/email"
	"go_rest_pg_starter/middlewares"
	"go_rest_pg_starter/models"
)

type UserWithToken struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	UserEmail string `json:"user_email"`
	Token     string `json:"token"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type Users struct {
	us      models.UserService
	emailer *email.Client
}

type SignupUser struct {
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginUser struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type ResetPasswordUser struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		us:      us,
		emailer: emailer,
	}
}

// POST /api/signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var signupUser SignupUser

	err := json.NewDecoder(r.Body).Decode(&signupUser)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot signup.")
		return
	}

	user := models.User{
		Username: signupUser.Username,
		Email:    signupUser.Email,
		Password: signupUser.Password,
	}

	// Create the user in our database
	err = u.us.Create(&user)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot signup.")
		return
	}

	// Send an email to the user
	err = u.emailer.Welcome(user.Username, user.Email)
	if err != nil {
		sendErrorResponse(w, http.StatusFound, "Could not send welcome email.")
		return
	}

	// Issue JWT and let the user sign-in
	signingKey := r.Context().Value("signingKey").(string)
	err = u.signIn(w, &user, signingKey)
	if err != nil {
		sendErrorResponse(w, http.StatusFound, "Cannot signin.")
		return
	}
}

// POST /api/login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser LoginUser

	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot login.")
		return
	}

	user, err := u.us.Authenticate(loginUser.Email, loginUser.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			sendErrorResponse(w, http.StatusNotFound, "User not found.")
		default:
			sendErrorResponse(w, http.StatusForbidden, "Cannot login.")
		}
		return
	}

	signingKey := r.Context().Value("signingKey").(string)
	err = u.signIn(w, user, signingKey)
	if err != nil {
		sendErrorResponse(w, http.StatusFound, "Cannot sign-in.")
		return
	}
}

// Set token if the user doesn't have one and set it to cookie
func (u *Users) signIn(w http.ResponseWriter, user *models.User, signingKey string) error {
	tokenString, err := auth.IssueJWT(user, signingKey)

	if err != nil {
		return err
	}

	response := UserWithToken{
		UserID:    user.ID,
		UserEmail: user.Email,
		Username:  user.Username,
		Token:     tokenString,
	}
	setSuccessStatus(w, http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
}

// GET /api/me
func (u *Users) Me(w http.ResponseWriter, r *http.Request) {
	user := middlewares.LookUpUserFromContext(r.Context())
	if user == nil {
		sendErrorResponse(w, http.StatusFound, "User not found.")
		return
	}

	setSuccessStatus(w, http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// POST /api/forgot_password
// Send user an email with a token to reset the password
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var resetPasswordUser ResetPasswordUser

	// Get user info
	err := json.NewDecoder(r.Body).Decode(&resetPasswordUser)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot reset. Could not get user info.")
		return
	}

	// Create a token to start resetting user password
	token, err := u.us.InitiateReset(resetPasswordUser.Email)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot reset. Could not issue a reset token.")
		return
	}

	// Email the user the password reset token
	err = u.emailer.ResetPassword(resetPasswordUser.Email, token)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot reset. Could not email to the user.")
		return
	}
}

// POST /api/update_password
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var resetPasswordUser ResetPasswordUser

	// Get user info
	err := json.NewDecoder(r.Body).Decode(&resetPasswordUser)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot get user info.")
		return
	}

	// Reset user password (Update with new password)
	token := r.URL.Query().Get("token") // Get token from URL param
	user, err := u.us.CompleteReset(token, resetPasswordUser.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusForbidden, "Cannot update with new password.")
		return
	}

	// Let user sign-in
	signingKey := r.Context().Value("signingKey").(string)
	u.signIn(w, user, signingKey)
	if err != nil {
		sendErrorResponse(w, http.StatusFound, "Updated user password. Cannot signin.")
		return
	}
}
