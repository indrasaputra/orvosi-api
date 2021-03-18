package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/indrasaputra/orvosi-api/entity"
)

// MedicalRecordSelector connects the database with medical record entity
// and only responsible for retrieving medical record data.
type MedicalRecordSelector struct {
	db *sql.DB
}

// NewMedicalRecordSelector creates an instance of MedicalRecordSelector.
func NewMedicalRecordSelector(db *sql.DB) *MedicalRecordSelector {
	return &MedicalRecordSelector{db: db}
}

// FindByID finds medical record by its id.
func (ms *MedicalRecordSelector) FindByID(ctx context.Context, id uint64) (*entity.MedicalRecord, *entity.Error) {
	query := "SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by, email FROM medical_records WHERE id = $1 LIMIT 1"
	row := ms.db.QueryRowContext(ctx, query, id)

	mr := &entity.MedicalRecord{
		User: &entity.User{},
	}
	if err := row.Scan(&mr.ID, &mr.Symptom, &mr.Diagnosis, &mr.Therapy, &mr.Result, &mr.CreatedAt, &mr.CreatedBy, &mr.UpdatedAt, &mr.UpdatedBy, &mr.User.Email); err != nil {
		return nil, entity.WrapError(entity.ErrInternalServer, err.Error())
	}
	return mr, nil
}

// FindByEmail finds all medical records bounded to specific email.
func (ms *MedicalRecordSelector) FindByEmail(ctx context.Context, email string, from uint64, limit uint) ([]*entity.MedicalRecord, *entity.Error) {
	query := "SELECT id, symptom, diagnosis, therapy, result, created_at, created_by, updated_at, updated_by FROM medical_records WHERE email = $1 AND id < $2 ORDER BY created_at DESC LIMIT $3"
	rows, err := ms.db.QueryContext(ctx, query, email, from, limit)
	if err != nil {
		return []*entity.MedicalRecord{}, entity.WrapError(entity.ErrInternalServer, err.Error())
	}
	defer rows.Close()

	var result []*entity.MedicalRecord
	for rows.Next() {
		var tmp entity.MedicalRecord
		if err := rows.Scan(&tmp.ID, &tmp.Symptom, &tmp.Diagnosis, &tmp.Therapy, &tmp.Result, &tmp.CreatedAt, &tmp.CreatedBy, &tmp.UpdatedAt, &tmp.UpdatedBy); err != nil {
			log.Printf("[MedicalRecordSelector-FindByEmail] scan rows error: %v", err)
			continue
		}

		result = append(result, &tmp)
	}
	if rows.Err() != nil {
		return []*entity.MedicalRecord{}, entity.WrapError(entity.ErrInternalServer, rows.Err().Error())
	}
	return result, nil
}
