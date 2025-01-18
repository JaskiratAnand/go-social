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
