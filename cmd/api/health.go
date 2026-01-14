package main

import (
	"log"
	"net/http"
)

// healthCheckHandler provides a simple health check endpoint.
//
//	@Summary		Health Check
//	@Description	Checks the health status of the API
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"version": "1.0.0",
		"env":     app.config.env,
	}
	if err := writeJson(w, http.StatusOK, data); err != nil {
		log.Print("failed to write health check response", "error", err)
		writeJson(w, http.StatusInternalServerError, err.Error())
	}
}
