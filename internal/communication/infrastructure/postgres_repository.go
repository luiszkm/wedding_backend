package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/luiszkm/wedding_backend/internal/communication/domain"
)

type PostgresComunicadoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresComunicadoRepository(db *pgxpool.Pool) *PostgresComunicadoRepository {
	return &PostgresComunicadoRepository{
		db: db,
	}
}

func (r *PostgresComunicadoRepository) Criar(ctx context.Context, comunicado *domain.Comunicado) error {
	query := `
		INSERT INTO comunicados (id, id_evento, titulo, mensagem, data_publicacao)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query,
		comunicado.ID(),
		comunicado.IDEvento(),
		comunicado.Titulo(),
		comunicado.Mensagem(),
		comunicado.DataPublicacao(),
	)

	return err
}

func (r *PostgresComunicadoRepository) BuscarPorEvento(ctx context.Context, idEvento uuid.UUID) ([]*domain.Comunicado, error) {
	query := `
		SELECT id, id_evento, titulo, mensagem, data_publicacao
		FROM comunicados
		WHERE id_evento = $1
		ORDER BY data_publicacao DESC`

	rows, err := r.db.Query(ctx, query, idEvento)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comunicados []*domain.Comunicado
	for rows.Next() {
		var id, idEventoRow uuid.UUID
		var titulo, mensagem string
		var dataPublicacao sql.NullTime

		err := rows.Scan(&id, &idEventoRow, &titulo, &mensagem, &dataPublicacao)
		if err != nil {
			return nil, err
		}

		comunicado := domain.HydrateComunicado(
			id,
			idEventoRow,
			titulo,
			mensagem,
			dataPublicacao.Time,
		)
		comunicados = append(comunicados, comunicado)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comunicados, nil
}

func (r *PostgresComunicadoRepository) BuscarPorID(ctx context.Context, id uuid.UUID) (*domain.Comunicado, error) {
	query := `
		SELECT id, id_evento, titulo, mensagem, data_publicacao
		FROM comunicados
		WHERE id = $1`

	var idRow, idEvento uuid.UUID
	var titulo, mensagem string
	var dataPublicacao sql.NullTime

	err := r.db.QueryRow(ctx, query, id).Scan(&idRow, &idEvento, &titulo, &mensagem, &dataPublicacao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrComunicadoNaoEncontrado
		}
		return nil, err
	}

	return domain.HydrateComunicado(
		idRow,
		idEvento,
		titulo,
		mensagem,
		dataPublicacao.Time,
	), nil
}

func (r *PostgresComunicadoRepository) Editar(ctx context.Context, comunicado *domain.Comunicado) error {
	query := `
		UPDATE comunicados 
		SET titulo = $1, mensagem = $2
		WHERE id = $3`

	result, err := r.db.Exec(ctx, query,
		comunicado.Titulo(),
		comunicado.Mensagem(),
		comunicado.ID(),
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrComunicadoNaoEncontrado
	}

	return nil
}

func (r *PostgresComunicadoRepository) Deletar(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comunicados WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrComunicadoNaoEncontrado
	}

	return nil
}
