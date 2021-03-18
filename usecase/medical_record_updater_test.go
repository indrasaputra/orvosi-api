package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/orvosi-api/entity"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/indrasaputra/orvosi-api/usecase"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordUpdaterExecutor struct {
	usecase *usecase.MedicalRecordUpdater
	repo    *mock_usecase.MockUpdateMedicalRecordRepository
}

func TestNewMedicalRecordUpdater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of MedicalRecordUpdater", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor(ctrl)
		assert.NotNil(t, exec.usecase)
	})
}

func TestMedicalRecordUpater_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repository returns error", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor(ctrl)

		record := createValidMedicalRecord()
		exec.repo.EXPECT().DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com").Return(false, entity.ErrInternalServer)

		err := exec.usecase.Update(context.Background(), "dummy@dummy.com", uint64(1), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
	})

	t.Run("medical record not found", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor(ctrl)

		record := createValidMedicalRecord()
		exec.repo.EXPECT().DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com").Return(false, nil)

		err := exec.usecase.Update(context.Background(), "dummy@dummy.com", uint64(1), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrMedicalRecordNotFound, err)
	})

	t.Run("medical record can't be updated", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor(ctrl)

		record := createValidMedicalRecord()
		exec.repo.EXPECT().DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com").Return(true, nil)
		exec.repo.EXPECT().Update(context.Background(), uint64(1), record).Return(entity.ErrInternalServer)

		err := exec.usecase.Update(context.Background(), "dummy@dummy.com", uint64(1), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
	})

	t.Run("successfully update medical record", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor(ctrl)

		record := createValidMedicalRecord()
		exec.repo.EXPECT().DoesRecordExist(context.Background(), uint64(1), "dummy@dummy.com").Return(true, nil)
		exec.repo.EXPECT().Update(context.Background(), uint64(1), record).Return(nil)

		err := exec.usecase.Update(context.Background(), "dummy@dummy.com", uint64(1), record)

		assert.Nil(t, err)
	})
}

func createMedicalRecordUpdaterExecutor(ctrl *gomock.Controller) *MedicalRecordUpdaterExecutor {
	r := mock_usecase.NewMockUpdateMedicalRecordRepository(ctrl)
	u := usecase.NewMedicalRecordUpdater(r)

	return &MedicalRecordUpdaterExecutor{
		usecase: u,
		repo:    r,
	}
}
