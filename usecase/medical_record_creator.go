package usecase

import (
	"context"
	"strings"

	"github.com/indrasaputra/orvosi-api/entity"
)

// CreateMedicalRecord defines the business logic
// to create a medical record.
type CreateMedicalRecord interface {
	// Create creates a new medical record.
	Create(ctx context.Context, record *entity.MedicalRecord) *entity.Error
}

// InsertMedicalRecordRepository defines the business logic
// to insert a medical record into a repository.
type InsertMedicalRecordRepository interface {
	// Insert inserts the medical record into the repository.
	// This operation MUST set the inserted ID back to the medical record object.
	Insert(ctx context.Context, record *entity.MedicalRecord) *entity.Error
}

// MedicalRecordCreator responsibles for medical record creation workflow.
type MedicalRecordCreator struct {
	repo InsertMedicalRecordRepository
}

// NewMedicalRecordCreator creates an instance of MedicalRecordCreator.
func NewMedicalRecordCreator(repo InsertMedicalRecordRepository) *MedicalRecordCreator {
	return &MedicalRecordCreator{
		repo: repo,
	}
}

// Create creates a new medical record and persist it into a repository
func (mrc *MedicalRecordCreator) Create(ctx context.Context, record *entity.MedicalRecord) *entity.Error {
	if err := validateMedicalRecord(record); err != nil {
		return err
	}

	return mrc.repo.Insert(ctx, record)
}

func validateMedicalRecord(record *entity.MedicalRecord) *entity.Error {
	if record == nil {
		return entity.ErrEmptyMedicalRecord
	}

	sanitizeMedicalRecord(record)
	if !isMedicalRecordAttributesValid(record) {
		return entity.ErrInvalidMedicalRecordAttribute
	}
	return nil
}

func sanitizeMedicalRecord(record *entity.MedicalRecord) {
	record.Symptom = strings.TrimSpace(record.Symptom)
	record.Diagnosis = strings.TrimSpace(record.Diagnosis)
	record.Therapy = strings.TrimSpace(record.Therapy)
}

func isMedicalRecordAttributesValid(record *entity.MedicalRecord) bool {
	return record.Symptom != "" &&
		record.Diagnosis != "" &&
		record.Therapy != "" &&
		record.User != nil
}
