package tool_test

import (
	"testing"

	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/tool"
	"github.com/stretchr/testify/assert"
)

const (
	audience = "test-audience"
	token    = "eyJhbGciOiJIUzI1NiIsImtpZCI6ImVlYTFiMWY0MjgwN2E4Y2MxMzZhMDNhM2MxNmQyOWRiODI5NmRhZjAiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoidGVzdC1hdWRpZW5jZSIsImF1ZCI6InRlc3QtYXVkaWVuY2UiLCJzdWIiOiIxMjM0NTY3ODkwIiwiaGQiOiJkdW1teS5jb20iLCJlbWFpbCI6ImR1bW15QGR1bW15LmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoia216d2E4YXdhYSIsIm5hbWUiOiJEdW1teSBBY2NvdW50IiwicGljdHVyZSI6Imh0dHBzOi8vZHVtbXkuY29tL3Bob3RvLmpwZyIsImdpdmVuX25hbWUiOiJEdW1teSIsImZhbWlseV9uYW1lIjoiQWNjb3VudCIsImxvY2FsZSI6ImVuIiwiaWF0IjoyNjExNjQzMjAzLCJleHAiOjI2MTE2NDY4MDMsImp0aSI6InNvbWVyYW5kb21hbHBoYW51bWVyaWsxMjM0NTY3ODkwIn0.raVMs-JvbG4h2lP5p8RfPC5YAyBDe4nCBCiGXGWE88g"
)

func TestNewIDTokenDecoder(t *testing.T) {
	t.Run("successfully create an instance of IDTokenDecoder", func(t *testing.T) {
		dec := tool.NewIDTokenDecoder(audience)
		assert.NotNil(t, dec)
	})
}

func TestIDTokenDecoder_Decode(t *testing.T) {
	t.Run("fail to decode invalid google id token", func(t *testing.T) {
		dec := tool.NewIDTokenDecoder(audience)
		user, err := dec.Decode(token)

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidGoogleToken.Code, err.Code)
		assert.Nil(t, user)
	})
}
