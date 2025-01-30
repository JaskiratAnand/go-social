package main

import (
	"context"
	"net/http"
	"time"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

// RegisterUser godoc
//
//	@Summary		Register user
//	@Description	Registers user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body RegisterUserPayload
//	@Success		201		{object}	string	true "User Registered"
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Router			/auth/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	// existing user
	user, err := app.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if user.Email != "" {
		// user already exists
	}

	// hash pwd
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// creating new user
	createUserParam := &store.CreateUserParams{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hash,
	}
	userID, err := app.store.CreateUser(ctx, *createUserParam)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// creating verification token
	invitationParam := &store.CreateInvitationParams{
		Token:   uuid.New(),
		UserID:  userID,
		Expiary: time.Now().Add(app.config.mail.exp),
	}

	if err := app.store.CreateInvitation(ctx, *invitationParam); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
