package database

import (
	"fmt"
	"log"

	"ms-parcel-core/internal/config"
	mypostgres "ms-parcel-core/internal/infrastructure/persistence/postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	mypostgres.RegisterTenantScope(db)

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Error enabling uuid-ossp: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}
