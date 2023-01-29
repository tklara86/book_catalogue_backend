package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tklara86/book_catalogue/internal/data"
)

// createAuthorHandler creates new author
func (app *application) createAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	author := &data.Author{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Description: input.Description,
	}

	authorId, err := app.models.Author.Insert(author)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	authorResult := fmt.Sprintf("Author: %q has been added to your catalogue", author.FirstName+" "+author.LastName)
	jsonResponse := map[string]any{
		"author_id":      authorId,
		"client_message": authorResult,
	}

	err = app.writeToJSON(w, http.StatusCreated, envelope{"success": jsonResponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showAuthorHandler get author by id
func (app *application) getAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	author, err := app.models.Author.GetAuthor(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"author": author}, nil)
	if err != nil {
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

}

// getBooksHandler get all books
func (app *application) getAuthorsHandler(w http.ResponseWriter, r *http.Request) {

	authors, err := app.models.Author.GetAuthors()

	dateLayout := "02/01/2006"
	for _, author := range authors {

		author.DateAdded = author.CreatedAt.UTC().Format(dateLayout)
		author.DateUpdated = author.UpdatedAt.UTC().Format(dateLayout)

		numberOfBooks, _ := app.models.Author.GetAuthorNumberOfBooks(author.AuthorID)

		author.AuthorBooks = numberOfBooks
	}

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"results": authors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ID []int `json:"ids"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	for _, id := range input.ID {
		err = app.models.Author.DeleteAuthor(int64(id))
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"message": "author successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// func (app *application) updateAuthorHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	category, err := app.models.Category.GetCategory(id)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = app.writeToJSON(w, http.StatusOK, envelope{"category": category}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

func (app *application) updateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	author, err := app.models.Author.GetAuthor(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Id          int     `json:"id"`
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FirstName != nil {
		author.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		author.LastName = *input.LastName
	}

	if input.Description != nil {
		author.Description = *input.Description
	}

	// v := validator.New()

	// if data.ValidateCategory(v, category); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }

	err = app.models.Author.UpdateAuthor(author)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	jsonResponse := map[string]any{
		"client_message": "author has been updated",
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"success": jsonResponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
