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

type BillSheetRepository struct {
	db *sqlx.DB
}

func NewBillSheetRepository(db *sqlx.DB) *BillSheetRepository {
	return &BillSheetRepository{db}
}

func (r *BillSheetRepository) Sheet(ctx context.Context, userID string, sheetID uuid.UUID) (*biller.BillSheet, error) {

	query := `
		SELECT
			id,
			user_id,
			name,
			(SELECT SUM(amount_due) FROM bill_sheet_entries bse WHERE bse.sheet_id = bs.id) AS amount_due,
			IFNULL((SELECT SUM(amount_paid) FROM bill_sheet_entries bse WHERE bse.sheet_id = bs.id), 0.00) AS amount_paid,
			ts_created,
			ts_updated
		FROM bill_sheets bs
		WHERE id = ? AND user_id = ?
	`

	var sheet = new(biller.BillSheet)
	err := r.db.GetContext(ctx, sheet, query, sheetID, userID)
	return sheet, err

}

func (r *BillSheetRepository) Sheets(ctx context.Context, userID string) ([]*biller.BillSheet, error) {

	query := `
		SELECT
			id,
			user_id,
			name,
			(SELECT SUM(amount_due) FROM bill_sheet_entries bse WHERE bse.sheet_id = bs.id) AS amount_due,
			IFNULL((SELECT SUM(amount_paid) FROM bill_sheet_entries bse WHERE bse.sheet_id = bs.id), 0.00) AS amount_paid,
			ts_created,
			ts_updated
		FROM bill_sheets bs
		WHERE user_id = ?
	`

	var sheets = make([]*biller.BillSheet, 0)
	err := r.db.SelectContext(ctx, &sheets, query, userID)
	return sheets, err

}

func (r *BillSheetRepository) CreateSheet(ctx context.Context, sheet *biller.BillSheet) error {

	now := time.Now()
	sheet.TSCreated = now
	sheet.TSUpdated = now

	query, args, err := sq.Insert("bill_sheets").SetMap(structs.Map(sheet)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err

}

func (r *BillSheetRepository) UpdateSheet(ctx context.Context, sheetID uuid.UUID, sheet *biller.BillSheet) error {

	sheet.ID = sheetID
	now := time.Now()
	sheet.TSUpdated = now

	query, args, err := sq.Update("bill_sheets").
		SetMap(structs.Map(sheet)).
		Where(sq.Eq{"id": sheetID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *BillSheetRepository) DeleteSheet(ctx context.Context, sheetID uuid.UUID) error {

	query, args, err := sq.Delete("bill_sheets").Where(sq.Eq{"id": sheetID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *BillSheetRepository) SheetEntry(ctx context.Context, sheetID uuid.UUID, entryID uuid.UUID) (*biller.BillSheetEntry, error) {

	query := `
		SELECT
			bse.entry_id,
			bse.sheet_id,
			bse.bill_id,
			b.name as bill_name,
			p.name as provider_name,
			p.web_address as provider_web_address,
			bse.date_due,
			bse.amount_due,
			bse.receipt_id,
			r.label as receipt_name,
			bse.date_paid,
			bse.amount_paid,
			bse.ts_created,
			bse.ts_updated
		FROM bill_sheet_entries bse
		LEFT JOIN bills b on bse.bill_id = b.id
		LEFT JOIN providers p on b.provider_id = p.id
		LEFT JOIN receipts r on bse.receipt_id = r.id
		WHERE sheet_id = ? and entry_id = ?
	`

	var sheet = new(biller.BillSheetEntry)
	err := r.db.GetContext(ctx, sheet, query, sheetID, entryID)
	return sheet, err

}

func (r *BillSheetRepository) SheetEntries(ctx context.Context, sheetID uuid.UUID) ([]*biller.BillSheetEntry, error) {

	query := `
		SELECT
			bse.entry_id,
			bse.sheet_id,
			bse.bill_id,
			b.name as bill_name,
			p.name as provider_name,
			p.web_address as provider_web_address,
			bse.date_due,
			bse.amount_due,
			bse.receipt_id,
			r.label as receipt_name,
			bse.date_paid,
			bse.amount_paid,
			bse.ts_created,
			bse.ts_updated
		FROM bill_sheet_entries bse
		LEFT JOIN bills b on bse.bill_id = b.id
		LEFT JOIN providers p on b.provider_id = p.id
		LEFT JOIN receipts r on bse.receipt_id = r.id
		WHERE sheet_id = ?
	`

	var sheets = make([]*biller.BillSheetEntry, 0)
	err := r.db.SelectContext(ctx, &sheets, query, sheetID)
	return sheets, err

}

func (r *BillSheetRepository) CreateSheetEntry(ctx context.Context, entry *biller.BillSheetEntry) error {

	now := time.Now()
	entry.TSCreated = now
	entry.TSUpdated = now

	query, args, err := sq.Insert("bill_sheet_entries").SetMap(structs.Map(entry)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *BillSheetRepository) UpdateSheetEntry(ctx context.Context, entryID uuid.UUID, entry *biller.BillSheetEntry) error {

	now := time.Now()
	entry.TSUpdated = now
	entry.EntryID = entryID

	query, args, err := sq.Update("bill_sheet_entries").
		SetMap(structs.Map(entry)).
		Where(sq.Eq{"entry_id": entryID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update sheet query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *BillSheetRepository) DeleteSheetEntry(ctx context.Context, sheetID uuid.UUID, entryID uuid.UUID) error {

	query, args, err := sq.Delete("bill_sheet_entries").Where(sq.Eq{
		"entry_id": entryID,
		"sheet_id": sheetID,
	}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete sheet entry query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
