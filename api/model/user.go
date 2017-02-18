package model

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/spolu/cumulo/api"
	"github.com/spolu/cumulo/lib/db"
	"github.com/spolu/cumulo/lib/errors"
	"github.com/spolu/cumulo/lib/token"
)

// User represents a user object. Users are uniquely assoiated with phone
// numbers (for now).
type User struct {
	Token   string
	Created time.Time

	Phone  string
	PubKey *string
}

// NewUserResource generates a new user resource.
func NewUserResource(
	ctx context.Context,
	user *User,
) api.UserResource {
	return api.UserResource{
		Token:   user.Token,
		Created: user.Created.UnixNano() / api.TimeResolutionNs,
		Phone:   user.Phone,
		PubKey:  user.PubKey,
	}
}

// CreateUser creates and stores a new User object.
func CreateUser(
	ctx context.Context,
	phone string,
) (*User, error) {
	user := User{
		Token:   token.New("user"),
		Created: time.Now().UTC(),
		Phone:   phone,
	}

	ext := db.Ext(ctx, "api")
	if _, err := sqlx.NamedExec(ext, `
INSERT INTO users
  (token, created, phone)
VALUES
  (:token, :created, :phone)
`, user); err != nil {
		switch err := err.(type) {
		case *pq.Error:
			if err.Code.Name() == "unique_violation" {
				return nil, errors.Trace(ErrUniqueConstraintViolation{err})
			}
		case sqlite3.Error:
			if err.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, errors.Trace(ErrUniqueConstraintViolation{err})
			}
		}
		return nil, errors.Trace(err)
	}

	return &user, nil
}

// Save updates the object database representation with the in-memory values.
func (u *User) Save(
	ctx context.Context,
) error {
	ext := db.Ext(ctx, "register")
	_, err := sqlx.NamedExec(ext, `
UPDATE users
SET phone = :phone, pubkey = :pubkey
WHERE token = :token
`, u)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// LoadUserByPhone attempts to load a user with the given phone.
func LoadUserByPhone(
	ctx context.Context,
	phone string,
) (*User, error) {
	user := User{
		Phone: phone,
	}

	ext := db.Ext(ctx, "api")
	if rows, err := sqlx.NamedQuery(ext, `
SELECT *
FROM users
WHERE phone = :phone
`, user); err != nil {
		return nil, errors.Trace(err)
	} else if !rows.Next() {
		return nil, nil
	} else if err := rows.StructScan(&user); err != nil {
		defer rows.Close()
		return nil, errors.Trace(err)
	} else if err := rows.Close(); err != nil {
		return nil, errors.Trace(err)
	}

	return &user, nil
}
