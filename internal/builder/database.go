package builder

import (
	"database/sql"
	"fmt"

	"github.com/indrasaputra/orvosi-api/internal/config"
)

// BuildSQLDatabase builds *sql.DB from given config.
func BuildSQLDatabase(driver string, cfg *config.Config) (*sql.DB, error) {
	sqlCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := sql.Open(driver, sqlCfg)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	return db, nil
}
