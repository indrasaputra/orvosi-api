package repository_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

type UserInserter_Executor struct {
	repo *repository.UserInserter
	sql  sqlmock.Sqlmock
}

func TestNewUserInserter(t *testing.T) {
	t.Run("successfully create an instance of UserInserter", func(t *testing.T) {
		exec := createUserInserterExecutor()
		assert.NotNil(t, exec.repo)
	})
}

func TestUserInserter_InsertOrIgnore(t *testing.T) {
	t.Run("can't proceed due to nil user", func(t *testing.T) {
		exec := createUserInserterExecutor()

		err := exec.repo.InsertOrIgnore(context.Background(), nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyUser, err)
	})

	t.Run("database returns error", func(t *testing.T) {
		exec := createUserInserterExecutor()

		exec.sql.ExpectExec(`INSERT INTO users \(name, email, google_id, created_at, updated_at, created_by, updated_by\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) ON CONFLICT \(email\) DO NOTHING`).
			WillReturnError(errors.New("fail to insert to database"))

		user := createValidUser()
		err := exec.repo.InsertOrIgnore(context.Background(), user)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer.Code, err.Code)
	})

	t.Run("successfully insert or ignore user", func(t *testing.T) {
		exec := createUserInserterExecutor()

		exec.sql.ExpectExec(`INSERT INTO users \(name, email, google_id, created_at, updated_at, created_by, updated_by\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) ON CONFLICT \(email\) DO NOTHING`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		user := createValidUser()
		err := exec.repo.InsertOrIgnore(context.Background(), user)

		assert.Nil(t, err)
	})
}

func createValidUser() *entity.User {
	return &entity.User{
		ID:       hashids.ID(1),
		Email:    "email@provider.com",
		Name:     "User 1",
		GoogleID: "super-long-google-id",
	}
}

func createUserInserterExecutor() *UserInserter_Executor {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Panicf("error opening a stub database connection: %v\n", err)
	}

	repo := repository.NewUserInserter(db)
	return &UserInserter_Executor{
		repo: repo,
		sql:  mock,
	}
}
