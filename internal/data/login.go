package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/robihdy/passman/internal/encryption"
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
	DB     *sql.DB
	aesKey string
}

func (m LoginModel) Insert(login *Login) error {
	query := `
        INSERT INTO logins (name, username, password, website) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []interface{}{login.Name, login.Username, encryption.Encrypt(login.Password, m.aesKey), login.Website}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&login.ID, &login.CreatedAt, &login.Version)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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

	login.Password = encryption.Decrypt(login.Password, m.aesKey)

	return &login, nil
}

func (m LoginModel) Update(login *Login) error {
	query := `
        UPDATE logins 
        SET name = $1, username = $2, password = $3, website = $4, version = version + 1
        WHERE id = $5
        RETURNING version`

	args := []interface{}{
		login.Name,
		login.Username,
		encryption.Encrypt(login.Password, m.aesKey),
		login.Website,
		login.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&login.Version); err != nil {
		return err
	}

	login.Password = encryption.Decrypt(login.Password, m.aesKey)

	return nil
}

func (m LoginModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM logins
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m LoginModel) GetAll(name string, username string, filters Filters) ([]*Login, error) {
	query := `
        SELECT id, created_at, name, username, password, website, version
        FROM logins
        ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logins := []*Login{}

	for rows.Next() {
		var login Login

		err := rows.Scan(
			&login.ID,
			&login.CreatedAt,
			&login.Name,
			&login.Username,
			&login.Password,
			&login.Website,
			&login.Version,
		)
		if err != nil {
			return nil, err
		}

		login.Password = encryption.Decrypt(login.Password, m.aesKey)

		logins = append(logins, &login)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logins, nil
}
