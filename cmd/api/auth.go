package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/go-chi/chi/v5"
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
//	@Param			payload	body	RegisterUserPayload	true
//	@Success		201
//	@Failure		400	{object}	error	"Bad Request"
//	@Failure		500	{object}	error	"Server encountered a problem"
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
	// user, err := app.store.GetUserByEmail(ctx, payload.Email)
	// if (err != nil) && !errors.Is(err, sql.ErrNoRows) {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	// if user.Verified {
	// 	err := errors.New("user already exists")
	// 	app.badRequestResponse(w, r, err)
	// }

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

// ActivateUser godoc
//
//	@Summary		Activate user
//	@Description	Activate user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	path	string	true	"Invite token"
//	@Success		204
//	@Failure		404	{object}	error	"Invalid token"
//	@Failure		500	{object}	error	"Server encountered a problem"
//	@Router			/auth/activate/{token} [post]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenParam := chi.URLParam(r, "token")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	token, err := uuid.Parse(tokenParam)
	if err != nil {
		app.customError(w, r, http.StatusBadRequest, "invalid token")
		return
	}

	invite, err := app.store.GetInvitationByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.recordNotFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if time.Now().After(invite.Expiary) {
		app.customError(w, r, http.StatusBadGateway, "invite token expired")
		return
	}

	err = app.store.ActivateUser(ctx, invite.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.DeleteInvitationByUserId(ctx, invite.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, "User Activated"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
