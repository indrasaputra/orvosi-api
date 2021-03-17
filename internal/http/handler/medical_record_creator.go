package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	"github.com/indrasaputra/orvosi-api/internal/http/response"
	"github.com/indrasaputra/orvosi-api/usecase"
)

// CreateMedicalRecordRequest represents medical record request.
type CreateMedicalRecordRequest struct {
	Symptom   string `json:"symptom"`
	Diagnosis string `json:"diagnosis"`
	Therapy   string `json:"therapy"`
}

// MedicalRecordCreator handles HTTP request and response
// for create medical record.
type MedicalRecordCreator struct {
	creator usecase.CreateMedicalRecord
}

// NewMedicalRecordCreator creates an instance of MedicalCreator.
func NewMedicalRecordCreator(creator usecase.CreateMedicalRecord) *MedicalRecordCreator {
	return &MedicalRecordCreator{
		creator: creator,
	}
}

// Create handles `POST /medical-records` endpoint.
func (mrc *MedicalRecordCreator) Create(ctx echo.Context) error {
	var request CreateMedicalRecordRequest
	if err := ctx.Bind(&request); err != nil {
		res := response.NewError(entity.ErrInvalidMedicalRecordRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return err
	}

	user, err := extractUserFromRequestContext(ctx.Request().Context())
	if err != nil {
		res := response.NewError(err)
		ctx.JSON(http.StatusInternalServerError, res)
		return err
	}

	record := createMedicalRecordFromCreateRequest(&request, user)
	if err := mrc.creator.Create(ctx.Request().Context(), record); err != nil {
		res := response.NewError(err)
		status := http.StatusInternalServerError
		if err.Code != entity.ErrInternalServer.Code {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, res)
		return err
	}

	ctx.JSON(http.StatusCreated, response.NewSuccess(nil, response.EmptyMeta{}))
	return nil
}

func extractUserFromRequestContext(ctx context.Context) (*entity.User, error) {
	val := ctx.Value(middleware.ContextKeyUser)
	user, ok := val.(*entity.User)
	if !ok {
		return nil, entity.ErrInternalServer
	}
	return user, nil
}

func createMedicalRecordFromCreateRequest(req *CreateMedicalRecordRequest, user *entity.User) *entity.MedicalRecord {
	return &entity.MedicalRecord{
		User:      user,
		Symptom:   req.Symptom,
		Diagnosis: req.Diagnosis,
		Therapy:   req.Therapy,
	}
}
