package db

import (
	"net/url"

	"github.com/gofrs/uuid"
)

type Database interface {
	// Register will take a first name, a last name, an email and an encrypted password. returns the id or an error.
	UserCreate(firstName string, lastName, email string, password string) (uuid.UUID, error)
	HashedPasswordGet(email string) (uuid.UUID, string, error)
	GenImageCreate(url url.URL) (uuid.UUID, error)
}
