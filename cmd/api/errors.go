package main

import (
	"errors"
	"net/http"
)

func (app *application) recordNotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusNotFound, "record not found")
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) customError(w http.ResponseWriter, r *http.Request, status int, errMessage string) {
	err := errors.New(errMessage)
	app.logger.Warnw("error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, status, err.Error())
}
