package main

import (
	"errors"
	"log"
	"net/http"
)

// errorResponse represents a JSON error envelope returned by handlers.
type errorResponse struct {
	Error string `json:"error"`
}

var ErrInvalidRequest = errors.New("invalid request payload")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAccountNotActivated = errors.New("account is not activated")

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeJsonError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) statusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Internal Server Error occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Bad Request occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Not Found occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusNotFound, "the requested resource was not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusConflict, err.Error())
}
