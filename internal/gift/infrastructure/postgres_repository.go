// file: internal/gift/infrastructure/postgres_repository.go
package infrastructure

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/gift/domain"
)

type PostgresPresenteRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPresenteRepository(db *pgxpool.Pool) domain.PresenteRepository {
	return &PostgresPresenteRepository{db: db}
}

func (r *PostgresPresenteRepository) Save(ctx context.Context, presente *domain.Presente) error {
	sql := `
		INSERT INTO presentes (
			id, id_casamento, nome, descricao, eh_favorito, status, 
			detalhes_tipo, detalhes_link_loja, foto_url, categoria -- Adicionada a nova coluna
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10); -- Adicionado o novo placeholder
	`
	detalhes := presente.Detalhes()
	_, err := r.db.Exec(ctx, sql,
		presente.ID(),
		presente.IDCasamento(),
		presente.Nome(),
		presente.Descricao(),
		presente.EhFavorito(),
		presente.Status(),
		detalhes.Tipo,
		detalhes.LinkDaLoja,
		presente.FotoURL(),
		presente.Categoria(),
	)

	if err != nil {
		return fmt.Errorf("falha ao inserir presente no banco de dados: %w", err)
	}
	return nil
}

func (r *PostgresPresenteRepository) ListarDisponiveisPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Presente, error) {
	// A query seleciona apenas os presentes com status 'DISPONIVEL' para um casamento específico.
	sql := `
		SELECT
			id, id_casamento, nome, descricao, foto_url, status, categoria, eh_favorito,
			detalhes_tipo, detalhes_link_loja
		FROM presentes
		WHERE id_casamento = $1 AND status = 'DISPONIVEL'
		ORDER BY eh_favorito DESC, nome ASC;
	`
	rows, err := r.db.Query(ctx, sql, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar presentes disponíveis: %w", err)
	}
	defer rows.Close()

	presentes := make([]*domain.Presente, 0)
	for rows.Next() {
		var d domain.DetalhesPresente
		var id, idCasamento uuid.UUID
		var nome, descricao, fotoURL, status, categoria, detalhesTipo, detalhesLinkLoja string
		var ehFavorito bool

		// Usamos ponteiros para campos que podem ser nulos no banco de dados.
		var pDesc, pFotoURL, pLinkLoja, pCategoria *string

		if err := rows.Scan(
			&id, &idCasamento, &nome, &pDesc, &pFotoURL, &status, &pCategoria, &ehFavorito,
			&detalhesTipo, &pLinkLoja,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha de presente: %w", err)
		}

		// Atribui valores dos ponteiros, se não forem nulos.
		if pDesc != nil {
			descricao = *pDesc
		}
		if pFotoURL != nil {
			fotoURL = *pFotoURL
		}
		if pLinkLoja != nil {
			detalhesLinkLoja = *pLinkLoja
		}
		if pCategoria != nil {
			categoria = *pCategoria
		}

		d.Tipo = detalhesTipo
		d.LinkDaLoja = detalhesLinkLoja

		presente := domain.HydratePresente(id, idCasamento, nome, descricao, fotoURL, status, categoria, ehFavorito, d)
		presentes = append(presentes, presente)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de presentes: %w", err)
	}

	return presentes, nil
}
