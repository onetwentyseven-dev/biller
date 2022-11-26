package mysql

import (
	"context"
	"time"

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

func (r *ProvidersRepository) Providers(ctx context.Context) ([]*biller.Provider, error) {

	query := `
		SELECT 
			id,
			name,
			ts_created,
			ts_updated
		FROM providers
	`

	var providers = make([]*biller.Provider, 0)
	err := r.db.SelectContext(ctx, &providers, query)
	return providers, err

}

func (r *ProvidersRepository) Provider(ctx context.Context, providerID uuid.UUID) (*biller.Provider, error) {

	query := `
		SELECT 
			id,
			name,
			ts_created,
			ts_updated
		FROM providers
		WHERE id = ?
	`

	var provider = new(biller.Provider)
	err := r.db.GetContext(ctx, provider, query, providerID)
	return provider, err

}

func (r *ProvidersRepository) CreateProvider(ctx context.Context, provider *biller.Provider) error {

	now := time.Now()
	provider.TSCreated = now
	provider.TSUpdated = now

	query := `
		INSERT INTO providers (
			id, name, ts_created,ts_updated
		) VALUES (:id, :name, :ts_created, :ts_updated)
	`

	_, err := r.db.NamedExecContext(ctx, query, provider)
	return err

}
