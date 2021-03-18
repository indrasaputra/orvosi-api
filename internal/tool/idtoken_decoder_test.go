package tool_test

import (
	"io/ioutil"
	"testing"

	"github.com/indrasaputra/orvosi-api/entity"
	"github.com/indrasaputra/orvosi-api/internal/tool"
	"github.com/stretchr/testify/assert"
)

const (
	audience = "test-audience"
)

func TestNewIDTokenDecoder(t *testing.T) {
	t.Run("successfully create an instance of IDTokenDecoder", func(t *testing.T) {
		dec := tool.NewIDTokenDecoder(audience)
		assert.NotNil(t, dec)
	})
}

func TestIDTokenDecoder_Decode(t *testing.T) {
	t.Run("fail to decode invalid google id token", func(t *testing.T) {
		token, ierr := ioutil.ReadFile("../../test/fixture/token.txt")
		assert.Nil(t, ierr)

		dec := tool.NewIDTokenDecoder(audience)
		user, err := dec.Decode(string(token))

		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInvalidGoogleToken.Code, err.Code)
		assert.Nil(t, user)
	})
}
