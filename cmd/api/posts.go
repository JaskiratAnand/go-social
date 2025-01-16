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

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	// setup user verification

	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// var userID uuid.UUID
	userID, _ := uuid.Parse("bf6ade0f-9ab6-49c3-bb5c-57808599c432")
	createPost := &store.CreatePostParams{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  userID,
		Tags:    payload.Tags,
	}

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	var post store.CreatePostRow
	var err error
	if post, err = app.store.CreatePost(ctx, *createPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	postID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	var post store.GetPostWithCommentsByIdRow
	if post, err = app.store.GetPostWithCommentsById(ctx, postID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.recordNotFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	// setup user verification

	idParam := chi.URLParam(r, "postID")

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	postID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	if err = app.store.DeletePostById(ctx, postID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.recordNotFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   string   `json:"title" validate:"omitempty,max=100"`
	Content string   `json:"content" validate:"omitempty,max=1000"`
	Tags    []string `json:"tags" validate:"omitempty"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	// setup user verification

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	idParam := chi.URLParam(r, "postID")

	postID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	var post store.GetPostsByIdRow
	post, err = app.store.GetPostsById(ctx, postID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	updatePost := &store.UpdatePostByIdParams{
		Title:     If(payload.Title != "", payload.Title, post.Title),
		Content:   If(payload.Content != "", payload.Content, post.Content),
		Tags:      If(payload.Tags != nil, payload.Tags, post.Tags),
		ID:        postID,
		UpdatedAt: post.UpdatedAt,
	}

	var updatedPost store.UpdatePostByIdRow
	updatedPost, err = app.store.UpdatePostById(ctx, *updatePost)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, updatedPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type createCommentPayload struct {
	UserID  uuid.UUID `json:"user_id" validate:"required"`
	Content string    `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	// setup user verification

	var payload createCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	idParam := chi.URLParam(r, "postID")

	postID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
	defer cancel()

	createComment := &store.CreateCommentParams{
		PostID:  postID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}

	var comment store.CreateCommentRow
	comment, err = app.store.CreateComment(ctx, *createComment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
