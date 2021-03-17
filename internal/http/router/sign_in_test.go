package router_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/router"
	mock_usecase "github.com/indrasaputra/orvosi-api/test/mock/usecase"
	"github.com/stretchr/testify/assert"
)

func TestSignerRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("all desired signer routes are registered", func(t *testing.T) {
		desired := map[string]string{
			"/sign-in": "POST",
		}

		h := createSigner(ctrl)
		routes := router.Signer(h)

		for _, route := range routes {
			assert.Equal(t, desired[route.Path], route.Method)
		}
	})
}

func createSigner(ctrl *gomock.Controller) *handler.Signer {
	m := mock_usecase.NewMockSignIn(ctrl)
	return handler.NewSigner(m)
}
