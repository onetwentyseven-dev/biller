package biller

import (
	"time"

	"github.com/google/uuid"
)

type Provider struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	TSCreated time.Time `db:"ts_created" json:"ts_created"`
	TSUpdated time.Time `db:"ts_updated" json:"ts_updated"`
}
