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
			web_address,
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
			web_address,
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
			id, name, web_address, ts_created,ts_updated
		) VALUES (:id, :name, :web_address, :ts_created, :ts_updated)
	`

	_, err := r.db.NamedExecContext(ctx, query, provider)
	return err

}

func (r *ProvidersRepository) UpdateProvider(ctx context.Context, providerID uuid.UUID, provider *biller.Provider) error {

	provider.ID = providerID
	provider.TSUpdated = time.Now()

	query := `
		UPDATE providers set name = :name, web_address = :web_address WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, provider)
	return err

}

func (r *ProvidersRepository) DeleteProvider(ctx context.Context, providerID uuid.UUID) error {
	query := `DELETE FROM providers where id = ?`
	_, err := r.db.ExecContext(ctx, query, providerID)
	return err
}
