package handler

import (
	"net/http"

	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/response"
	"github.com/indrasaputra/orvosi-api/usecase"
	"github.com/labstack/echo/v4"
)

// UpdateMedicalRecordRequest represents medical record request.
type UpdateMedicalRecordRequest struct {
	Symptom   string `json:"symptom"`
	Diagnosis string `json:"diagnosis"`
	Therapy   string `json:"therapy"`
	Result    string `json:"result"`
}

// MedicalRecordUpdater handles HTTP request and response
// for update medical record.
type MedicalRecordUpdater struct {
	updater usecase.UpdateMedicalRecord
}

// NewMedicalRecordUpdater updates an instance of MedicalUpdater.
func NewMedicalRecordUpdater(updater usecase.UpdateMedicalRecord) *MedicalRecordUpdater {
	return &MedicalRecordUpdater{
		updater: updater,
	}
}

// Update handles `PUT /medical-records/:id` endpoint.
func (mru *MedicalRecordUpdater) Update(ctx echo.Context) error {
	str := ctx.Param("id")
	id, herr := hashids.DecodeHash([]byte(str))
	if herr != nil {
		res := response.NewError(entity.ErrInvalidID)
		ctx.JSON(http.StatusInternalServerError, res)
		return herr
	}

	var request UpdateMedicalRecordRequest
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

	record := createMedicalRecordFromUpdateRequest(&request, user)
	if err := mru.updater.Update(ctx.Request().Context(), user.Email, uint64(id), record); err != nil {
		res := response.NewError(err)
		status := http.StatusNotFound
		if err.Code == entity.ErrInternalServer.Code {
			status = http.StatusInternalServerError
		}
		ctx.JSON(status, res)
		return err
	}

	ctx.JSON(http.StatusOK, response.NewSuccess(nil, response.EmptyMeta{}))
	return nil
}

func createMedicalRecordFromUpdateRequest(req *UpdateMedicalRecordRequest, user *entity.User) *entity.MedicalRecord {
	return &entity.MedicalRecord{
		User:      user,
		Symptom:   req.Symptom,
		Diagnosis: req.Diagnosis,
		Therapy:   req.Therapy,
		Result:    req.Result,
	}
}
