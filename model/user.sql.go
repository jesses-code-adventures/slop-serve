// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package model

import (
	"context"

	"github.com/gofrs/uuid"
)

const hashedPasswordGet = `-- name: HashedPasswordGet :one
select users.password_hash, users.id user_id
from users
where users.email = $1
`

type HashedPasswordGetRow struct {
	PasswordHash string
	UserID       uuid.UUID
}

func (q *Queries) HashedPasswordGet(ctx context.Context, email string) (HashedPasswordGetRow, error) {
	row := q.db.QueryRow(ctx, hashedPasswordGet, email)
	var i HashedPasswordGetRow
	err := row.Scan(&i.PasswordHash, &i.UserID)
	return i, err
}

const userCreate = `-- name: UserCreate :one
INSERT INTO users (first_name, last_name, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type UserCreateParams struct {
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
}

func (q *Queries) UserCreate(ctx context.Context, arg UserCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, userCreate,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.PasswordHash,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const userCreateWithId = `-- name: UserCreateWithId :one
INSERT INTO users (first_name, last_name, email, password_hash, id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type UserCreateWithIdParams struct {
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	ID           uuid.UUID
}

func (q *Queries) UserCreateWithId(ctx context.Context, arg UserCreateWithIdParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, userCreateWithId,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.PasswordHash,
		arg.ID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
