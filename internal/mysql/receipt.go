package mysql

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/onetwentyseven-dev/biller"
)

type ReceiptRepository struct {
	db *sqlx.DB
}

func NewReceiptRepository(db *sqlx.DB) *ReceiptRepository {
	return &ReceiptRepository{db}
}

func (r *ReceiptRepository) Receipt(ctx context.Context, userID string, receiptID uuid.UUID) (*biller.Receipt, error) {

	query := `
		SELECT
			id,
			user_id,
			provider_id,
			label,
			date_paid,
			amount_paid,
			ts_created,
			ts_updated
		FROM receipts
		WHERE id = ? AND user_id = ?
	`

	var receipt = new(biller.Receipt)
	err := r.db.GetContext(ctx, receipt, query, receiptID, userID)
	return receipt, err

}

func (r *ReceiptRepository) ReceiptsByProviderID(ctx context.Context, userID string, providerID uuid.UUID) ([]*biller.Receipt, error) {

	query := `
		SELECT
			id,
			user_id, 
			provider_id,
			label,
			date_paid,
			amount_paid,
			ts_created,
			ts_updated
		FROM receipts
		WHERE id = ? and user_id = ?
	`

	var receipts = make([]*biller.Receipt, 0)
	err := r.db.SelectContext(ctx, &receipts, query, providerID, userID)
	return receipts, err

}

func (r *ReceiptRepository) Receipts(ctx context.Context, userID string) ([]*biller.Receipt, error) {

	query := `
		SELECT
			id,
			user_id,
			provider_id,
			label,
			date_paid,
			amount_paid,
			ts_created,
			ts_updated
		FROM receipts
		WHERE user_id = ?
	`

	var receipts = make([]*biller.Receipt, 0)
	err := r.db.SelectContext(ctx, &receipts, query, userID)
	return receipts, err

}

func (r *ReceiptRepository) CreateReceipt(ctx context.Context, receipt *biller.Receipt) error {

	now := time.Now()
	receipt.TSCreated = now
	receipt.TSUpdated = now

	query, args, err := sq.Insert("receipts").SetMap(structs.Map(receipt)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create receipt query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err

}

func (r *ReceiptRepository) UpdateReceipt(ctx context.Context, receiptID uuid.UUID, receipt *biller.Receipt) error {

	receipt.ID = receiptID
	now := time.Now()
	receipt.TSUpdated = now

	query, args, err := sq.Update("receipts").
		SetMap(structs.Map(receipt)).
		Where(sq.Eq{"id": receiptID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *ReceiptRepository) DeleteReceipt(ctx context.Context, receiptID uuid.UUID) error {

	query, args, err := sq.Delete("receipts").Where(sq.Eq{"id": receiptID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
