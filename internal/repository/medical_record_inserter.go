package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
)

// MedicalRecordInserter connects the database with medical record entity
// and only responsible for inserting a new data.
type MedicalRecordInserter struct {
	db *sql.DB
}

// NewMedicalRecordInserter creates an instance of MedicalRecordInserter.
func NewMedicalRecordInserter(db *sql.DB) *MedicalRecordInserter {
	return &MedicalRecordInserter{db: db}
}

// Insert inserts a new medical record data into the database.
func (mri *MedicalRecordInserter) Insert(ctx context.Context, record *entity.MedicalRecord) *entity.Error {
	if record == nil {
		return entity.ErrEmptyMedicalRecord
	}

	query := "INSERT INTO " +
		"medical_records (email, symptom, diagnosis, therapy, result, created_at, updated_at, created_by, updated_by) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	row := mri.db.QueryRow(query,
		record.User.Email,
		record.Symptom,
		record.Diagnosis,
		record.Therapy,
		"",
		time.Now(),
		time.Now(),
		record.User.Email,
		record.User.Email,
	)

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return entity.WrapError(entity.ErrInternalServer, "[MedicalRecordInserter-Insert] exec insert query: "+err.Error())
	}

	record.ID = hashids.ID(id)
	return nil
}
