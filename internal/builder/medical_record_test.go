package builder_test

import (
	"database/sql"
	"testing"

	"github.com/indrasaputra/orvosi-api/internal/builder"
	"github.com/indrasaputra/orvosi-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildMedicalRecordCreator(t *testing.T) {
	t.Run("successfully build medical record creator", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)

		db := &sql.DB{}

		routes := builder.BuildMedicalRecordCreator(cfg, db)
		assert.NotEmpty(t, routes)
	})
}

func TestBuildMedicalRecordFinder(t *testing.T) {
	t.Run("successfully build medical record finder", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)

		db := &sql.DB{}

		routes := builder.BuildMedicalRecordFinder(cfg, db)
		assert.NotEmpty(t, routes)
	})
}

func TestBuildMedicalRecordUpdater(t *testing.T) {
	t.Run("successfully build medical record updater", func(t *testing.T) {
		cfg, err := config.NewConfig("../../test/fixture/env.valid")
		assert.Nil(t, err)

		db := &sql.DB{}

		routes := builder.BuildMedicalRecordUpdater(cfg, db)
		assert.NotEmpty(t, routes)
	})
}
