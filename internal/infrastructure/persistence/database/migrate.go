package database

import (
	"ms-parcel-core/internal/infrastructure/persistence/postgres"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&postgres.DBParcel{},
		&postgres.DBParcelItem{},
		&postgres.DBParcelPayment{},
		&postgres.DBTrackingEvent{},
		&postgres.DBPrintRecord{},
		&postgres.DBPriceRule{},
	)
	if err != nil {
		return err
	}
	return nil
}
