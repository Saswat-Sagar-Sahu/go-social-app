package main

import (
	"log"
	"net/http"
)

func (app *application) statusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Internal Server Error occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Bad Request occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Bad Request occurred %s path : %s error : %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusNotFound, "the requested resource was not found")
}
