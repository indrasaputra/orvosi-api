package builder

import (
	"database/sql"

	"github.com/indrasaputra/orvosi-api/internal/config"
	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/router"
	"github.com/indrasaputra/orvosi-api/internal/repository"
	"github.com/indrasaputra/orvosi-api/usecase"
)

// BuildMedicalRecordCreator builds medical record creation workflow
// starting from handler down to repository.
func BuildMedicalRecordCreator(cfg *config.Config, db *sql.DB) []*router.Route {
	ins := repository.NewMedicalRecordInserter(db)
	uc := usecase.NewMedicalRecordCreator(ins)
	hdr := handler.NewMedicalRecordCreator(uc)
	return router.MedicalRecordCreator(hdr)
}

// BuildMedicalRecordFinder builds medical record find workflow
// starting from handler down to repository.
func BuildMedicalRecordFinder(cfg *config.Config, db *sql.DB) []*router.Route {
	sel := repository.NewMedicalRecordSelector(db)
	uc := usecase.NewMedicalRecordFinder(sel)
	hdr := handler.NewMedicalRecordFinder(uc)
	return router.MedicalRecordFinder(hdr)
}

// BuildMedicalRecordUpdater builds medical record update workflow
// starting from handler down to repository.
func BuildMedicalRecordUpdater(cfg *config.Config, db *sql.DB) []*router.Route {
	up := repository.NewMedicalRecordUpdater(db)
	uc := usecase.NewMedicalRecordUpdater(up)
	hdr := handler.NewMedicalRecordUpdater(uc)
	return router.MedicalRecordUpdater(hdr)
}
