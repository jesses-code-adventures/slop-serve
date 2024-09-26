// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package model

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type GeneratedImage struct {
	ID        uuid.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Url       string
}

type User struct {
	ID           uuid.UUID
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
}
