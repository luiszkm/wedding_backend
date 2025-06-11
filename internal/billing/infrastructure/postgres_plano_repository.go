// file: internal/billing/infrastructure/postgres_plano_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/billing/domain"
)

type PostgresPlanoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPlanoRepository(db *pgxpool.Pool) domain.PlanoRepository {
	return &PostgresPlanoRepository{db: db}
}

func (r *PostgresPlanoRepository) ListAll(ctx context.Context) ([]*domain.Plano, error) {
	// Query para buscar todos os planos, ordenados pelo preço do mais barato ao mais caro.
	sql := `
		SELECT id, nome, preco_em_centavos, numero_maximo_eventos, duracao_em_dias, id_stripe_price 
		FROM planos 
		ORDER BY preco_em_centavos ASC;
	`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar planos: %w", err)
	}
	defer rows.Close()

	planos := make([]*domain.Plano, 0)
	for rows.Next() {
		var id uuid.UUID
		var nome string
		var preco, eventos, dias int
		var idStripe string

		if err := rows.Scan(&id, &nome, &preco, &eventos, &dias, &idStripe); err != nil { // <-- scan do novo campo
			return nil, fmt.Errorf("falha ao escanear linha de plano: %w", err)
		}
		plano := domain.HydratePlano(id, nome, idStripe, preco, eventos, dias)
		planos = append(planos, plano)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de planos: %w", err)
	}

	return planos, nil
}

func (r *PostgresPlanoRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Plano, error) {
	sql := `SELECT id, nome, preco_em_centavos, numero_maximo_eventos, duracao_em_dias, id_stripe_price FROM planos WHERE id = $1`
	row := r.db.QueryRow(ctx, sql, id)

	var planoID uuid.UUID
	var nome, idStripe string
	var preco, eventos, dias int

	err := row.Scan(&planoID, &nome, &preco, &eventos, &dias, &idStripe)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlanoNaoEncontrado
		}
		return nil, fmt.Errorf("falha ao buscar plano por id: %w", err)
	}

	return domain.HydratePlano(planoID, nome, idStripe, preco, eventos, dias), nil
}
