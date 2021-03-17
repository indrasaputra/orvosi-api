package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordCreator_Executor struct {
	handler *handler.MedicalRecordCreator
	usecase *mock_usecase.MockCreateMedicalRecord
}

func TestNewMedicalRecordCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of MedicalRecordCreator", func(t *testing.T) {
		exec := createMedicalRecordCreatorExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestMedicalRecordCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("can't process invalid medical record request", func(t *testing.T) {
		body, _ := json.Marshal("invalid request body")
		req := httptest.NewRequest(http.MethodPost, "/medical-records", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordCreatorExecutor(ctrl)
		exec.handler.Create(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"02-003","message":"Medical record request is invalid. Please, check the JSON request"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("can't extract user information from request context", func(t *testing.T) {
		mr := createValidCreateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPost, "/medical-records", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordCreatorExecutor(ctrl)
		exec.handler.Create(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("creator service returns 4xx error", func(t *testing.T) {
		mr := createValidCreateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPost, "/medical-records", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordCreatorExecutor(ctrl)
		exec.usecase.EXPECT().Create(ctx.Request().Context(), createMedicalRecordFromCreateRequest(mr, user)).Return(entity.ErrInvalidMedicalRecordAttribute)
		exec.handler.Create(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"02-002","message":"Medical record's attributes are invalid. Please, check all attributes"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("creator service returns 5xx error", func(t *testing.T) {
		mr := createValidCreateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPost, "/medical-records", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordCreatorExecutor(ctrl)
		exec.usecase.EXPECT().Create(ctx.Request().Context(), createMedicalRecordFromCreateRequest(mr, user)).Return(entity.ErrInternalServer)
		exec.handler.Create(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("creator service successfully process the request", func(t *testing.T) {
		mr := createValidCreateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPost, "/medical-records", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordCreatorExecutor(ctrl)
		exec.usecase.EXPECT().Create(ctx.Request().Context(), createMedicalRecordFromCreateRequest(mr, user)).Return(nil)
		exec.handler.Create(ctx)

		assert.Equal(t, http.StatusCreated, rec.Code)
		str := fmt.Sprintf("%s\n", `{"data":null,"meta":{}}`)
		assert.Equal(t, str, rec.Body.String())
	})
}

func createValidCreateMedicalRecordRequest() *handler.CreateMedicalRecordRequest {
	return &handler.CreateMedicalRecordRequest{
		Symptom:   "symptom",
		Diagnosis: "diagnosis",
		Therapy:   "therapy",
	}
}

func createUserInformation() *entity.User {
	return &entity.User{
		ID:       hashids.ID(1),
		Email:    "user@email.com",
		Name:     "Edward Jenner",
		GoogleID: "12345678901234567890",
	}
}

func createMedicalRecordFromCreateRequest(req *handler.CreateMedicalRecordRequest, user *entity.User) *entity.MedicalRecord {
	return &entity.MedicalRecord{
		User:      user,
		Symptom:   req.Symptom,
		Diagnosis: req.Diagnosis,
		Therapy:   req.Therapy,
	}
}

func createMedicalRecordCreatorExecutor(ctrl *gomock.Controller) *MedicalRecordCreator_Executor {
	u := mock_usecase.NewMockCreateMedicalRecord(ctrl)
	h := handler.NewMedicalRecordCreator(u)
	return &MedicalRecordCreator_Executor{
		handler: h,
		usecase: u,
	}
}
