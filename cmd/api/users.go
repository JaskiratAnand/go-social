package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserResponseType struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Verified  bool      `json:"verified"`
}

// GetUserById godoc
//
//	@Summary		Fetches user profile
//	@Description	Fetches user profile by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"User ID"
//	@Success		200		{object}	UserResponseType
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		404		{object}	error	"Record Not Found"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	ctx := r.Context()

	userID, err := uuid.Parse(idParam)
	if err != nil {
		app.customErrorResponse(w, r, http.StatusBadRequest, "invalid post-id")
		return
	}

	user, err := app.store.GetUserByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.recordNotFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	userResponse := &UserResponseType{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		Verified:  user.Verified,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userResponse); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetUserByUsername godoc
//
//	@Summary		Fetches user profile
//	@Description	Fetches user profile by username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string	true	"Username"
//	@Success		200			{object}	UserResponseType
//	@Failure		400			{object}	error	"Bad Request"
//	@Failure		404			{object}	error	"Record Not Found"
//	@Failure		500			{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	// @Router			/users/{username} [get]
func (app *application) getUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	ctx := r.Context()

	user, err := app.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.recordNotFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	userResponse := &UserResponseType{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		Verified:  user.Verified,
	}

	if err := app.jsonResponse(w, http.StatusOK, userResponse); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type FollowUnfollowUserPayload struct {
	FollowID uuid.UUID `json:"followID" validate:"required"`
}

// FollowUser godoc
//
//	@Summary		Follow user
//	@Description	Follow user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"Follow ID"
//	@Success		200		{object}	nil
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	ctx := r.Context()

	followID, err := uuid.Parse(idParam)
	if err != nil {
		app.customErrorResponse(w, r, http.StatusBadRequest, "invalid post-id")
		return
	}

	user := app.GetUserFromCtx(r)

	userFollowData := &store.FollowUserParams{
		UserID:   user.ID,
		FollowID: followID,
	}

	err = app.store.FollowUser(ctx, *userFollowData)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollow user
//	@Description	Unfollow user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"Unfollow ID"
//	@Success		200		{object}	nil
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	ctx := r.Context()

	unfollowID, err := uuid.Parse(idParam)
	if err != nil {
		app.customErrorResponse(w, r, http.StatusBadRequest, "invalid post-id")
		return
	}

	user := app.GetUserFromCtx(r)

	userUnfollowData := &store.UnfollowUserParams{
		UserID:   user.ID,
		FollowID: unfollowID,
	}

	err = app.store.UnfollowUser(ctx, *userUnfollowData)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
