package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/indrasaputra/orvosi-api/entity"
)

// UserInserter connects the database with user entity
// and only responsible for inserting a new data.
type UserInserter struct {
	db *sql.DB
}

// NewUserInserter creates an instance of UserInserter.
func NewUserInserter(db *sql.DB) *UserInserter {
	return &UserInserter{db: db}
}

// InsertOrIgnore inserts a new data into the database.
// If the user already exists in the database, it does nothing.
func (ui *UserInserter) InsertOrIgnore(ctx context.Context, user *entity.User) *entity.Error {
	if user == nil {
		return entity.ErrEmptyUser
	}

	query := "INSERT INTO users (name, email, google_id, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (email) DO NOTHING"
	_, err := ui.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.GoogleID,
		time.Now(),
		time.Now(),
		user.Email,
		user.Email,
	)

	if err != nil {
		return entity.WrapError(entity.ErrInternalServer, err.Error())
	}
	return nil
}
