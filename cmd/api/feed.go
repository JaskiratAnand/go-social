package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	// get user id
	userID, _ := uuid.Parse("9415cb50-29b8-486c-beb7-8a149e75cde1")

	feed, err := app.store.GetUserFeed(ctx, userID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
