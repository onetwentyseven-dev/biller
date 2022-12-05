package biller

import (
	"time"

	"github.com/google/uuid"
)

type Provider struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     string    `db:"user_id" structs:"user_id" json:"user_id"`
	Name       string    `db:"name" json:"name"`
	WebAddress *string   `db:"web_address" json:"web_address,omitempty"`
	TSCreated  time.Time `db:"ts_created" json:"ts_created"`
	TSUpdated  time.Time `db:"ts_updated" json:"ts_updated"`
}
