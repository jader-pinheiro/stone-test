package db

import (
	"stone-test/internal/infra/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := "postgres://postgres:passwd@127.0.0.1:5432/stone_test?sslmode=disable"

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
