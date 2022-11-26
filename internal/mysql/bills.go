package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/onetwentyseven-dev/biller"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
)

type BillsRepository struct {
	db *sqlx.DB
}

func NewBillsRepository(db *sqlx.DB) *BillsRepository {

	return &BillsRepository{db}
}

func (r *BillsRepository) Bill(ctx context.Context, billID uuid.UUID) (*biller.Bill, error) {

	query := `
		SELECT
			id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
		WHERE id = ?
	`

	var bill = new(biller.Bill)
	err := r.db.GetContext(ctx, bill, query, billID)
	return bill, err

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

	query, args, err := sq.Insert("bills").SetMap(structs.Map(bill)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create bill query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err

}

func (r *BillsRepository) UpdateBill(ctx context.Context, billID uuid.UUID, bill *biller.Bill) error {

	now := time.Now()
	bill.ID = billID
	bill.TSUpdated = now

	query, args, err := sq.Update("bills").
		SetMap(structs.Map(bill)).
		Where(sq.Eq{"id": billID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update bill query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err

}

func (r *BillsRepository) DeleteBill(ctx context.Context, billID uuid.UUID) error {

	query := `
		SELECT
			id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, billID)
	return err

}
