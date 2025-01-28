package main

import (
	"net/http"
)

// Health godoc
//
//	@Summary		Health endpoint
//	@Description	Check server health
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error	"Server encountered a problem"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
