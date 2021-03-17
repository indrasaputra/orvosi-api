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

type MedicalRecordFinder_Executor struct {
	usecase *usecase.MedicalRecordFinder
	repo    *mock_usecase.MockFindMedicalRecordRepository
}

func TestNewMedicalRecordFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of MedicalRecordFinder", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)
		assert.NotNil(t, exec.usecase)
	})
}

func TestMedicalRecordFinder_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("repo returns error", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		exec.repo.EXPECT().FindByID(context.Background(), uint64(1)).Return(&entity.MedicalRecord{}, entity.ErrInternalServer)
		res, err := exec.usecase.FindByID(context.Background(), uint64(1), "dummy@dummy.com")

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
		assert.Nil(t, res)
	})

	t.Run("medical record is not owned by the user", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		exec.repo.EXPECT().FindByID(context.Background(), uint64(1)).Return(&entity.MedicalRecord{User: &entity.User{Email: "notdummy@dummy.com"}}, nil)
		res, err := exec.usecase.FindByID(context.Background(), uint64(1), "dummy@dummy.com")

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrUnauthorized, err)
		assert.Nil(t, res)
	})

	t.Run("successfully find medical records bounded to specific email", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		exec.repo.EXPECT().FindByID(context.Background(), uint64(1)).Return(&entity.MedicalRecord{User: &entity.User{Email: "dummy@dummy.com"}}, nil)
		res, err := exec.usecase.FindByID(context.Background(), uint64(1), "dummy@dummy.com")

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestMedicalRecordFinder_FindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("email doesn't contain username", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		res, err := exec.usecase.FindByEmail(context.Background(), "@dummy.com", 0)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidEmail, err)
		assert.Empty(t, res)
	})

	t.Run("email doesn't contain domain", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		res, err := exec.usecase.FindByEmail(context.Background(), "dummy@", 0)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidEmail, err)
		assert.Empty(t, res)
	})

	t.Run("email contains made-up domain", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		res, err := exec.usecase.FindByEmail(context.Background(), "dummy@dummy-domain.com", 0)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidEmail, err)
		assert.Empty(t, res)
	})

	t.Run("repo returns error", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		exec.repo.EXPECT().FindByEmail(context.Background(), "dummy@dummy.com", uint64(0), uint(10)).Return([]*entity.MedicalRecord{}, entity.ErrInternalServer)
		res, err := exec.usecase.FindByEmail(context.Background(), "dummy@dummy.com", 0)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInternalServer, err)
		assert.Empty(t, res)
	})

	t.Run("successfully find medical records bounded to specific email", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)

		exec.repo.EXPECT().FindByEmail(context.Background(), "dummy@dummy.com", uint64(0), uint(10)).Return([]*entity.MedicalRecord{&entity.MedicalRecord{}}, nil)
		res, err := exec.usecase.FindByEmail(context.Background(), "dummy@dummy.com", 0)

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, len(res))
	})
}

func createMedicalRecordFinderExecutor(ctrl *gomock.Controller) *MedicalRecordFinder_Executor {
	r := mock_usecase.NewMockFindMedicalRecordRepository(ctrl)
	u := usecase.NewMedicalRecordFinder(r)

	return &MedicalRecordFinder_Executor{
		usecase: u,
		repo:    r,
	}
}
