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
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordUpdater_Executor struct {
	handler *handler.MedicalRecordUpdater
	usecase *mock_usecase.MockUpdateMedicalRecord
}

func TestNewMedicalRecordUpdater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of MedicalRecordUpdater", func(t *testing.T) {
		exec := createMedicalRecordUpdaterExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestMedicalRecordUpdater_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("entity id is not hashids.ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/", nil)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1234")

		exec := createMedicalRecordUpdaterExecutor(ctrl)
		exec.handler.Update(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-004","message":"Entity ID is invalid"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("can't process invalid medical record request", func(t *testing.T) {
		body, _ := json.Marshal("invalid request body")
		req := httptest.NewRequest(http.MethodPut, "/medical-records/:id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordUpdaterExecutor(ctrl)
		exec.handler.Update(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"02-003","message":"Medical record request is invalid. Please, check the JSON request"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("can't extract user information from request context", func(t *testing.T) {
		mr := createValidUpdateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordUpdaterExecutor(ctrl)
		exec.handler.Update(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("wanted record not found", func(t *testing.T) {
		mr := createValidUpdateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordUpdaterExecutor(ctrl)
		exec.usecase.EXPECT().Update(ctx.Request().Context(), user.Email, uint64(1), createMedicalRecordFromUpdateRequest(mr, user)).Return(entity.ErrMedicalRecordNotFound)
		exec.handler.Update(ctx)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"02-005","message":"Medical record not found"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("update usecase returns 5xx", func(t *testing.T) {
		mr := createValidUpdateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordUpdaterExecutor(ctrl)
		exec.usecase.EXPECT().Update(ctx.Request().Context(), user.Email, uint64(1), createMedicalRecordFromUpdateRequest(mr, user)).Return(entity.ErrInternalServer)
		exec.handler.Update(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("successfully update record", func(t *testing.T) {
		mr := createValidUpdateMedicalRecordRequest()
		body, _ := json.Marshal(mr)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordUpdaterExecutor(ctrl)
		exec.usecase.EXPECT().Update(ctx.Request().Context(), user.Email, uint64(1), createMedicalRecordFromUpdateRequest(mr, user)).Return(nil)
		exec.handler.Update(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		str := fmt.Sprintf("%s\n", `{"data":null,"meta":{}}`)
		assert.Equal(t, str, rec.Body.String())
	})
}

func createValidUpdateMedicalRecordRequest() *handler.UpdateMedicalRecordRequest {
	return &handler.UpdateMedicalRecordRequest{
		Symptom:   "symptom",
		Diagnosis: "diagnosis",
		Therapy:   "therapy",
		Result:    "result",
	}
}

func createMedicalRecordFromUpdateRequest(req *handler.UpdateMedicalRecordRequest, user *entity.User) *entity.MedicalRecord {
	return &entity.MedicalRecord{
		User:      user,
		Symptom:   req.Symptom,
		Diagnosis: req.Diagnosis,
		Therapy:   req.Therapy,
		Result:    req.Result,
	}
}

func createMedicalRecordUpdaterExecutor(ctrl *gomock.Controller) *MedicalRecordUpdater_Executor {
	u := mock_usecase.NewMockUpdateMedicalRecord(ctrl)
	h := handler.NewMedicalRecordUpdater(u)
	return &MedicalRecordUpdater_Executor{
		handler: h,
		usecase: u,
	}
}
