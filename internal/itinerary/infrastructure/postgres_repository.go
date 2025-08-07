package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/itinerary/domain"
)

type PostgresItineraryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresItineraryRepository(db *pgxpool.Pool) domain.ItineraryRepository {
	return &PostgresItineraryRepository{db: db}
}

func (r *PostgresItineraryRepository) Save(ctx context.Context, userID uuid.UUID, item *domain.ItineraryItem) error {
	sql := `
		INSERT INTO itens_roteiro (id, id_evento, horario, titulo_atividade, descricao_atividade, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, $6, $7
		FROM eventos e
		WHERE e.id = $2 AND e.id_usuario = $8
	`
	
	cmdTag, err := r.db.Exec(ctx, sql,
		item.ID(),
		item.IDEvento(),
		item.Horario(),
		item.TituloAtividade(),
		item.DescricaoAtividade(),
		item.CreatedAt(),
		item.UpdatedAt(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("falha ao inserir item do roteiro: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("evento não encontrado ou usuário não tem permissão")
	}

	return nil
}

func (r *PostgresItineraryRepository) FindByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.ItineraryItem, error) {
	sql := `
		SELECT id, id_evento, horario, titulo_atividade, descricao_atividade, created_at, updated_at
		FROM itens_roteiro
		WHERE id_evento = $1
		ORDER BY horario ASC
	`

	rows, err := r.db.Query(ctx, sql, eventID)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar itens do roteiro por evento: %w", err)
	}
	defer rows.Close()

	var items []*domain.ItineraryItem
	for rows.Next() {
		var id, idEvento uuid.UUID
		var titulo string
		var descricao *string
		var horario, createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &idEvento, &horario, &titulo, &descricao, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha do item do roteiro: %w", err)
		}

		item := domain.HydrateItineraryItem(id, idEvento, horario, titulo, descricao, createdAt, updatedAt)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas dos itens do roteiro: %w", err)
	}

	return items, nil
}

func (r *PostgresItineraryRepository) FindByID(ctx context.Context, userID, itemID uuid.UUID) (*domain.ItineraryItem, error) {
	sql := `
		SELECT i.id, i.id_evento, i.horario, i.titulo_atividade, i.descricao_atividade, i.created_at, i.updated_at
		FROM itens_roteiro i
		JOIN eventos e ON i.id_evento = e.id
		WHERE i.id = $1 AND e.id_usuario = $2
	`

	var id, idEvento uuid.UUID
	var titulo string
	var descricao *string
	var horario, createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, sql, itemID, userID).Scan(&id, &idEvento, &horario, &titulo, &descricao, &createdAt, &updatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, domain.ErrItemRoteiroNaoEncontrado
		}
		return nil, fmt.Errorf("falha ao consultar item do roteiro por ID: %w", err)
	}

	return domain.HydrateItineraryItem(id, idEvento, horario, titulo, descricao, createdAt, updatedAt), nil
}

func (r *PostgresItineraryRepository) Update(ctx context.Context, userID uuid.UUID, item *domain.ItineraryItem) error {
	sql := `
		UPDATE itens_roteiro SET 
			horario = $1, 
			titulo_atividade = $2, 
			descricao_atividade = $3, 
			updated_at = $4
		FROM eventos e
		WHERE itens_roteiro.id = $5 
			AND itens_roteiro.id_evento = e.id 
			AND e.id_usuario = $6
	`

	cmdTag, err := r.db.Exec(ctx, sql,
		item.Horario(),
		item.TituloAtividade(),
		item.DescricaoAtividade(),
		item.UpdatedAt(),
		item.ID(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("falha ao atualizar item do roteiro: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrItemRoteiroNaoEncontrado
	}

	return nil
}

func (r *PostgresItineraryRepository) Delete(ctx context.Context, userID, itemID uuid.UUID) error {
	sql := `
		DELETE FROM itens_roteiro
		USING eventos e
		WHERE itens_roteiro.id = $1 
			AND itens_roteiro.id_evento = e.id 
			AND e.id_usuario = $2
	`

	cmdTag, err := r.db.Exec(ctx, sql, itemID, userID)
	if err != nil {
		return fmt.Errorf("falha ao deletar item do roteiro: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrItemRoteiroNaoEncontrado
	}

	return nil
}

func (r *PostgresItineraryRepository) FindByEventIDAndUserID(ctx context.Context, eventID, userID uuid.UUID) ([]*domain.ItineraryItem, error) {
	sql := `
		SELECT i.id, i.id_evento, i.horario, i.titulo_atividade, i.descricao_atividade, i.created_at, i.updated_at
		FROM itens_roteiro i
		JOIN eventos e ON i.id_evento = e.id
		WHERE i.id_evento = $1 AND e.id_usuario = $2
		ORDER BY i.horario ASC
	`

	rows, err := r.db.Query(ctx, sql, eventID, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar itens do roteiro por evento e usuário: %w", err)
	}
	defer rows.Close()

	var items []*domain.ItineraryItem
	for rows.Next() {
		var id, idEvento uuid.UUID
		var titulo string
		var descricao *string
		var horario, createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &idEvento, &horario, &titulo, &descricao, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha do item do roteiro: %w", err)
		}

		item := domain.HydrateItineraryItem(id, idEvento, horario, titulo, descricao, createdAt, updatedAt)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas dos itens do roteiro: %w", err)
	}

	return items, nil
}