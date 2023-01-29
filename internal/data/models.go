package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Book     BookModel
	Author   AuthorModel
	Category CategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Book:     BookModel{DB: db},
		Author:   AuthorModel{DB: db},
		Category: CategoryModel{DB: db},
	}
}
