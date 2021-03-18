package handler_test

import (
	"context"
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

type SignerExecutor struct {
	handler *handler.Signer
	usecase *mock_usecase.MockSignIn
}

func TestNewSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Signer", func(t *testing.T) {
		exec := createSignerExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestSigner_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("can't extract user information from request context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createSignerExecutor(ctrl)
		exec.handler.SignIn(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("somehow signer service doesn't receive user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createSignerExecutor(ctrl)
		exec.usecase.EXPECT().SignIn(ctx.Request().Context(), user).Return(entity.ErrEmptyUser)
		exec.handler.SignIn(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"03-001","message":"User is empty"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("signer service returns internal error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createSignerExecutor(ctrl)
		exec.usecase.EXPECT().SignIn(ctx.Request().Context(), user).Return(entity.ErrInternalServer)
		exec.handler.SignIn(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := fmt.Sprintf("%s\n", `{"errors":[{"code":"01-001","message":"Internal server error"}],"meta":null}`)
		assert.Equal(t, str, rec.Body.String())
	})

	t.Run("successfully sign in", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		user := createUserInformation()
		req = req.WithContext(context.WithValue(context.Background(), middleware.ContextKeyUser, user))

		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		exec := createSignerExecutor(ctrl)
		exec.usecase.EXPECT().SignIn(ctx.Request().Context(), user).Return(nil)
		exec.handler.SignIn(ctx)

		assert.Equal(t, http.StatusCreated, rec.Code)
		str := fmt.Sprintf("%s\n", `{"data":null,"meta":{}}`)
		assert.Equal(t, str, rec.Body.String())
	})
}

func createSignerExecutor(ctrl *gomock.Controller) *SignerExecutor {
	u := mock_usecase.NewMockSignIn(ctrl)
	h := handler.NewSigner(u)
	return &SignerExecutor{
		handler: h,
		usecase: u,
	}
}
