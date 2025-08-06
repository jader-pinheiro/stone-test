package utils_test

import (
	"context"
	"stone-test/internal/infra/entity"
	"stone-test/internal/utils"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

func TestParseBrazilianFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10,00", 10.00},
		{"1.000,50", 1000.50},
		{"136.170,000", 136170.0},
	}

	for _, test := range tests {
		result, err := utils.ParseBrazilianFloat(test.input)
		if err != nil {
			t.Errorf("Erro ao parsear '%s': %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("Esperado %f, obtido %f para entrada '%s'", test.expected, result, test.input)
		}
	}
}

func TestParseLine(t *testing.T) {
	line := "2025-08-01;WINV25;0;136170,000;1;090000002;10;1;2025-08-01;4090;93"

	stock, err := utils.ParseLine(line)
	if err != nil {
		t.Fatalf("Erro ao parsear linha: %v", err)
	}

	if stock.InstrumentCode != "WINV25" {
		t.Errorf("InstrumentCode esperado: WINV25, obtido: %s", stock.InstrumentCode)
	}

	expectedPrice := 136170.0
	if stock.BusinessPrice != expectedPrice {
		t.Errorf("Preço esperado: %f, obtido: %f", expectedPrice, stock.BusinessPrice)
	}

	expectedQty := int64(1)
	if stock.NegotiatedQuantity != expectedQty {
		t.Errorf("Quantidade esperada: %d, obtida: %d", expectedQty, stock.NegotiatedQuantity)
	}

	loc, _ := time.LoadLocation("America/Sao_Paulo")
	expectedDate := time.Date(2025, 8, 1, 0, 0, 0, 0, loc)

	if !stock.BusinessDate.Equal(expectedDate) {
		t.Errorf("Data esperada: %s, obtida: %s", expectedDate, stock.BusinessDate)
	}

	expectedTime := time.Date(2025, 8, 1, 9, 0, 0, 2_000_000, time.UTC)
	if !stock.ClosingTime.Equal(expectedTime) {
		t.Errorf("Hora de fechamento esperada: %v, obtida: %v", expectedTime, stock.ClosingTime)
	}
}

type MockInserter struct {
	mock.Mock
}

func (m *MockInserter) InsertBatch(ctx context.Context, conn *pgxpool.Pool, batch []entity.Stocks) error {
	args := m.Called(ctx, conn, batch)
	return args.Error(0)
}
