package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck route
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// books routes
	router.HandlerFunc(http.MethodPost, "/v1/book", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", app.getBooksHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.getBookHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.deleteBookHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.updateBookHandler)

	// authors routes
	router.HandlerFunc(http.MethodPost, "/v1/author", app.createAuthorHandler)
	router.HandlerFunc(http.MethodGet, "/v1/authors", app.getAuthorsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/authors/:id", app.getAuthorHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/authors", app.deleteAuthorHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/authors/:id", app.updateAuthorHandler)

	// categories routes
	router.HandlerFunc(http.MethodPost, "/v1/category", app.createCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories", app.getCategoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", app.getCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/categories", app.deleteCategoryHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/categories/:id", app.updateCategoryHandler)

	return app.enableCORS(router)
}
