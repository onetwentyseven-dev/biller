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

func (r *BillsRepository) Bill(ctx context.Context, userID string, billID uuid.UUID) (*biller.Bill, error) {

	query := `
		SELECT
			id,
			user_id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
		WHERE id = ? and user_id = ?
	`

	var bill = new(biller.Bill)
	err := r.db.GetContext(ctx, bill, query, billID, userID)
	return bill, err

}

func (r *BillsRepository) Bills(ctx context.Context, userID string) ([]*biller.Bill, error) {

	query := `
		SELECT
			id,
			user_id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
		WHERE user_id = ?
	`

	var bills = make([]*biller.Bill, 0)
	err := r.db.SelectContext(ctx, &bills, query, userID)
	return bills, err

}

func (r *BillsRepository) BillsByProvider(ctx context.Context, userID string, providerID uuid.UUID) ([]*biller.Bill, error) {
	query := `
		SELECT
			id,
			user_id,
			provider_id,
			name,
			ts_created,
			ts_updated
		FROM bills
		WHERE provider_id = ? AND user_id = ?
	`

	var bills = make([]*biller.Bill, 0)
	err := r.db.SelectContext(ctx, &bills, query, providerID, userID)
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

	query, args, err := sq.Delete("bills").Where(sq.Eq{"id": billID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update bill query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
