package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/indrasaputra/orvosi-api/usecase"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordCreator_Executor struct {
	usecase *usecase.MedicalRecordCreator
	repo    *mock_usecase.MockInsertMedicalRecordRepository
}

func TestNewMedicalRecordCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of MedicalRecordCreator", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		assert.NotNil(t, exec.usecase)
	})
}

func TestMedicalRecordCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("medical record entity is empty/nil", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)

		err := exec.usecase.Create(context.Background(), nil)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrEmptyMedicalRecord, err)
	})

	t.Run("medical record's symptom is empty", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		record := createValidMedicalRecord()
		record.Symptom = "   "

		err := exec.usecase.Create(context.Background(), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidMedicalRecordAttribute, err)
	})

	t.Run("medical record's diagnosis is empty", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		record := createValidMedicalRecord()
		record.Diagnosis = "   "

		err := exec.usecase.Create(context.Background(), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidMedicalRecordAttribute, err)
	})

	t.Run("medical record's therapy is empty", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		record := createValidMedicalRecord()
		record.Therapy = ""

		err := exec.usecase.Create(context.Background(), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidMedicalRecordAttribute, err)
	})

	t.Run("medical record's user is nil", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		record := createValidMedicalRecord()
		record.User = nil

		err := exec.usecase.Create(context.Background(), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidMedicalRecordAttribute, err)
	})

	t.Run("medical record repo fails", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		record := createValidMedicalRecord()

		exec.repo.EXPECT().Insert(context.Background(), record).Return(entity.ErrInternalServer)

		err := exec.usecase.Create(context.Background(), record)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
	})

	t.Run("successfully create a new medical record", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		record := createValidMedicalRecord()

		exec.repo.EXPECT().Insert(context.Background(), record).Return(nil)

		err := exec.usecase.Create(context.Background(), record)

		assert.Nil(t, err)
	})
}

func createValidMedicalRecord() *entity.MedicalRecord {
	return &entity.MedicalRecord{
		ID:        hashids.ID(1),
		Symptom:   "symptom",
		Diagnosis: "diagnosis",
		Therapy:   "therapy",
		Result:    "result",
		User: &entity.User{
			ID:       hashids.ID(1),
			Email:    "email@provider.com",
			Name:     "User 1",
			GoogleID: "super-long-google-id",
		},
	}
}

func createMedicalRecordCreatorExecutor(ctrl *gomock.Controller) *MedicalRecordCreator_Executor {
	r := mock_usecase.NewMockInsertMedicalRecordRepository(ctrl)
	u := usecase.NewMedicalRecordCreator(r)

	return &MedicalRecordCreator_Executor{
		usecase: u,
		repo:    r,
	}
}
