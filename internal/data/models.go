package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Logins LoginModel
}

func NewModels(db *sql.DB, aesKey string) Models {
	return Models{
		Logins: LoginModel{DB: db, aesKey: aesKey},
	}
}
