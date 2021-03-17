package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/indrasaputra/orvosi-api/usecase"
	"github.com/stretchr/testify/assert"
)

type SignIn_Executor struct {
	usecase *usecase.Signer
	repo    *mock_usecase.MockInsertUserRepository
}

func TestNewSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of SignIn", func(t *testing.T) {
		exec := createSignInExecutor(ctrl)
		assert.NotNil(t, exec.usecase)
	})
}

func TestSigner_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("user is empty/nil", func(t *testing.T) {
		exec := createSignInExecutor(ctrl)

		err := exec.usecase.SignIn(context.Background(), nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyUser, err)
	})

	t.Run("repo returns error", func(t *testing.T) {
		exec := createSignInExecutor(ctrl)

		user := createValidUser()
		exec.repo.EXPECT().InsertOrIgnore(context.Background(), user).Return(entity.ErrInternalServer)

		err := exec.usecase.SignIn(context.Background(), user)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
	})

	t.Run("successfully insert or ignore user", func(t *testing.T) {
		exec := createSignInExecutor(ctrl)

		user := createValidUser()
		exec.repo.EXPECT().InsertOrIgnore(context.Background(), user).Return(nil)

		err := exec.usecase.SignIn(context.Background(), user)

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

func createSignInExecutor(ctrl *gomock.Controller) *SignIn_Executor {
	r := mock_usecase.NewMockInsertUserRepository(ctrl)
	u := usecase.NewSigner(r)

	return &SignIn_Executor{
		usecase: u,
		repo:    r,
	}
}
