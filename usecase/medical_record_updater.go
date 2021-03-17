package usecase

import (
	"context"

	"github.com/indrasaputra/orvosi-api/entity"
)

// UpdateMedicalRecord defines the business logic
// to update a medical record.
type UpdateMedicalRecord interface {
	// Update updates a new medical record.
	Update(ctx context.Context, email string, id uint64, record *entity.MedicalRecord) *entity.Error
}

// UpdateMedicalRecordRepository defines the business logic
// to update a medical record into a repository.
type UpdateMedicalRecordRepository interface {
	// FindByID finds medical records by its id.
	DoesRecordExist(ctx context.Context, id uint64, email string) (bool, *entity.Error)
	// Update updates certain medical record.
	Update(ctx context.Context, id uint64, record *entity.MedicalRecord) *entity.Error
}

// MedicalRecordUpdater responsibles for medical record update workflow.
type MedicalRecordUpdater struct {
	repo UpdateMedicalRecordRepository
}

// NewMedicalRecordUpdater creates an instance of MedicalRecordUpdater.
func NewMedicalRecordUpdater(repo UpdateMedicalRecordRepository) *MedicalRecordUpdater {
	return &MedicalRecordUpdater{
		repo: repo,
	}
}

// Update updates the medical record.
func (mu *MedicalRecordUpdater) Update(ctx context.Context, email string, id uint64, record *entity.MedicalRecord) *entity.Error {
	found, err := mu.repo.DoesRecordExist(ctx, id, email)
	if err != nil {
		return err
	}
	if !found {
		return entity.ErrMedicalRecordNotFound
	}

	return mu.repo.Update(ctx, id, record)
}
