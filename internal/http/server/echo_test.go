package server_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	"github.com/indrasaputra/orvosi-api/internal/http/router"
	"github.com/indrasaputra/orvosi-api/internal/http/server"
	"github.com/indrasaputra/orvosi-api/internal/tool"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("successfully create an instance of Server", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		srv := createServer(ctrl)
		assert.NotNil(t, srv)
	})
}

func createMedicalRecordCreator(ctrl *gomock.Controller) *handler.MedicalRecordCreator {
	m := mock_usecase.NewMockCreateMedicalRecord(ctrl)
	return handler.NewMedicalRecordCreator(m)
}

func createServer(ctrl *gomock.Controller) *server.Server {
	c := createMedicalRecordCreator(ctrl)
	r := router.MedicalRecordCreator(c)
	d := tool.NewIDTokenDecoder("audience")
	m := middleware.WithJWTDecoder(d.Decode)
	return server.NewServer(m, r)
}
