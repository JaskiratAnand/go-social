package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JaskiratAnand/go-social/internal/mailer"
	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

type ReturnUserID struct {
	UserID uuid.UUID `json:"userID"`
}

// RegisterUser godoc
//
//	@Summary		Register user
//	@Description	Registers user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User Signup detailes"
//	@Success		201		{object}	ReturnUserID
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		401		{object}	error	"Invalid Credentials"
//	@Failure		409		{object}	error	"User Already Verified"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Router			/auth/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// existing user
	existingUser := false
	user, err := app.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.internalServerError(w, r, err)
			return
		}
	} else {
		existingUser = true
	}

	// verify password
	if existingUser {
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password)); err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		if user.Verified {
			app.customErrorResponse(w, r, http.StatusConflict, "user already verified")
			return
		}
	}

	var userID uuid.UUID
	if !existingUser { // creating new user
		// hash pwd
		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		createUserParam := &store.CreateUserParams{
			Username: payload.Username,
			Email:    payload.Email,
			Password: hash,
			RoleID:   1,
		}
		userID, err = app.store.CreateUser(ctx, *createUserParam)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	} else {
		userID = user.ID
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

	// send email
	isProdEnv := app.config.env == "production"
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, invitationParam.Token.String()) // redirect to /auth/activate/{token} from FE
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      payload.Username,
		ActivationURL: activationURL,
	}

	statusCode, err := app.mailer.Send(
		mailer.UserWelcomeTemplate,
		payload.Username,
		payload.Email,
		vars,
		!isProdEnv,
	)
	if err != nil {
		app.logger.Errorw("error sending activation email", "error", err)
		// rollback on email fail
		if err := app.store.DeleteUser(ctx, userID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		return
	}

	app.logger.Infow("Email sent", "status code", statusCode)

	if err := app.jsonResponse(w, http.StatusCreated, &ReturnUserID{UserID: userID}); err != nil {
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
//	@Failure		502	{object}	error	"Invite token expired"
//	@Router			/auth/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenParam := chi.URLParam(r, "token")

	ctx := r.Context()

	token, err := uuid.Parse(tokenParam)
	if err != nil {
		app.customErrorResponse(w, r, http.StatusBadRequest, "invalid token")
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
		app.customErrorResponse(w, r, http.StatusBadGateway, "invite token expired")
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

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// CreateToken godoc
//
//	@Summary		Create a Token
//	@Description	Create a Token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error					"Bad Request"
//	@Failure		401		{object}	error					"Unauthorized"
//	@Failure		500		{object}	error					"Server encountered a problem"
//	@Router			/auth/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// check existing user
	user, err := app.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.unauthorizedErrorResponse(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if !user.Verified {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}
	// verify user password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password)); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// gen token => add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send to client
	if err := app.jsonResponse(w, http.StatusOK, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
