package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// GetUserById godoc
//
//	@Summary		Fetches user profile
//	@Description	Fetches user profile by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"User ID"
//	@Success		200		{object}	store.Users
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		404		{object}	error	"Record Not Found"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	userID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
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

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
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
//	@Success		200			{object}	store.Users
//	@Failure		400			{object}	error	"Bad Request"
//	@Failure		404			{object}	error	"Record Not Found"
//	@Failure		500			{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{username} [get]
func (app *application) getUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	user, err := app.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.recordNotFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Sets Follow user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			followID	body		string	true	"Follow ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	userID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	// followid from payload
	// var payload.followID
	var followID = uuid.New()

	userFollowData := &store.FollowUserParams{
		UserID:   userID,
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
//	@Summary		Unfollows a user
//	@Description	Sets Unfollow user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			followID	body		string	true	"Follow ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	userID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	// followid from payload
	// var payload.followID
	var followID = uuid.New()

	userUnfollowData := &store.UnfollowUserParams{
		UserID:   userID,
		FollowID: followID,
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
