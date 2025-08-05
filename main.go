package main

import (
	"context"
	"log"
	"os"
	"stone-test/internal/infra/db"
	"stone-test/internal/ui/routes"
	"stone-test/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := db.Connect()

	if err != nil {
		log.Fatalf("CRITICAL: Failed to connect to the database: %v", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalf("Failed to get raw database connection: %v", err)
	}

	//connection to PostgreSQL using pgx
	ctx := context.Background()
	dsnPgx := os.Getenv("GORM_DSN")
	pool, err := pgxpool.New(ctx, dsnPgx)
	if err != nil {
		log.Fatalf("Erro ao criar pool: %v", err)
	}
	defer pool.Close()

	_, errPgx := utils.ProcessFileContent(ctx, pool)
	if errPgx != nil {
		log.Fatalf("Erro no processamento: %v", err)
	}

	defer sqlDB.Close()
	routes.GetRoutes(conn)
}
