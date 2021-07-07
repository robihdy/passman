package data

import (
	"database/sql"
	"time"

	"github.com/robihdy/passman/internal/validator"
	"gopkg.in/guregu/null.v4"
)

type Login struct {
	ID        int64       `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	Name      string      `json:"name"`
	Username  string      `json:"username"`
	Password  string      `json:"password"`
	Website   null.String `json:"website"`
	Version   int32       `json:"version"`
}

func ValidateLogin(v *validator.Validator, l *Login) {
	v.Check(l.Name != "", "name", "must be provided")
	v.Check(len(l.Name) <= 255, "name", "must not be more than 255 bytes long")

	v.Check(l.Username != "", "username", "must be provided")
	v.Check(len(l.Username) <= 255, "username", "must not be more than 255 bytes long")

	v.Check(l.Password != "", "password", "must be provided")
	v.Check(len(l.Password) <= 255, "password", "must not be more than 255 bytes long")
	v.Check(len(l.Password) >= 8, "password", "must be more than or equal to 8 bytes long")
}

type LoginModel struct {
	DB *sql.DB
}
