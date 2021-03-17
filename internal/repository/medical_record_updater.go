package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/indrasaputra/orvosi-api/entity"
)

// MedicalRecordUpdater connects the database with medical record entity
// and only responsible for usecase of updating a data.
type MedicalRecordUpdater struct {
	db *sql.DB
}

// NewMedicalRecordUpdater creates an instance of MedicalRecordUpdater.
func NewMedicalRecordUpdater(db *sql.DB) *MedicalRecordUpdater {
	return &MedicalRecordUpdater{db: db}
}

// DoesRecordExist checks whether medical record which has certain id and email exists.
func (mu *MedicalRecordUpdater) DoesRecordExist(ctx context.Context, id uint64, email string) (bool, *entity.Error) {
	query := "SELECT id FROM medical_records WHERE id = $1 AND email = $2 LIMIT 1"
	row := mu.db.QueryRowContext(ctx, query, id, email)

	var tmp uint64
	err := row.Scan(&tmp)
	if err != nil && err != sql.ErrNoRows {
		return false, entity.WrapError(entity.ErrInternalServer, err.Error())
	} else if err != nil && err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

// Update updates the whole record data.
func (mu *MedicalRecordUpdater) Update(ctx context.Context, id uint64, record *entity.MedicalRecord) *entity.Error {
	query := "UPDATE medical_records SET symptom = $1, diagnosis = $2, therapy = $3, result = $4, updated_at = $5, updated_by = $6 WHERE id = $7"
	_, err := mu.db.ExecContext(ctx, query,
		record.Symptom,
		record.Diagnosis,
		record.Therapy,
		record.Result,
		time.Now(),
		record.User.Email,
		id,
	)

	if err != nil {
		return entity.WrapError(entity.ErrInternalServer, err.Error())
	}
	return nil
}
