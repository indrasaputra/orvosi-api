package repository_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/orvosi/api/entity"
	"github.com/orvosi/api/internal/repository"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordUpdater_Executor struct {
	repo *repository.MedicalRecordUpdater
	sql  sqlmock.Sqlmock
}

func TestNewMedicalRecordUpdater(t *testing.T) {
	t.Run("successfully create an instance of MedicalRecordUpdater", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor()
		assert.NotNil(t, exec.repo)
	})
}

func TestMedicalRecordUpdater_DoesRecordExist(t *testing.T) {
	t.Run("query returns internal error", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor()

		exec.sql.ExpectQuery(`SELECT id FROM medical_records WHERE id = \$1 AND email = \$2 LIMIT 1`).
			WillReturnError(errors.New("fail to select from database"))
		found, err := exec.repo.DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com")

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
		assert.False(t, found)
	})
}

func createMedicalRecordUpdaterExecutor() *MedicalRecordUpdater_Executor {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Panicf("error opening a stub database connection: %v\n", err)
	}

	repo := repository.NewMedicalRecordUpdater(db)
	return &MedicalRecordUpdater_Executor{
		repo: repo,
		sql:  mock,
	}
}
