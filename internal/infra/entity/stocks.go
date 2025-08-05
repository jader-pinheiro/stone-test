package entity

import (
	"time"
)

type Stocks struct {
	ID                 uint
	BusinessDate       time.Time `gorm:"index:idx_instrument_date,priority:2" json:"DataNegocio"`
	InstrumentCode     string    `gorm:"type:varchar(50);index:idx_instrument_date,priority:1" json:"CodigoInstrumento"`
	BusinessPrice      float64   `json:"PrecoNegocio"`
	NegotiatedQuantity int64     `json:"QuantidadeNegociada"`
	ClosingTime        time.Time `json:"HoraFechamento"`
}
