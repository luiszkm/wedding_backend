// file: internal/gallery/infrastructure/postgres_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/gallery/domain"
)

type PostgresFotoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresFotoRepository(db *pgxpool.Pool) domain.FotoRepository {
	return &PostgresFotoRepository{db: db}
}

// SalvarMultiplas usa bulk inserts (CopyFrom) para alta performance.
func (r *PostgresFotoRepository) SalvarMultiplas(ctx context.Context, fotos []*domain.Foto) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação de fotos: %w", err)
	}
	defer tx.Rollback(ctx)

	// Prepara os dados para o bulk insert na tabela 'fotos'
	fotosRows := make([][]any, len(fotos))
	for i, f := range fotos {
		fotosRows[i] = []any{f.ID(), f.IDCasamento(), f.StorageKey(), f.URLPublica(), f.EhFavorito(), f.CreatedAt()}
	}

	// Insere todas as fotos de uma vez
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"fotos"},
		[]string{"id", "id_casamento", "storage_key", "url_publica", "eh_favorito", "created_at"},
		pgx.CopyFromRows(fotosRows),
	)
	if err != nil {
		return fmt.Errorf("falha no bulk insert em 'fotos': %w", err)
	}

	// Prepara os dados para o bulk insert na tabela 'fotos_rotulos'
	rotulosRows := [][]any{}
	for _, f := range fotos {
		for _, rotulo := range f.Rotulos() {
			rotulosRows = append(rotulosRows, []any{f.ID(), rotulo})
		}
	}

	if len(rotulosRows) > 0 {
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"fotos_rotulos"},
			[]string{"id_foto", "nome_rotulo"},
			pgx.CopyFromRows(rotulosRows),
		)
		if err != nil {
			return fmt.Errorf("falha no bulk insert em 'fotos_rotulos': %w", err)
		}
	}

	return tx.Commit(ctx)
}
func (r *PostgresFotoRepository) ListarPublicasPorCasamento(ctx context.Context, casamentoID uuid.UUID, filtroRotulo domain.Rotulo) ([]*domain.Foto, error) {
	// A base da query faz o JOIN e agrega os rótulos de cada foto em um array.
	baseSQL := `
		SELECT
			f.id, f.id_casamento, f.storage_key, f.url_publica, f.eh_favorito, f.created_at,
			COALESCE(ARRAY_AGG(fr.nome_rotulo) FILTER (WHERE fr.nome_rotulo IS NOT NULL), '{}') as rotulos
		FROM fotos f
		LEFT JOIN fotos_rotulos fr ON f.id = fr.id_foto
		WHERE f.id_casamento = $1
	`
	args := []interface{}{casamentoID}

	// Adiciona o filtro de rótulo dinamicamente, se ele for fornecido.
	if filtroRotulo != "" {
		// Esta subquery verifica se o rótulo filtrado existe para a foto.
		baseSQL += ` AND EXISTS (SELECT 1 FROM fotos_rotulos fr_filter WHERE fr_filter.id_foto = f.id AND fr_filter.nome_rotulo = $2)`
		args = append(args, filtroRotulo)
	}

	// Finaliza a query com o agrupamento e a ordenação exigida pela documentação.
	finalSQL := baseSQL + ` GROUP BY f.id ORDER BY f.eh_favorito DESC, f.created_at DESC`

	rows, err := r.db.Query(ctx, finalSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar fotos públicas: %w", err)
	}
	defer rows.Close()

	fotos := make([]*domain.Foto, 0)
	for rows.Next() {
		var id, idCasamento uuid.UUID
		var storageKey, urlPublica string
		var ehFavorito bool
		var createdAt time.Time
		var rotulos []string // pgx consegue escanear um array do SQL para um slice de string

		if err := rows.Scan(&id, &idCasamento, &storageKey, &urlPublica, &ehFavorito, &createdAt, &rotulos); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha de foto: %w", err)
		}

		domainRotulos := make([]domain.Rotulo, len(rotulos))
		for i, r := range rotulos {
			domainRotulos[i] = domain.Rotulo(r)
		}

		foto := domain.HydrateFoto(id, idCasamento, storageKey, urlPublica, ehFavorito, createdAt, domainRotulos)
		fotos = append(fotos, foto)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de fotos: %w", err)
	}

	return fotos, nil
}
func (r *PostgresFotoRepository) FindByID(ctx context.Context, fotoID uuid.UUID) (*domain.Foto, error) {
	// Query para buscar a foto e agregar seus rótulos
	sql := `
		SELECT
			f.id, f.id_casamento, f.storage_key, f.url_publica, f.eh_favorito, f.created_at,
			COALESCE(ARRAY_AGG(fr.nome_rotulo) FILTER (WHERE fr.nome_rotulo IS NOT NULL), '{}') as rotulos
		FROM fotos f
		LEFT JOIN fotos_rotulos fr ON f.id = fr.id_foto
		WHERE f.id = $1
		GROUP BY f.id;
	`
	row := r.db.QueryRow(ctx, sql, fotoID)

	var id, idCasamento uuid.UUID
	var storageKey, urlPublica string
	var ehFavorito bool
	var createdAt time.Time
	var rotulos []string

	if err := row.Scan(&id, &idCasamento, &storageKey, &urlPublica, &ehFavorito, &createdAt, &rotulos); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFotoNaoEncontrada // Crie este erro no seu pacote de domínio
		}
		return nil, fmt.Errorf("falha ao escanear foto por id: %w", err)
	}

	domainRotulos := make([]domain.Rotulo, len(rotulos))
	for i, r := range rotulos {
		domainRotulos[i] = domain.Rotulo(r)
	}

	foto := domain.HydrateFoto(id, idCasamento, storageKey, urlPublica, ehFavorito, createdAt, domainRotulos)
	return foto, nil
}
func (r *PostgresFotoRepository) Update(ctx context.Context, foto *domain.Foto) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação para update de foto: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Atualiza os campos na tabela principal 'fotos'.
	updateFotoSQL := `UPDATE fotos SET eh_favorito = $1 WHERE id = $2`
	if _, err := tx.Exec(ctx, updateFotoSQL, foto.EhFavorito(), foto.ID()); err != nil {
		return fmt.Errorf("falha ao atualizar dados da foto: %w", err)
	}

	// 2. Sincroniza os rótulos usando a estratégia "delete-then-insert".
	// Primeiro, remove todas as associações de rótulos existentes para esta foto.
	if _, err := tx.Exec(ctx, "DELETE FROM fotos_rotulos WHERE id_foto = $1", foto.ID()); err != nil {
		return fmt.Errorf("falha ao deletar rótulos antigos: %w", err)
	}

	// Segundo, insere a lista atual de rótulos do agregado.
	rotulos := foto.Rotulos()
	if len(rotulos) > 0 {
		rotulosRows := make([][]any, len(rotulos))
		for i, r := range rotulos {
			rotulosRows[i] = []any{foto.ID(), r}
		}

		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"fotos_rotulos"},
			[]string{"id_foto", "nome_rotulo"},
			pgx.CopyFromRows(rotulosRows),
		)
		if err != nil {
			return fmt.Errorf("falha no bulk insert de novos rótulos: %w", err)
		}
	}

	// 3. Confirma a transação.
	return tx.Commit(ctx)
}
func (r *PostgresFotoRepository) Delete(ctx context.Context, fotoID uuid.UUID) error {
	sql := `DELETE FROM fotos WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, sql, fotoID)
	if err != nil {
		return fmt.Errorf("falha ao deletar foto: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrFotoNaoEncontrada // Reutilizando nosso erro
	}
	return nil
}
