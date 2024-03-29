package data

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/tklara86/book_catalogue/internal/validator"
)

type Book struct {
	ID             int64       `json:"id"`
	Title          string      `json:"title"`
	Subtitle       string      `json:"subtitle"`
	Description    string      `json:"description"`
	Image          string      `json:"image"`
	ISBN           string      `json:"isbn"`
	PageCount      int         `json:"page_count"`
	PublishedDate  string      `json:"published_date"`
	Status         int         `json:"status,omitempty"`
	StatusName     string      `json:"status_name"`
	StatusID       int         `json:"status_id"`
	Authors        []int       `json:"authors,omitempty"`
	Categories     []int       `json:"categories,omitempty"`
	BookCategories []*Category `json:"book_categories"`
	BookAuthors    []*Author   `json:"book_authors"`
	DateAdded      string      `json:"date_added"`
	DateUpdated    string      `json:"date_updated"`
	CreatedAt      time.Time   `json:"-"`
	UpdatedAt      time.Time   `json:"-"`
}

func ValidateBook(v *validator.Validator, book *Book) {

	v.Check(book.Title != "", "title", "Title cannot be empty")

	//	v.Check(book.Status > 0, "status", "Invalid Status")

	// v.Check(book.Authors != nil, "authors", "must be provided")
	// v.Check(len(book.Authors) >= 1, "authors", "must contain at least 1 author")

	// v.Check(len(book.Categories) >= 1, "categories", "must contain at least 1 category")
	// v.Check(book.Categories != nil, "categories", "must be provided")

	// v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type BookModel struct {
	DB *sql.DB
}

// Insert new book and returns new book id
func (b *BookModel) Insert(book *Book) (int, error) {
	query := `
    INSERT INTO cg_books(title,subtitle,description,page_count,image,published_date,isbn,status,status_id,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?, UTC_TIMESTAMP(), UTC_TIMESTAMP())
  `
	args := []any{book.Title, book.Subtitle, book.Description, book.PageCount, book.Image, book.PublishedDate, book.ISBN, book.Status, book.StatusID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (b *BookModel) GetBooks(qs url.Values) ([]*Book, error) {
	query := `SELECT DISTINCT(b.id), b.title, b.status, b.subtitle, b.description, b.page_count, b.image, b.published_date, b.isbn, b.status_id, b.created_at, b.updated_at FROM cg_books b`

	authors := strings.Split(qs.Get("authors"), ",")
	categories := strings.Split(qs.Get("categories"), ",")

	var authorIds []string
	var categoryIds []string

	if qs.Get("authors") != "" {
		query += ` LEFT JOIN cg_book_authors ba ON ba.book_id = b.id`
	}
	if qs.Get("categories") != "" {
		query += ` LEFT JOIN cg_book_categories bc ON bc.book_id = b.id`
	}

	if qs.Get("authors") != "" {
		authorIds = append(authorIds, authors...)
		query += ` WHERE ba.author_id IN ` + `(` + strings.Join(authorIds, ",") + `)`
	}
	if qs.Get("categories") != "" {
		categoryIds = append(categoryIds, categories...)
		query += ` AND bc.category_id IN ` + `(` + strings.Join(categoryIds, ",") + `)`
	}

	query += " ORDER BY b.title ASC"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	results, err := b.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer results.Close()

	books := []*Book{}

	for results.Next() {
		bk := &Book{}

		err := results.Scan(&bk.ID, &bk.Title, &bk.StatusName, &bk.Subtitle, &bk.Description, &bk.PageCount, &bk.Image, &bk.PublishedDate, &bk.ISBN, &bk.StatusID, &bk.CreatedAt, &bk.UpdatedAt)
		if err != nil {
			return nil, err
		}

		books = append(books, bk)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (b *BookModel) GetFilteredBooks(title string, authors []string, categories []string, filters Filters) ([]*Book, error) {
	query := `SELECT * FROM cg_books`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	results, err := b.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer results.Close()

	books := []*Book{}

	for results.Next() {
		bk := &Book{}

		err := results.Scan(&bk.ID, &bk.Title, &bk.StatusName, &bk.StatusID, &bk.CreatedAt, &bk.UpdatedAt)
		if err != nil {
			return nil, err
		}

		books = append(books, bk)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (b *BookModel) DeleteBook(id int64) error {

	if id < 0 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM cg_books WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}

func (b *BookModel) GetBook(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT * FROM cg_books WHERE id = ?`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).Scan(&book.ID, &book.Title, &book.StatusName, &book.StatusID, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil

}

func (b *BookModel) UpdateBook(book *Book) error {
	query := `UPDATE cg_books SET title = ?, status = ?, status_id = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := b.DB.ExecContext(ctx, query, book.Title, book.Status, book.StatusID, book.ID)
	if err != nil {
		return err
	}
	return nil
}
