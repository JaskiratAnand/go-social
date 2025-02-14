package main

import (
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

// CreatePost godoc
//
//	@Summary		Create a Post
//	@Description	Creates new posts
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			CreatePost	body		CreatePostPayload	true	"Create Post Payload"
//	@Success		201			{object}	store.CreatePostRow
//	@Failure		400			{object}	error	"Bad Request"
//	@Failure		500			{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// get user id
	user := app.GetUserFromCtx(r)

	createPost := &store.CreatePostParams{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  user.ID,
		Tags:    payload.Tags,
	}

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

// GetPost godoc
//
//	@Summary		Fetch Post
//	@Description	Fetch post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		string	true	"Post ID"
//	@Success		200		{object}	store.GetPostWithCommentsByIdRow
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		404		{object}	error	"Record Not Found"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")

	ctx := r.Context()

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

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Delete Post
//	@Description	Deletes post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	string	true	"Post ID"
//	@Success		204
//	@Failure		400	{object}	error	"Bad Request"
//	@Failure		404	{object}	error	"Record Not Found"
//	@Failure		500	{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")

	ctx := r.Context()

	user := app.GetUserFromCtx(r)

	postID, err := uuid.Parse(idParam)
	if err != nil {
		err = errors.New("invalid post-id")
		app.badRequestResponse(w, r, err)
		return
	}

	deletePostParam := &store.DeletePostByIdParams{
		ID:     postID,
		UserID: user.ID,
	}

	if err = app.store.DeletePostById(ctx, *deletePostParam); err != nil {
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

// UpdatePost godoc
//
//	@Summary		Fetch Post
//	@Description	Fetch post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID		path		string				true	"Post ID"
//	@Param			updatePost	body		UpdatePostPayload	true	"Update Post Payload"
//	@Success		200			{object}	store.UpdatePostByIdRow
//	@Failure		400			{object}	error	"Bad Request"
//	@Failure		401			{object}	error	"Unauthorized"
//	@Failure		500			{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx := r.Context()

	user := app.GetUserFromCtx(r)

	var post store.GetPostsByIdRow
	post, err = app.store.GetPostsById(ctx, postID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if user.ID != post.UserID {
		app.customErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
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

	if err := app.jsonResponse(w, http.StatusOK, updatedPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateCommentPayload struct {
	Content string `json:"content" validate:"omitempty,max=1000"`
}

// CreateComment godoc
//
//	@Summary		Create Comment
//	@Description	Create comment on posts
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		string	true	"Post ID"
//	@Param			content	body		string	true	"Content Payload"
//	@Success		200		{object}	store.CreateCommentRow
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server encountered a problem"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
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

	ctx := r.Context()

	user := app.GetUserFromCtx(r)

	createComment := &store.CreateCommentParams{
		PostID:  postID,
		UserID:  user.ID,
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
