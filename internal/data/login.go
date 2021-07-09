package data

import (
	"database/sql"
	"errors"
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

func (m LoginModel) Insert(login *Login) error {
	query := `
        INSERT INTO logins (name, username, password, website) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []interface{}{login.Name, login.Username, login.Password, login.Website}

	return m.DB.QueryRow(query, args...).Scan(&login.ID, &login.CreatedAt, &login.Version)
}

func (m LoginModel) Get(id int64) (*Login, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, name, username, password, website, version
        FROM logins
        WHERE id = $1`

	var login Login

	err := m.DB.QueryRow(query, id).Scan(
		&login.ID,
		&login.CreatedAt,
		&login.Name,
		&login.Username,
		&login.Password,
		&login.Website,
		&login.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &login, nil
}

// Add a placeholder method for updating a specific record in the logins table.
func (m LoginModel) Update(login *Login) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the logins table.
func (m LoginModel) Delete(id int64) error {
	return nil
}
