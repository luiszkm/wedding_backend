// file: internal/billing/infrastructure/postgres_assinatura_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/billing/domain" // Ajuste o path se necessário
)

type PostgresAssinaturaRepository struct {
	db *pgxpool.Pool
}

// FindByStripeSubscriptionID implements domain.AssinaturaRepository.
func (r *PostgresAssinaturaRepository) FindByStripeSubscriptionID(ctx context.Context, stripeID string) (*domain.Assinatura, error) {
	sql := `SELECT id, id_usuario, id_plano, id_stripe_subscription, data_inicio, data_fim, status FROM assinaturas WHERE id_stripe_subscription = $1`
	row := r.db.QueryRow(ctx, sql, stripeID)

	var assinaturaID, usuarioID, planoID uuid.UUID
	var idStripeSub string
	var dataInicio, dataFim *time.Time
	var status domain.StatusAssinatura

	err := row.Scan(&assinaturaID, &usuarioID, &planoID, &idStripeSub, &dataInicio, &dataFim, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		}
		return nil, fmt.Errorf("falha ao escanear assinatura por id stripe: %w", err)
	}

	var di, df time.Time
	if dataInicio != nil {
		di = *dataInicio
	}
	if dataFim != nil {
		df = *dataFim
	}

	return domain.HydrateAssinatura(assinaturaID, usuarioID, planoID, idStripeSub, di, df, status), nil
}

// FindByID implements domain.AssinaturaRepository.
func (r *PostgresAssinaturaRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Assinatura, error) {
	sql := `SELECT id, id_usuario, id_plano, id_stripe_subscription, data_inicio, data_fim, status FROM assinaturas WHERE id = $1`
	row := r.db.QueryRow(ctx, sql, id)

	var assinaturaID, usuarioID, planoID uuid.UUID
	var idStripeSub *string // Usamos ponteiros para campos que podem ser nulos
	var dataInicio, dataFim *time.Time
	var status domain.StatusAssinatura

	err := row.Scan(&assinaturaID, &usuarioID, &planoID, &idStripeSub, &dataInicio, &dataFim, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAssinaturaNaoEncontrada
		}
		return nil, fmt.Errorf("falha ao escanear assinatura por id: %w", err)
	}

	var di, df time.Time
	var stripeSubID string
	if dataInicio != nil {
		di = *dataInicio
	}
	if dataFim != nil {
		df = *dataFim
	}
	if idStripeSub != nil {
		stripeSubID = *idStripeSub
	}

	return domain.HydrateAssinatura(assinaturaID, usuarioID, planoID, stripeSubID, di, df, status), nil
}

// Update implements domain.AssinaturaRepository.
func (r *PostgresAssinaturaRepository) Update(ctx context.Context, assinatura *domain.Assinatura) error {
	sql := `UPDATE assinaturas SET status = $1, data_inicio = $2, data_fim = $3 WHERE id = $4`
	cmdTag, err := r.db.Exec(ctx, sql,
		assinatura.Status(),
		assinatura.DataInicio(),
		assinatura.DataFim(),
		assinatura.ID(),
	)
	if err != nil {
		return fmt.Errorf("falha ao atualizar assinatura: %w", err)
	}
	// Verifica se alguma linha foi de fato atualizada.
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrAssinaturaNaoEncontrada
	}
	return nil
}
func NewPostgresAssinaturaRepository(db *pgxpool.Pool) domain.AssinaturaRepository {
	return &PostgresAssinaturaRepository{db: db}
}

func (r *PostgresAssinaturaRepository) Save(ctx context.Context, assinatura *domain.Assinatura) error {
	sql := `INSERT INTO assinaturas (id, id_usuario, id_plano, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, sql,
		assinatura.ID(),
		assinatura.IDUsuario(),
		assinatura.IDPlano(),
		assinatura.Status(),
	)
	if err != nil {
		return fmt.Errorf("falha ao salvar assinatura: %w", err)
	}
	return nil
}
