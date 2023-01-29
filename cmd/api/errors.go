package main

import (
	"fmt"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeToJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	app.errorResponse(w, r, http.StatusInternalServerError, errorMessage{
		Message: "the server encountered a problem and could not process your request",
		Status:  http.StatusInternalServerError,
	})
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	app.errorResponse(w, r, http.StatusNotFound, errorMessage{
		Message: "the requested resource could not be found",
		Status:  http.StatusNotFound,
	})
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	app.errorResponse(w, r, http.StatusMethodNotAllowed, errorMessage{
		Message: fmt.Sprintf("the %s method is not supported for this resource", r.Method),
		Status:  http.StatusMethodNotAllowed,
	})
}
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
