package router_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/router"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/stretchr/testify/assert"
)

func TestMedicalRecordCreatorRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("all desired medical record creator routes are registered", func(t *testing.T) {
		desired := map[string]string{
			"/medical-records": "POST",
		}

		h := createMedicalRecordCreator(ctrl)
		routes := router.MedicalRecordCreator(h)

		for _, route := range routes {
			assert.Equal(t, desired[route.Path], route.Method)
			assert.NotEmpty(t, route.Middlewares)
		}
	})
}

func TestMedicalRecordFinderRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("all desired medical record finder routes are registered", func(t *testing.T) {
		desired := map[string]string{
			"/medical-records":     "GET",
			"/medical-records/:id": "GET",
		}

		h := createMedicalRecordFinder(ctrl)
		routes := router.MedicalRecordFinder(h)

		for _, route := range routes {
			assert.Equal(t, desired[route.Path], route.Method)
			assert.Empty(t, route.Middlewares)
		}
	})
}

func TestMedicalRecordUpdaterRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("all desired medical record updater routes are registered", func(t *testing.T) {
		desired := map[string]string{
			"/medical-records/:id": "PUT",
		}

		h := createMedicalRecordUpdater(ctrl)
		routes := router.MedicalRecordUpdater(h)

		for _, route := range routes {
			assert.Equal(t, desired[route.Path], route.Method)
			assert.NotEmpty(t, route.Middlewares)
		}
	})
}

func createMedicalRecordCreator(ctrl *gomock.Controller) *handler.MedicalRecordCreator {
	m := mock_usecase.NewMockCreateMedicalRecord(ctrl)
	return handler.NewMedicalRecordCreator(m)
}

func createMedicalRecordFinder(ctrl *gomock.Controller) *handler.MedicalRecordFinder {
	m := mock_usecase.NewMockFindMedicalRecord(ctrl)
	return handler.NewMedicalRecordFinder(m)
}

func createMedicalRecordUpdater(ctrl *gomock.Controller) *handler.MedicalRecordUpdater {
	m := mock_usecase.NewMockUpdateMedicalRecord(ctrl)
	return handler.NewMedicalRecordUpdater(m)
}
