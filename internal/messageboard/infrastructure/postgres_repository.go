// file: internal/messageboard/infrastructure/postgres_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/messageboard/domain"
)

type PostgresRecadoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRecadoRepository(db *pgxpool.Pool) domain.RecadoRepository {
	return &PostgresRecadoRepository{db: db}
}

func (r *PostgresRecadoRepository) Save(ctx context.Context, recado *domain.Recado) error {
	sql := `
		INSERT INTO recados (
			id, id_casamento, id_grupo_de_convidados, nome_do_autor, texto, status, eh_favorito, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := r.db.Exec(ctx, sql,
		recado.ID(),
		recado.IDCasamento(),
		recado.IDGrupoDeConvidados(),
		recado.NomeDoAutor(),
		recado.Texto(),
		recado.Status(),
		recado.EhFavorito(),
		recado.DataDeCriacao(),
	)

	if err != nil {
		return fmt.Errorf("falha ao inserir recado no banco de dados: %w", err)
	}
	return nil
}
func (r *PostgresRecadoRepository) ListarPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Recado, error) {
	sql := `
		SELECT
			id, id_casamento, id_grupo_de_convidados, nome_do_autor, texto, status, eh_favorito, created_at
		FROM recados
		WHERE id_casamento = $1
		ORDER BY created_at DESC;
	`
	rows, err := r.db.Query(ctx, sql, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar recados por casamento: %w", err)
	}
	defer rows.Close()

	recados := make([]*domain.Recado, 0)
	for rows.Next() {
		var id, idCasamento, idGrupo uuid.UUID
		var nomeAutor, texto, status string
		var ehFavorito bool
		var createdAt time.Time

		if err := rows.Scan(
			&id, &idCasamento, &idGrupo, &nomeAutor, &texto, &status, &ehFavorito, &createdAt,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha de recado: %w", err)
		}

		recado := domain.HydrateRecado(id, idCasamento, idGrupo, nomeAutor, texto, status, ehFavorito, createdAt)
		recados = append(recados, recado)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de recados: %w", err)
	}

	return recados, nil
}
func (r *PostgresRecadoRepository) FindByID(ctx context.Context, recadoID uuid.UUID) (*domain.Recado, error) {
	sql := `SELECT id, id_casamento, id_grupo_de_convidados, nome_do_autor, texto, status, eh_favorito, created_at FROM recados WHERE id = $1`

	row := r.db.QueryRow(ctx, sql, recadoID)

	var id, idCasamento, idGrupo uuid.UUID
	var nomeAutor, texto, status string
	var ehFavorito bool
	var createdAt time.Time

	if err := row.Scan(&id, &idCasamento, &idGrupo, &nomeAutor, &texto, &status, &ehFavorito, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRecadoNaoEncontrado // Um novo erro de domínio que você pode criar
		}
		return nil, fmt.Errorf("falha ao escanear recado por id: %w", err)
	}

	recado := domain.HydrateRecado(id, idCasamento, idGrupo, nomeAutor, texto, status, ehFavorito, createdAt)
	return recado, nil
}
func (r *PostgresRecadoRepository) Update(ctx context.Context, recado *domain.Recado) error {
	sql := `UPDATE recados SET status = $1, eh_favorito = $2 WHERE id = $3`

	cmdTag, err := r.db.Exec(ctx, sql, recado.Status(), recado.EhFavorito(), recado.ID())
	if err != nil {
		return fmt.Errorf("falha ao atualizar recado: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrRecadoNaoEncontrado // Erro caso o ID não exista mais
	}

	return nil
}

func (r *PostgresRecadoRepository) ListarAprovadosPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Recado, error) {
	// A query implementa as regras de negócio da documentação:
	// - Filtra por status 'APROVADO'
	// - Ordena por 'eh_favorito' e depois pela data de criação
	sql := `
		SELECT
			id, id_casamento, id_grupo_de_convidados, nome_do_autor, texto, status, eh_favorito, created_at
		FROM recados
		WHERE id_casamento = $1 AND status = 'APROVADO'
		ORDER BY eh_favorito DESC, created_at DESC;
	`
	rows, err := r.db.Query(ctx, sql, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar recados aprovados: %w", err)
	}
	defer rows.Close()

	// A lógica de "escanear" as linhas é idêntica à de ListarPorCasamento
	recados := make([]*domain.Recado, 0)
	for rows.Next() {
		var id, idCasamento, idGrupo uuid.UUID
		var nomeAutor, texto, status string
		var ehFavorito bool
		var createdAt time.Time

		if err := rows.Scan(
			&id, &idCasamento, &idGrupo, &nomeAutor, &texto, &status, &ehFavorito, &createdAt,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha de recado aprovado: %w", err)
		}

		recado := domain.HydrateRecado(id, idCasamento, idGrupo, nomeAutor, texto, status, ehFavorito, createdAt)
		recados = append(recados, recado)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de recados aprovados: %w", err)
	}

	return recados, nil
}
