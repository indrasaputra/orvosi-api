package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/repository"
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
		assert.Equal(t, entity.ErrInternalServer.Code, err.Code)
		assert.False(t, found)
	})

	t.Run("record not found in repository", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor()

		exec.sql.ExpectQuery(`SELECT id FROM medical_records WHERE id = \$1 AND email = \$2 LIMIT 1`).
			WillReturnError(sql.ErrNoRows)
		found, err := exec.repo.DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com")

		assert.Nil(t, err)
		assert.False(t, found)
	})

	t.Run("successfully found the record", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor()

		exec.sql.ExpectQuery(`SELECT id FROM medical_records WHERE id = \$1 AND email = \$2 LIMIT 1`).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(999),
			)
		found, err := exec.repo.DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com")

		assert.Nil(t, err)
		assert.True(t, found)
	})
}

func TestMedicalRecordUpdater_Update(t *testing.T) {
	t.Run("query returns internal error", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor()

		exec.sql.ExpectExec(`UPDATE medical_records SET symptom = \$1, diagnosis = \$2, therapy = \$3, result = \$4, updated_at = \$5, updated_by = \$6 WHERE id = \$7`).
			WillReturnError(errors.New("fail to select from database"))
		err := exec.repo.Update(context.Background(), uint64(1), createValidMedicalRecord())

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer.Code, err.Code)
	})

	t.Run("successfully update the medical record", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor()

		exec.sql.ExpectExec(`UPDATE medical_records SET symptom = \$1, diagnosis = \$2, therapy = \$3, result = \$4, updated_at = \$5, updated_by = \$6 WHERE id = \$7`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := exec.repo.Update(context.Background(), uint64(1), createValidMedicalRecord())

		assert.Nil(t, err)
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
