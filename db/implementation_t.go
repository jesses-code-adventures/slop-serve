package db

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type TestDatabase struct {
	logger *slog.Logger
}

func NewTestDatabase(logger *slog.Logger) TestDatabase {
	return TestDatabase{logger: logger}
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func getTestUserUUID() uuid.UUID {
	return Must(uuid.FromString(os.Getenv("TEST_USER_ID")))
}

func getTestGenImageUUID() uuid.UUID {
	return Must(uuid.FromString("TEST_GEN_IMAGE_ID"))
}

func (d TestDatabase) UserRegister(firstName string, lastName, email string, password string) (id uuid.UUID, err error) {
	d.logger.Debug(fmt.Sprintf("db entry:\nfirst_name: %s\nlast_name: %s, email: %s, password: %s", firstName, lastName, email, password))
	id = getTestUserUUID()
	return
}

func (d TestDatabase) UserLogin(email string, password string) (id uuid.UUID, err error) {
	d.logger.Debug(fmt.Sprintf("login:\nemail: %s\npassword: %s", email, password))
	id = getTestUserUUID()
	return
}

// TODO: should return the whole image
func (d TestDatabase) GenImage(url url.URL) (id uuid.UUID, err error) {
	createdAt := time.Now()
	d.logger.Debug(fmt.Sprintf("create gen image db\ncreated_at: %s, url: %s", createdAt, url.String()))
	id = getTestGenImageUUID()
	return
}
