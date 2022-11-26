package mysql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/onetwentyseven-dev/biller"
)

type BillsRepository struct {
	db *sqlx.DB
}

func NewBillsRepository(db *sqlx.DB) *BillsRepository {
	return &BillsRepository{db}
}

func (r *BillsRepository) Bills(ctx context.Context) ([]*biller.Bill, error) {

	query := `
		SELECT
			id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
	`

	var bills = make([]*biller.Bill, 0)
	err := r.db.SelectContext(ctx, &bills, query)
	return bills, err

}

func (r *BillsRepository) BillsByProvider(ctx context.Context, providerID uuid.UUID) ([]*biller.Bill, error) {
	query := `
		SELECT
			id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
		WHERE provider_id = ?
	`

	var bills = make([]*biller.Bill, 0)
	err := r.db.SelectContext(ctx, &bills, query, providerID)
	return bills, err
}

func (r *BillsRepository) CreateBill(ctx context.Context, bill *biller.Bill) error {

	now := time.Now()
	bill.TSCreated = now
	bill.TSUpdated = now

	query := `
		INSERT INTO bills (
			id, provider_id, name, ts_created, ts_updated
		) VALUES (
			:id, :provider_id, :name, :ts_created, :ts_updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, bill)

	return err

}
