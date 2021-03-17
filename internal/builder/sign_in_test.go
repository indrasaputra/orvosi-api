package builder_test

import (
	"database/sql"
	"testing"

	"github.com/indrasaputra/orvosi-api/internal/builder"
	"github.com/indrasaputra/orvosi-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildrSigner(t *testing.T) {
	t.Run("successfully build signer", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)

		db := &sql.DB{}

		routes := builder.BuildSigner(cfg, db)
		assert.NotEmpty(t, routes)
	})
}
