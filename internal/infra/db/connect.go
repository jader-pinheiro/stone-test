package db

import (
	"fmt"
	"os"
	"stone-test/internal/infra/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("GORM_DSN")

	if dsn == "" {
		return nil, fmt.Errorf("GORM_DSN environment variable not set")
	}

	// Conectando ao banco de dados PostgreSQL com GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.Stocks{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
