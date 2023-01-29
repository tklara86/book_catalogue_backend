package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tklara86/book_catalogue/internal/data"
	"github.com/tklara86/book_catalogue/internal/validator"
)

// createCategoryHandler creates new category
func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	category := &data.Category{
		Name: input.Name,
	}

	categoryId, err := app.models.Category.Insert(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	categoryResult := fmt.Sprintf("%q has been added to your catalogue", category.Name)
	jsonResponse := map[string]any{
		"client_message": categoryResult,
		"category_id":    categoryId,
	}

	err = app.writeToJSON(w, http.StatusCreated, envelope{"success": jsonResponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showCategoryHandler get category by id
func (app *application) showCategoryHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	fmt.Fprintf(w, "show category of id %d", id)
}

// getCategoriesHandler get all categories
func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	categories, err := app.models.Category.GetCategories()
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	dateLayout := "02/01/2006"
	for _, cat := range categories {
		cat.DateAdded = cat.CreatedAt.UTC().Format(dateLayout)
		cat.DateUpdated = cat.UpdatedAt.UTC().Format(dateLayout)
		bookCategories, _ := app.models.Category.GetBooksInCategory(cat.ID)

		cat.BooksInCategory = bookCategories
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"results": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ID []int `json:"ids"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	for _, id := range input.ID {
		err = app.models.Category.DeleteCategory(int64(id))
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

	err = app.writeToJSON(w, http.StatusOK, envelope{"message": "category successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Category.GetCategory(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Category.GetCategory(id)
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
		Id   int     `json:"id"`
		Name *string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Category.UpdateCategory(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	jsonResponse := map[string]any{
		"client_message": "category name has been updated",
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"success": jsonResponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
