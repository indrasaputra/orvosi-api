package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWithJWTDecoder(t *testing.T) {
	t.Run("request doesn't contain authorization", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		hdr := createHandler()
		dec := createNormalDecoder()
		hdr = middleware.WithJWTDecoder(dec)(hdr)

		err := hdr(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, entity.ErrUnauthorized, err)
	})

	t.Run("token length is not 2 (invalid)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "invalid-token")
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		hdr := createHandler()
		dec := createNormalDecoder()
		hdr = middleware.WithJWTDecoder(dec)(hdr)

		err := hdr(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, entity.ErrUnauthorized, err)
	})

	t.Run("token prefix is not `Bearer`", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "not-bearer jwt-token")
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		hdr := createHandler()
		dec := createNormalDecoder()
		hdr = middleware.WithJWTDecoder(dec)(hdr)

		err := hdr(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, entity.ErrUnauthorized, err)
	})

	t.Run("jwt decoder func returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer jwt-token")
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		hdr := createHandler()
		dec := createErrorDecoder()
		hdr = middleware.WithJWTDecoder(dec)(hdr)

		err := hdr(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, entity.ErrUnauthorized, err)
	})

	t.Run("successfully decode jwt", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer jwt-token")
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		hdr := createHandler()
		dec := createNormalDecoder()
		hdr = middleware.WithJWTDecoder(dec)(hdr)

		err := hdr(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		val := ctx.Request().Context().Value(middleware.ContextKeyUser)
		user, ok := val.(*entity.User)
		assert.True(t, ok)
		assert.NotNil(t, user)
		assert.Equal(t, "dummy@jwtmiddleware.com", user.Email)
	})
}

func TestWithContentType(t *testing.T) {
	t.Run("request can't be continued due to wrong content type", func(t *testing.T) {
		tables := []struct {
			expected string
			actual   string
		}{
			{echo.MIMEApplicationForm, echo.MIMEApplicationJSON},
			{echo.MIMEApplicationJSON, echo.MIMEApplicationJavaScriptCharsetUTF8},
			{echo.MIMEApplicationJSON, echo.MIMEApplicationProtobuf},
			{echo.MIMETextXML, echo.MIMEMultipartForm},
		}

		for _, table := range tables {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, table.actual)
			rec := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(req, rec)

			hdr := createHandler()
			hdr = middleware.WithContentType(table.expected)(hdr)

			err := hdr(ctx)

			assert.NotNil(t, err)
			assert.Equal(t, entity.ErrWrongContentType, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("successfully continue the request when content-type is right", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		hdr := createHandler()
		hdr = middleware.WithContentType(echo.MIMEApplicationJSON)(hdr)

		err := hdr(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func createErrorDecoder() middleware.JWTDecoder {
	return func(token string) (*entity.User, *entity.Error) {
		return nil, entity.ErrUnauthorized
	}
}

func createNormalDecoder() middleware.JWTDecoder {
	return func(token string) (*entity.User, *entity.Error) {
		return &entity.User{Email: "dummy@jwtmiddleware.com"}, nil
	}
}

func createHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
}
