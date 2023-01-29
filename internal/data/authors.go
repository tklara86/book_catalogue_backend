package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Author struct {
	AuthorID    int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	AuthorName  string    `json:"author_name"`
	AuthorBooks int       `json:"author_books,omitempty"`
	Description string    `json:"description,omitempty"`
	DateAdded   string    `json:"date_added"`
	DateUpdated string    `json:"date_updated"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type BookAuthor struct {
	BookId    int64     `json:"book_id"`
	AuthorId  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type AuthorModel struct {
	DB *sql.DB
}

func (a *AuthorModel) Insert(author *Author) (int, error) {

	query := `INSERT INTO cg_authors(first_name, last_name, description,created_at,updated_at) VALUES(TRIM(?), TRIM(?), TRIM(?), UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	args := []any{author.FirstName, author.LastName, author.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}

func (a *AuthorModel) InsertBookAuthors(ba []BookAuthor) (int, error) {
	query := `INSERT INTO cg_book_authors (book_id, author_id, created_at, updated_at) VALUES`

	args := []any{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	for _, v := range ba {
		args = append(args, v.BookId, v.AuthorId)
		numFields := 1

		for j := 0; j < numFields; j++ {
			query += `(?,` + `?` + `, UTC_TIMESTAMP(), UTC_TIMESTAMP()),`
		}
		query = query[:len(query)-1] + `,`
	}
	query = query[:len(query)-1]

	result, err := a.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil

}

func (a *AuthorModel) GetAuthors() ([]*Author, error) {
	query := `SELECT id, CONCAT(first_name, ' ' ,last_name) as author_name, first_name, last_name, description, created_at, updated_at FROM cg_authors`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	authors := []*Author{}

	for rows.Next() {
		auth := &Author{}

		err := rows.Scan(&auth.AuthorID, &auth.AuthorName, &auth.FirstName, &auth.LastName, &auth.Description, &auth.CreatedAt, &auth.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, auth)

	}

	return authors, nil

}

func (a *AuthorModel) GetBookAuthors(id int64) ([]*Author, error) {
	query := `SELECT CONCAT(a.first_name, ' ', a.last_name) as author_name, a.id, a.first_name, a.last_name, a.description, a.created_at, a.updated_at FROM cg_authors a
						LEFT JOIN cg_book_authors bk  ON bk.author_id = a.id
						WHERE bk.book_id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	results, err := a.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer results.Close()

	authors := []*Author{}

	for results.Next() {
		auth := &Author{}

		err := results.Scan(&auth.AuthorName, &auth.AuthorID, &auth.FirstName, &auth.LastName, &auth.Description, &auth.CreatedAt, &auth.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, auth)
	}

	return authors, nil

}

func (a *AuthorModel) DeleteAuthor(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM cg_authors WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
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

func (a *AuthorModel) GetAuthor(id int64) (*Author, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, CONCAT(first_name, ' ', last_name) as author_name, first_name, last_name, description, created_at, updated_at FROM cg_authors WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var author Author

	err := a.DB.QueryRowContext(ctx, query, id).Scan(&author.AuthorID, &author.AuthorName, &author.FirstName, &author.LastName, &author.Description, &author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &author, nil

}

func (a *AuthorModel) UpdateAuthor(author *Author) error {
	query := `UPDATE cg_authors SET first_name = ?, last_name = ?, description = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := a.DB.ExecContext(ctx, query, author.FirstName, author.LastName, author.Description, author.AuthorID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthorModel) GetAuthorNumberOfBooks(id int64) (int, error) {
	query := `SELECT COUNT(b.id) FROM cg_books b
						LEFT JOIN cg_book_authors ba ON ba.book_id = b.id
						WHERE ba.author_id = ?`

	var bookAuthorNumber int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, id).Scan(&bookAuthorNumber)
	if err != nil {
		return 0, err
	}

	return bookAuthorNumber, nil

}

func (a *AuthorModel) DeleteBookAuthors(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM cg_book_authors WHERE book_id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
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
