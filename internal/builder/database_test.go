package builder_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/indrasaputra/orvosi-api/internal/builder"
	"github.com/indrasaputra/orvosi-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildSQLDatabase(t *testing.T) {
	t.Run("fail to build sql.DB due to unknown driver", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)

		db, err := builder.BuildSQLDatabase("unknown", cfg)

		assert.NotNil(t, err)
		assert.Nil(t, db)
	})

	t.Run("successfully build sql.DB using known driver (postgres)", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)

		db, err := builder.BuildSQLDatabase("postgres", cfg)

		assert.Nil(t, err)
		assert.NotNil(t, db)
	})
}
