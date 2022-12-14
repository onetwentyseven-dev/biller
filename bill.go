package biller

import (
	"time"

	"github.com/google/uuid"
)

type Bill struct {
	ID                 uuid.UUID `db:"id" structs:"id" json:"id"`
	UserID             string    `db:"user_id" structs:"user_id" json:"user_id"`
	ProviderID         uuid.UUID `db:"provider_id" structs:"provider_id" json:"provider_id"`
	ProviderName       *string   `db:"provider_name" structs:"-" json:"provider_name"`
	ProviderWebAddress *string   `db:"provider_web_address" structs:"-" json:"provider_web_address,omitempty"`
	Name               string    `db:"name" structs:"name" json:"name"`
	TSCreated          time.Time `db:"ts_created" structs:"ts_created" json:"ts_created"`
	TSUpdated          time.Time `db:"ts_updated" structs:"ts_updated" json:"ts_updated"`
}

type BillSheet struct {
	ID         uuid.UUID `db:"id" structs:"id" json:"id"`
	UserID     string    `db:"user_id" structs:"user_id" json:"user_id"`
	Name       string    `db:"name" structs:"name" json:"name"`
	AmountDue  *float64  `db:"amount_due" structs:"-" json:"amount_due"`
	AmountPaid *float64  `db:"amount_paid" structs:"-" json:"amount_paid"`
	TSCreated  time.Time `db:"ts_created" structs:"ts_created" json:"ts_created"`
	TSUpdated  time.Time `db:"ts_updated" structs:"ts_updated" json:"ts_updated"`
}

type BillSheetEntry struct {
	EntryID            uuid.UUID  `db:"entry_id" structs:"entry_id" json:"entry_id"`
	SheetID            uuid.UUID  `db:"sheet_id" structs:"sheet_id" json:"sheet_id"`
	BillID             uuid.UUID  `db:"bill_id" structs:"bill_id" json:"bill_id"`
	BillName           *string    `db:"bill_name" structs:"-" json:"bill_name,omitempty"`
	ProviderID         *uuid.UUID `db:"provider_id" structs:"-" json:"provider_id,omitempty"`
	ProviderName       *string    `db:"provider_name" structs:"-" json:"provider_name,omitempty"`
	ProviderWebAddress *string    `db:"provider_web_address" structs:"-" json:"provider_web_address,omitempty"`
	DateDue            time.Time  `db:"date_due" structs:"date_due" json:"date_due"`
	AmoutDue           float64    `db:"amount_due" structs:"amount_due" json:"amount_due"`
	ReceiptID          *uuid.UUID `db:"receipt_id" structs:"receipt_id" json:"receipt_id,omitempty"`
	ReceiptLabel       *string    `db:"receipt_name" structs:"-" json:"receipt_name,omitempty"`
	DatePaid           *time.Time `db:"date_paid" structs:"date_paid" json:"date_paid,omitempty"`
	AmountPaid         *float64   `db:"amount_paid" structs:"amount_paid" json:"amount_paid,omitempty"`
	TSCreated          time.Time  `db:"ts_created" structs:"ts_created" json:"ts_created"`
	TSUpdated          time.Time  `db:"ts_updated" structs:"ts_updated" json:"ts_updated"`
}
