package biller

import (
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	ID         uuid.UUID `db:"id" structs:"id" json:"id"`
	ProviderID uuid.UUID `db:"provider_id" structs:"provider_id" json:"provider_id"`
	Label      string    `db:"label" structs:"label" json:"label"`
	DatePaid   time.Time `db:"date_paid" structs:"date_paid" json:"date_paid"`
	AmountPaid float64   `db:"amount_paid" structs:"amount_paid" json:"amount_paid"`
	TSCreated  time.Time `db:"ts_created" structs:"ts_created" json:"ts_created"`
	TSUpdated  time.Time `db:"ts_updated" structs:"ts_updated" json:"ts_updated"`
}
