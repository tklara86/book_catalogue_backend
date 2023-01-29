package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/tklara86/book_catalogue/internal/data"
	"github.com/tklara86/book_catalogue/internal/validator"
)

// createBookHandler creates new book
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Title      string `json:"title"`
		Status     int    `json:"status"`
		StatusID   int    `json:"status_id"`
		Authors    []int  `json:"authors"`
		Categories []int  `json:"categories"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:      input.Title,
		Status:     input.Status,
		StatusID:   input.StatusID,
		Authors:    input.Authors,
		Categories: input.Categories,
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	bookId, err := app.models.Book.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	bookAuthors := []data.BookAuthor{}
	for _, authorId := range input.Authors {
		ba := []data.BookAuthor{
			{
				BookId:   int64(bookId),
				AuthorId: int64(authorId),
			},
		}
		bookAuthors = append(bookAuthors, ba...)
	}

	_, err = app.models.Author.InsertBookAuthors(bookAuthors)
	if err != nil {
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	bookCategories := []data.BookCategory{}
	for _, categoryId := range input.Categories {

		bc := []data.BookCategory{
			{
				BookId:     int64(bookId),
				CategoryId: int64(categoryId),
			},
		}
		bookCategories = append(bookCategories, bc...)
	}

	_, err = app.models.Category.InsertBookCategories(bookCategories)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	bookResult := fmt.Sprintf("%q has been aded to your collection!", book.Title)

	jsonResponse := map[string]any{
		"book_id":        bookId,
		"book_title":     book.Title,
		"book_status":    book.Status,
		"client_message": bookResult,
	}

	err = app.writeToJSON(w, http.StatusCreated, envelope{"success": jsonResponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// showBookHandler get book by id
func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Book.GetBook(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	categories, err := app.models.Category.GetBookCategories(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	book.BookCategories = categories

	authors, err := app.models.Author.GetBookAuthors(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	book.BookAuthors = authors

	err = app.writeToJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getBooksHandler get all books
func (app *application) getBooksHandler(w http.ResponseWriter, r *http.Request) {

	books, err := app.models.Book.GetBooks()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	dateLayout := "02/01/2006"

	for _, b := range books {
		b.DateAdded = b.CreatedAt.UTC().Format(dateLayout)
		b.DateUpdated = b.UpdatedAt.UTC().Format(dateLayout)

		category, err := app.models.Category.GetBookCategories(b.ID)

		if err != nil {
			log.Fatal(err)
			return
		}
		for _, cat := range category {
			b.Categories = append(b.Categories, int(cat.ID))
		}
		b.BookCategories = append(b.BookCategories, category...) // append book categories

		bookAuthors, err := app.models.Author.GetBookAuthors(b.ID)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, aut := range bookAuthors {
			b.Authors = append(b.Authors, int(aut.AuthorID))
		}
		b.BookAuthors = append(b.BookAuthors, bookAuthors...) // append authors

	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"results": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ID []int `json:"ids"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	for _, id := range input.ID {
		err = app.models.Book.DeleteBook(int64(id))
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

	err = app.writeToJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Book.GetBook(id)
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
		Id         int     `json:"id"`
		Title      *string `json:"title"`
		Status     *int    `json:"status"`
		Categories []int   `json:"updated_categories"`
		Authors    []int   `json:"updated_authors"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Status != nil {
		book.Status = *input.Status
	}

	// update categories
	if input.Categories != nil {
		err = app.models.Category.DeleteBookCategories(id)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		bookCategories := []data.BookCategory{}
		for _, categoryId := range input.Categories {

			bc := []data.BookCategory{
				{
					BookId:     int64(id),
					CategoryId: int64(categoryId),
				},
			}
			bookCategories = append(bookCategories, bc...)
		}

		_, err = app.models.Category.InsertBookCategories(bookCategories)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	}

	// update authors
	if input.Authors != nil {
		err = app.models.Author.DeleteBookAuthors(id)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		bookAuthors := []data.BookAuthor{}
		for _, authorId := range input.Authors {

			ba := []data.BookAuthor{
				{
					BookId:   int64(id),
					AuthorId: int64(authorId),
				},
			}
			bookAuthors = append(bookAuthors, ba...)
		}

		_, err = app.models.Author.InsertBookAuthors(bookAuthors)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	}
	err = app.models.Book.UpdateBook(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	jsonResponse := map[string]any{
		"client_message": "book has been updated",
	}

	err = app.writeToJSON(w, http.StatusOK, envelope{"success": jsonResponse}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
