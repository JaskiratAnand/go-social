package main

import (
	"context"
	"net/http"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/google/uuid"
)

// GetUserFeed godoc
//
//	@Summary		Fetches user feed
//	@Description	Fetches user feed with following posts
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"User ID"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Success		200		{object}	[]store.GetUserFeedRow
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	// get user id
	userID, _ := uuid.Parse("9415cb50-29b8-486c-beb7-8a149e75cde1")

	getUserFeedParams := &store.GetUserFeedParams{
		UserID:  userID,
		Column2: fq.Search, // Search
		Tags:    fq.Tags,
		Limit:   int64(fq.Limit),
		Offset:  int64(fq.Offset),
	}

	feed, err := app.store.GetUserFeed(ctx, *getUserFeedParams)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
