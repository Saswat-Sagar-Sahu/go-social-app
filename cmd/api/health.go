package main

import (
	"log"
	"net/http"
)

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
