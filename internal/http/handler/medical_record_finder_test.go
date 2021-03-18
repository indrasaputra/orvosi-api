package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/hashids"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MedicalRecordFinderExecutor struct {
	handler *handler.MedicalRecordFinder
	usecase *mock_usecase.MockFindMedicalRecord
}

const (
	maxUint64 = 1<<32 - 1
)

func TestNewMedicalRecordFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of MedicalRecordFinder", func(t *testing.T) {
		exec := createMedicalRecordFinderExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestMedicalRecordFinder_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("entity id is not hashids.ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1234")

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.handler.FindByID(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-004","message":"Entity ID is invalid"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("can't extract user information from request context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.handler.FindByID(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("medical record is not owned by given email", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.usecase.EXPECT().FindByID(ctx.Request().Context(), uint64(1), user.Email).Return(nil, entity.ErrUnauthorized)
		exec.handler.FindByID(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-002","message":"Request is unauthorized"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("finder service returns 5xx", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.usecase.EXPECT().FindByID(ctx.Request().Context(), uint64(1), user.Email).Return(nil, entity.ErrInternalServer)
		exec.handler.FindByID(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("successfully get medical record owned by certain email", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/medical-records/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("oWx0b8DZ1a")

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.usecase.EXPECT().FindByID(ctx.Request().Context(), uint64(1), user.Email).Return(createMedicalRecords()[0], nil)
		exec.handler.FindByID(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		str := fmt.Sprintf("%s\n", `{"data":{"id":"oWx0b8DZ1a","symptom":"Symptom","diagnosis":"Diagnosis","therapy":"Therapy","result":"Result","created_by":"user@dummy.com","created_at":"2021-01-28T15:00:00Z","updated_by":"user@dummy.com","updated_at":"2021-01-28T15:00:00Z"},"meta":{}}`)
		assert.Equal(t, str, rec.Body.String())
	})
}

func TestMedicalRecordFinder_FindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("can't extract user information from request context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/medical-records", nil)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.handler.FindByEmail(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("query param 'from' is invalid (not hashids)", func(t *testing.T) {
		paths := []string{"/medical-records?from=abc", "/medical-records?from=-10", "/medical-records?from=-12345"}
		for _, path := range paths {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			user := createUserInformation()
			req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			exec := createMedicalRecordFinderExecutor(ctrl)
			exec.handler.FindByEmail(ctx)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			str := fmt.Sprintf("%s\n", `{"errors":[{"code":"02-006","message":"Query param(s) is invalid"}],"meta":null}`)
			assert.Equal(t, str, rec.Body.String())
		}
	})

	t.Run("finder service returns 4xx error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/medical-records?from=oWx0b8DZ1a", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.usecase.EXPECT().FindByEmail(ctx.Request().Context(), user.Email, uint64(1)).Return([]*entity.MedicalRecord{}, entity.ErrInvalidEmail)
		exec.handler.FindByEmail(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"02-004","message":"Email is invalid. Please, check the email"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("finder service returns 5xx error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/medical-records?from=oWx0b8DZ1a", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createMedicalRecordFinderExecutor(ctrl)
		exec.usecase.EXPECT().FindByEmail(ctx.Request().Context(), user.Email, uint64(1)).Return([]*entity.MedicalRecord{}, entity.ErrInternalServer)
		exec.handler.FindByEmail(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("successfully get medical records from certain user (email)", func(t *testing.T) {
		tables := []struct {
			path  string
			param uint64
		}{
			{"/medical-records?from=oWx0b8DZ1a", 1},
			{"/medical-records", maxUint64},
			{"/medical-records?from=", maxUint64},
		}

		for _, table := range tables {
			req := httptest.NewRequest(http.MethodGet, table.path, nil)
			user := createUserInformation()
			req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			exec := createMedicalRecordFinderExecutor(ctrl)
			mrs := createMedicalRecords()
			exec.usecase.EXPECT().FindByEmail(ctx.Request().Context(), user.Email, table.param).Return(mrs, nil)
			exec.handler.FindByEmail(ctx)

			assert.Equal(t, http.StatusOK, rec.Code)
			str := fmt.Sprintf("%s\n", `{"data":[{"id":"oWx0b8DZ1a","symptom":"Symptom","diagnosis":"Diagnosis","therapy":"Therapy","result":"Result","created_by":"user@dummy.com","created_at":"2021-01-28T15:00:00Z","updated_by":"user@dummy.com","updated_at":"2021-01-28T15:00:00Z"}],"meta":{}}`)
			assert.Equal(t, str, rec.Body.String())
		}
	})
}

func createMedicalRecords() []*entity.MedicalRecord {
	return []*entity.MedicalRecord{
		{
			ID: hashids.ID(1),
			User: &entity.User{
				ID:       hashids.ID(1),
				Email:    "user@dummy.com",
				Name:     "User Dummy",
				GoogleID: "1234567890",
			},
			Symptom:   "Symptom",
			Diagnosis: "Diagnosis",
			Therapy:   "Therapy",
			Result:    "Result",
			Auditable: entity.Auditable{
				CreatedBy: "user@dummy.com",
				CreatedAt: time.Date(2021, time.January, 28, 15, 00, 00, 00, time.UTC),
				UpdatedBy: "user@dummy.com",
				UpdatedAt: time.Date(2021, time.January, 28, 15, 00, 00, 00, time.UTC),
			},
		},
	}
}

func createMedicalRecordFinderExecutor(ctrl *gomock.Controller) *MedicalRecordFinderExecutor {
	u := mock_usecase.NewMockFindMedicalRecord(ctrl)
	h := handler.NewMedicalRecordFinder(u)
	return &MedicalRecordFinderExecutor{
		handler: h,
		usecase: u,
	}
}
