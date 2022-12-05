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

type ProvidersRepository struct {
	db *sqlx.DB
}

func NewProviderRepository(db *sqlx.DB) *ProvidersRepository {
	return &ProvidersRepository{db}
}

func (r *ProvidersRepository) Providers(ctx context.Context, userID string) ([]*biller.Provider, error) {

	query := `
		SELECT 
			id,
			user_id,
			name,
			web_address,
			ts_created,
			ts_updated
		FROM providers
		WHERE user_id = ?
	`

	var providers = make([]*biller.Provider, 0)
	err := r.db.SelectContext(ctx, &providers, query, userID)
	return providers, err

}

func (r *ProvidersRepository) Provider(ctx context.Context, userID string, providerID uuid.UUID) (*biller.Provider, error) {

	query := `
		SELECT 
			id,
			user_id,
			name,
			web_address,
			ts_created,
			ts_updated
		FROM providers
		WHERE id = ? AND user_id = ?
	`

	var provider = new(biller.Provider)
	err := r.db.GetContext(ctx, provider, query, providerID, userID)
	return provider, err

}

func (r *ProvidersRepository) CreateProvider(ctx context.Context, provider *biller.Provider) error {

	now := time.Now()
	provider.TSCreated = now
	provider.TSUpdated = now

	query, args, err := sq.Insert("providers").SetMap(structs.Map(provider)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create provider query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *ProvidersRepository) UpdateProvider(ctx context.Context, providerID uuid.UUID, provider *biller.Provider) error {

	provider.ID = providerID
	provider.TSUpdated = time.Now()

	query, args, err := sq.Update("providers").SetMap(structs.Map(provider)).Where(sq.Eq{"id": providerID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update provider query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}

func (r *ProvidersRepository) DeleteProvider(ctx context.Context, providerID uuid.UUID) error {

	query, args, err := sq.Delete("providers").Where(sq.Eq{"id": providerID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete provider query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err

}
