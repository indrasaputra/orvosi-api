package entity_test

import (
	"fmt"
	"testing"

	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	t.Run("successfully create an instance of Error", func(t *testing.T) {
		err := entity.NewError("01-001", "Internal server error")

		assert.NotNil(t, err)
		assert.Equal(t, "01-001", err.Code)
		assert.Equal(t, "Internal server error", err.Message)
		assert.Equal(t, "Internal server error", err.Error())
	})
}

func TestError_Error(t *testing.T) {
	t.Run("internal message is derived from public message", func(t *testing.T) {
		messages := []string{
			"Internal server error",
			"Unauthorized",
		}

		for _, msg := range messages {
			err := entity.NewError("01-001", msg)
			assert.Equal(t, msg, err.Error())
		}
	})

	t.Run("internal message is concated from public message and wrapped message", func(t *testing.T) {
		messages := []string{
			"wrapped message #1",
			"wrapped message #2",
		}

		for _, msg := range messages {
			err := entity.NewError("01-001", "public message")
			err = entity.WrapError(err, msg)

			assert.Equal(t, fmt.Sprintf("public message. %s", msg), err.Error())
		}
	})
}

func TestWrapError(t *testing.T) {
	t.Run("new error has the same code and public message as base", func(t *testing.T) {
		ori := entity.NewError("01-001", "initial message")

		err := entity.WrapError(ori, "additional message #1")
		assert.Equal(t, ori.Code, err.Code)
		assert.Equal(t, ori.Message, err.Message)
	})

	t.Run("message is wrapped exactly at the end of current message and separated by :", func(t *testing.T) {
		err := entity.NewError("01-001", "initial message")

		err = entity.WrapError(err, "additional message #1")
		assert.Equal(t, "initial message. additional message #1", err.Error())

		err = entity.WrapError(err, "additional message #2")
		assert.Equal(t, "initial message. additional message #1. additional message #2", err.Error())
	})
}
