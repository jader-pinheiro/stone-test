package data

import (
	"context"
	"fmt"
	"stone-test/internal/infra/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

func InsertBatch(db *gorm.DB, batch []entity.Stocks, ctx context.Context) error {
	return db.Session(&gorm.Session{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}).WithContext(ctx).CreateInBatches(batch, len(batch)).Error
}

func InsertBatchCopy(ctx context.Context, pool *pgxpool.Pool, batch []entity.Stocks) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("falha ao adquirir conexão do pool: %w", err)
	}
	defer conn.Release()

	rows := make([][]interface{}, len(batch))
	for i, stock := range batch {
		rows[i] = []interface{}{
			stock.BusinessDate,
			stock.InstrumentCode,
			stock.BusinessPrice,
			stock.NegotiatedQuantity,
			stock.ClosingTime,
		}
	}

	_, err = conn.Conn().CopyFrom(ctx,
		pgx.Identifier{"stocks"},
		[]string{"business_date", "instrument_code", "business_price", "negotiated_quantity", "closing_time"},
		pgx.CopyFromRows(rows),
	)
	return err
}

type ticker struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value,omitempty"`
	MaxDailyVolume int64   `json:"max_daily_volume,omitempty"`
}

func GetTickerData(db *gorm.DB, ctx context.Context, instrumentCode string, businessDate *time.Time) (ticker, error) {
	var result ticker

	query := db.WithContext(ctx).
		Table("stocks").
		Select("instrument_code AS ticker, MAX(business_price) AS max_range_value, MAX(negotiated_quantity) AS max_daily_volume").
		Where("instrument_code = ?", instrumentCode)

	if businessDate != nil {
		query = query.Where("business_date >= ?", *businessDate)
	} else {
		end := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		start := end.AddDate(0, 0, -6)
		query = query.Where("business_date BETWEEN ? AND ?", start, end)
	}

	query = query.Group("instrument_code")

	if err := query.Scan(&result).Error; err != nil {
		return ticker{}, fmt.Errorf("erro ao buscar dados do ticker %s: %w", instrumentCode, err)
	}

	return result, nil
}
