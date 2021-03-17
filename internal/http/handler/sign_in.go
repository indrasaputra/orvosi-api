package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/http/response"
	"github.com/indrasaputra/orvosi-api/usecase"
)

// Signer handles HTTP request and response
// for sign in.
type Signer struct {
	signin usecase.SignIn
}

// NewSigner creates an instance of Medical.
func NewSigner(signin usecase.SignIn) *Signer {
	return &Signer{
		signin: signin,
	}
}

// SignIn handles `POST /sign-in` endpoint.
func (s *Signer) SignIn(ctx echo.Context) error {
	user, err := extractUserFromRequestContext(ctx.Request().Context())
	if err != nil {
		res := response.NewError(err)
		ctx.JSON(http.StatusInternalServerError, res)
		return err
	}

	if err := s.signin.SignIn(ctx.Request().Context(), user); err != nil {
		res := response.NewError(err)
		status := http.StatusInternalServerError
		if err.Code == entity.ErrEmptyUser.Code {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, res)
		return err
	}

	ctx.JSON(http.StatusCreated, response.NewSuccess(nil, response.EmptyMeta{}))
	return nil
}
