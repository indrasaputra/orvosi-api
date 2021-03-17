package usecase

import (
	"context"

	"github.com/indrasaputra/orvosi-api/entity"
)

// SignIn defines the business logic to sign in.
type SignIn interface {
	// SignIn signs a user in to the system.
	SignIn(ctx context.Context, user *entity.User) *entity.Error
}

// InsertUserRepository defines the business logic
// to insert a user into a repository.
type InsertUserRepository interface {
	// InsertOrIgnore inserts a user into the repository.
	// If the user already exists in the repository, it ignores and returns nil error.
	InsertOrIgnore(ctx context.Context, user *entity.User) *entity.Error
}

// Signer responsibles for sign-in workflow.
type Signer struct {
	repo InsertUserRepository
}

// NewSigner creates an instance of Signer.
func NewSigner(repo InsertUserRepository) *Signer {
	return &Signer{
		repo: repo,
	}
}

// SignIn signs in a user to the system.
// If the user doesn't exist yet in the system, it will register the user.
func (s *Signer) SignIn(ctx context.Context, user *entity.User) *entity.Error {
	if user == nil {
		return entity.ErrEmptyUser
	}

	return s.repo.InsertOrIgnore(ctx, user)
}
