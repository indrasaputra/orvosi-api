package repository_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordSelectorExecutor struct {
	repo *repository.MedicalRecordSelector
	sql  sqlmock.Sqlmock
}

func TestNewMedicalRecordSelector(t *testing.T) {
	t.Run("successfully create an instance of MedicalRecordSelector", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()
		assert.NotNil(t, exec.repo)
	})
}

func TestMedicalRecordSelector_FindByID(t *testing.T) {
	t.Run("select query returns error", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by, email FROM medical_records WHERE id = \$1 LIMIT 1`).
			WillReturnError(errors.New("fail to select from database"))

		res, err := exec.repo.FindByID(context.Background(), uint64(1))

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer.Code, err.Code)
		assert.Empty(t, res)
	})

	t.Run("row scan returns error", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by, email FROM medical_records WHERE id = \$1 LIMIT 1`).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "symptom", "diagnosis", "therapy", "result", "created_at", "created_by", "updated_at", "updated_by", "email"}).
				AddRow(1, "Symptom", "Diagnosis", "Therapy", "Result", "time.Now()", "dummy@dummy.com", "time.Now()", "dummy@dummy.com", "dummy@dummy.com"),
			)

		res, err := exec.repo.FindByID(context.Background(), uint64(1))

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("successfully retrieve one medical record", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by, email FROM medical_records WHERE id = \$1 LIMIT 1`).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "symptom", "diagnosis", "therapy", "result", "created_at", "created_by", "updated_at", "updated_by", "email"}).
				AddRow(1, "Symptom", "Diagnosis", "Therapy", "Result", time.Now(), "dummy@dummy.com", time.Now(), "dummy@dummy.com", "dummy@dummy.com"),
			)

		res, err := exec.repo.FindByID(context.Background(), uint64(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uint64(1), uint64(res.ID))
	})
}

func TestMedicalRecordSelector_FindByEmail(t *testing.T) {
	t.Run("select query returns error", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by FROM medical_records WHERE email = \$1 AND id < \$2 ORDER BY created_at DESC LIMIT \$3`).
			WillReturnError(errors.New("fail to select from database"))

		res, err := exec.repo.FindByEmail(context.Background(), "dummy@dummy.com", 100, 10)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer.Code, err.Code)
		assert.Empty(t, res)
	})

	t.Run("row scan returns error", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by FROM medical_records WHERE email = \$1 AND id < \$2 ORDER BY created_at DESC LIMIT \$3`).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "symptom", "diagnosis", "therapy", "result", "created_at", "created_by", "updated_at", "updated_by"}).
				AddRow(1, "Symptom", "Diagnosis", "Therapy", "Result", time.Now(), "dummy@dummy.com", time.Now(), "dummy@dummy.com").
				AddRow(2, "Symptom", "Diagnosis", "Therapy", "Result", "time.Now()", "dummy@dummy.com", "time.Now()", "dummy@dummy.com"),
			)

		res, err := exec.repo.FindByEmail(context.Background(), "dummy@dummy.com", 100, 10)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res))
	})

	t.Run("rows error occurs after scanning", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by FROM medical_records WHERE email = \$1 AND id < \$2 ORDER BY created_at DESC LIMIT \$3`).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "symptom", "diagnosis", "therapy", "result", "created_at", "created_by", "updated_at", "updated_by"}).
				AddRow(1, "Symptom", "Diagnosis", "Therapy", "Result", time.Now(), "dummy@dummy.com", time.Now(), "dummy@dummy.com").
				AddRow(2, "Symptom", "Diagnosis", "Therapy", "Result", "time.Now()", "dummy@dummy.com", "time.Now()", "dummy@dummy.com").
				RowError(1, errors.New("rows error")),
			)

		res, err := exec.repo.FindByEmail(context.Background(), "dummy@dummy.com", 100, 10)

		assert.NotNil(t, err)
		assert.Empty(t, res)
	})

	t.Run("successfully retrieve all rows", func(t *testing.T) {
		exec := createMedicalRecordSelectorExecutor()

		exec.sql.ExpectQuery(`SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by FROM medical_records WHERE email = \$1 AND id < \$2 ORDER BY created_at DESC LIMIT \$3`).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "symptom", "diagnosis", "therapy", "result", "created_at", "created_by", "updated_at", "updated_by"}).
				AddRow(1, "Symptom", "Diagnosis", "Therapy", "Result", time.Now(), "dummy@dummy.com", time.Now(), "dummy@dummy.com").
				AddRow(2, "Symptom", "Diagnosis", "Therapy", "Result", time.Now(), "dummy@dummy.com", time.Now(), "dummy@dummy.com"),
			)

		res, err := exec.repo.FindByEmail(context.Background(), "dummy@dummy.com", 100, 10)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 2, len(res))
	})
}

func createMedicalRecordSelectorExecutor() *MedicalRecordSelectorExecutor {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Panicf("[createMedicalRecordSelectorExecutor] error opening a stub database connection: %v\n", err)
	}

	repo := repository.NewMedicalRecordSelector(db)
	return &MedicalRecordSelectorExecutor{
		repo: repo,
		sql:  mock,
	}
}
