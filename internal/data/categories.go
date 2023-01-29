package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/tklara86/book_catalogue/internal/validator"
)

type Category struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	BooksInCategory int       `json:"books_in_category,omitempty"`
	DateAdded       string    `json:"date_added"`
	DateUpdated     string    `json:"date_updated"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

type BookCategory struct {
	ID         int64     `json:"id"`
	BookId     int64     `json:"book_id"`
	CategoryId int64     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type CategoryModel struct {
	DB *sql.DB
}

func ValidateCategory(v *validator.Validator, category *Category) {

	v.Check(category.Name != "", "name", "Name cannot be empty")

}

func (c *CategoryModel) Insert(category *Category) (int, error) {
	query := `INSERT INTO cg_categories (name,created_at,updated_at) VALUES (TRIM(?), UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	args := []any{category.Name}

	result, err := c.DB.Exec(query, args...)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *CategoryModel) InsertBookCategories(bc []BookCategory) (int, error) {
	query := `
	INSERT INTO cg_book_categories (book_id, category_id, created_at, updated_at)
	VALUES `

	args := []any{}

	for _, v := range bc {
		args = append(args, v.BookId, v.CategoryId)
		numFields := 1

		for j := 0; j < numFields; j++ {
			query += `(?,` + `?` + `, UTC_TIMESTAMP(), UTC_TIMESTAMP()),`
		}
		query = query[:len(query)-1] + `,`
	}
	query = query[:len(query)-1]

	result, err := c.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *CategoryModel) GetCategories() ([]*Category, error) {
	query := `SELECT * FROM cg_categories`

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		cat := &Category{}

		err = rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil

}

func (c *CategoryModel) GetCategory(id int64) (*Category, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT * FROM cg_categories WHERE id = ?`

	var category Category

	err := c.DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil

}

func (c *CategoryModel) GetBooksInCategory(id int64) (int, error) {
	query := `SELECT COUNT(b.id) FROM cg_books b
						LEFT JOIN cg_book_categories bc ON bc.book_id = b.id
						WHERE bc.category_id = ?`

	var bookCategoryNumber int

	err := c.DB.QueryRow(query, id).Scan(&bookCategoryNumber)
	if err != nil {
		return 0, err
	}

	return bookCategoryNumber, nil

}

func (c *CategoryModel) GetBookCategories(id int64) ([]*Category, error) {
	query := `SELECT c.id, c.name, c.created_at, c.updated_at FROM cg_categories c
						LEFT JOIN cg_book_categories bc ON bc.category_id = c.id
						WHERE bc.book_id = ?`

	results, err := c.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer results.Close()

	categories := []*Category{}

	for results.Next() {
		cat := &Category{}

		err := results.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil

}

func (c *CategoryModel) UpdateCategory(category *Category) error {
	query := `UPDATE cg_categories SET name = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`

	_, err := c.DB.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryModel) DeleteCategory(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM cg_categories WHERE id = ?`

	result, err := c.DB.Exec(query, id)
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

func (c *CategoryModel) DeleteBookCategories(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM cg_book_categories WHERE book_id = ?`

	result, err := c.DB.Exec(query, id)
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
