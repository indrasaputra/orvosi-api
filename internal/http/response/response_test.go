package response_test

import (
	"testing"

	"github.com/indrasaputra/orvosi-api/internal/http/response"
	"github.com/stretchr/testify/assert"
)

type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return e.Message
}

func TestNewSuccess(t *testing.T) {
	t.Run("successfully create an instance of success response", func(t *testing.T) {
		type dummy struct{}

		pairs := []struct {
			data interface{}
			meta interface{}
		}{
			{"success data", "success meta"},
			{1, 1},
			{3.14, 6.66},
			{"success data in string", 1.23},
			{dummy{}, dummy{}},
			{dummy{}, "success in meta"},
			{2, nil},
			{nil, 7.89},
		}

		for _, pair := range pairs {
			res := response.NewSuccess(pair.data, pair.meta)

			assert.NotNil(t, res)
			assert.Equal(t, pair.data, res.Data)
			assert.Equal(t, pair.meta, res.Meta)
		}
	})
}

func TestNewError(t *testing.T) {
	t.Run("successfully create an instance of error response", func(t *testing.T) {
		errors := [][]*Error{
			{
				{
					Message: "message from error detail #1",
					Code:    "01-404",
				},
			},
			{
				{
					Message: "message from error detail #2",
					Code:    "02-504",
				},
				{
					Message: "message from error detail #3",
					Code:    "03-503",
				},
			},
		}

		for _, errs := range errors {
			list := buildListOfError(errs)
			res := response.NewError(list...)

			assert.NotNil(t, res)
			assert.Equal(t, list, res.Errors)
			assert.Equal(t, nil, res.Meta)
		}
	})
}

func buildListOfError(errs []*Error) []error {
	res := make([]error, len(errs))
	for i, err := range errs {
		res[i] = err
	}
	return res
}
