package biller

import (
	"time"

	"github.com/google/uuid"
)

type Bill struct {
	ID         uuid.UUID `db:"id" json:"id"`
	ProviderID uuid.UUID `db:"provider_id" json:"provider_id"`
	Name       string    `db:"name" json:"name"`
	TSCreated  time.Time `db:"ts_created" json:"ts_created"`
	TSUpdated  time.Time `db:"ts_updated" json:"ts_updated"`
}
