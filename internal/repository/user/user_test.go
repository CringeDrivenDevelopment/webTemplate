package user

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexedwards/argon2id"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/model"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, mock)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, gormDB)

	userRepo := NewRepository(gormDB)

	passwordHash, err := argon2id.CreateHash("Very_strong_password1235", argon2id.DefaultParams)
	id := ulid.Make().String()
	email := "goodemail@gmail.com"
	assert.Nil(t, err)

	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash") VALUES ($1,$2,$3)`)).WithArgs(id, email, passwordHash).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash") VALUES ($1,$2,$3)`)).WithArgs(id, email, passwordHash).WillReturnError(fmt.Errorf("user with that email already exists"))
	mock.ExpectRollback()
	err = userRepo.Create(context.Background(), model.User{ID: id, Email: email, PasswordHash: passwordHash})
	assert.Nil(t, err)
	err = userRepo.Create(context.Background(), model.User{ID: id, Email: email, PasswordHash: passwordHash})
	assert.NotNil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, mock)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, gormDB)

	userRepo := NewRepository(gormDB)

	passwordHash, err := argon2id.CreateHash("Very_strong_password1235", argon2id.DefaultParams)
	id := ulid.Make().String()
	email := "goodemail@gmail.com"
	assert.Nil(t, err)

	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash") VALUES ($1,$2,$3)`)).WithArgs(id, email, passwordHash).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE LOWER(users.email) = $1 ORDER BY  "users"."id" LIMIT $2`)).WithArgs(email, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).AddRow(id, email, passwordHash))

	err = userRepo.Create(context.Background(), model.User{ID: id, Email: email, PasswordHash: passwordHash})
	assert.Nil(t, err)
	user, err := userRepo.GetUserByEmail(context.Background(), email)
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
	assert.Equal(t, model.User{ID: id, Email: email, PasswordHash: passwordHash}, user)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, mock)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, gormDB)

	userRepo := NewRepository(gormDB)

	passwordHash, err := argon2id.CreateHash("Very_strong_password1235", argon2id.DefaultParams)
	id := ulid.Make().String()
	email := "goodemail@gmail.com"
	assert.Nil(t, err)

	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash") VALUES ($1,$2,$3)`)).WithArgs(id, email, passwordHash).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY  "users"."id" LIMIT $2`)).WithArgs(id, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).AddRow(id, email, passwordHash))

	err = userRepo.Create(context.Background(), model.User{ID: id, Email: email, PasswordHash: passwordHash})
	assert.Nil(t, err)
	user, err := userRepo.GetUserByID(context.Background(), id)
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
	assert.Equal(t, model.User{ID: id, Email: email, PasswordHash: passwordHash}, user)
}
