package db

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jesses-code-adventures/slop/model"
)

type PgDatabase struct {
	ctx     context.Context
	conn    *pgx.Conn
	queries *model.Queries
}

func NewPgDatabase(ctx context.Context) PgDatabase {
	conn, err := pgx.Connect(ctx, fmt.Sprintf("user=%s dbname=%s sslmode=verify-full", os.Getenv("DB_USER"), os.Getenv("DB_NAME")))
	if err != nil {
		panic(err)
	}
	queries := model.New(conn)
	return PgDatabase{
		ctx:     ctx,
		conn:    conn,
		queries: queries,
	}
}

func (d PgDatabase) Close() error {
	return d.conn.Close(d.ctx)
}

func (d PgDatabase) UserCreate(firstName, lastName, email, hashedPassword string) (uuid.UUID, error) {
	return d.queries.UserCreateWithId(d.ctx, model.UserCreateWithIdParams{FirstName: firstName, LastName: lastName, Email: email, PasswordHash: hashedPassword, ID: uuid.FromStringOrNil(os.Getenv("TEST_USER_ID"))})
	// return d.queries.UserCreate(d.ctx, model.UserCreateParams{FirstName: firstName, LastName: lastName, Email: email, PasswordHash: hashedPassword})
}

func (d PgDatabase) HashedPasswordGet(email string) (uuid.UUID, string, error) {
	resp, err := d.queries.HashedPasswordGet(d.ctx, email)
	if err != nil {
		return uuid.UUID{}, "", err
	}
	return resp.UserID, resp.PasswordHash, nil
}

func (d PgDatabase) GenImageCreate(url string) (uuid.UUID, error) {
	return d.queries.GenImageCreate(d.ctx, url)
}
