package infrastructure

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	if presente.EhFracionado() {
		return r.SaveWithCotas(ctx, presente)
	}

	sql := `
		INSERT INTO presentes (
			id, id_evento, nome, descricao, eh_favorito, status, 
			detalhes_tipo, detalhes_link_loja, foto_url, categoria,
			tipo, valor_total_presente
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
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
		presente.Tipo(),
		presente.ValorTotal(),
	)

	if err != nil {
		return fmt.Errorf("falha ao inserir presente no banco de dados: %w", err)
	}
	return nil
}

func (r *PostgresPresenteRepository) SaveWithCotas(ctx context.Context, presente *domain.Presente) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// Inserir o presente
	sqlPresente := `
		INSERT INTO presentes (
			id, id_evento, nome, descricao, eh_favorito, status, 
			detalhes_tipo, detalhes_link_loja, foto_url, categoria,
			tipo, valor_total_presente
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`
	detalhes := presente.Detalhes()
	_, err = tx.Exec(ctx, sqlPresente,
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
		presente.Tipo(),
		presente.ValorTotal(),
	)

	if err != nil {
		return fmt.Errorf("falha ao inserir presente: %w", err)
	}

	// Inserir as cotas
	sqlCota := `
		INSERT INTO cotas_de_presentes (
			id, id_presente, numero_cota, valor_cota, status, id_selecao
		) VALUES ($1, $2, $3, $4, $5, $6);
	`

	for _, cota := range presente.Cotas() {
		_, err = tx.Exec(ctx, sqlCota,
			cota.ID(),
			cota.IDPresente(),
			cota.NumeroCota(),
			cota.ValorCota(),
			cota.Status(),
			cota.IDSelecao(),
		)
		if err != nil {
			return fmt.Errorf("falha ao inserir cota %d: %w", cota.NumeroCota(), err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("falha ao confirmar transação: %w", err)
	}

	return nil
}

func (r *PostgresPresenteRepository) Update(ctx context.Context, presente *domain.Presente) error {
	sql := `
		UPDATE presentes 
		SET status = $1
		WHERE id = $2;
	`
	_, err := r.db.Exec(ctx, sql, presente.Status(), presente.ID())
	if err != nil {
		return fmt.Errorf("falha ao atualizar presente: %w", err)
	}
	return nil
}

func (r *PostgresPresenteRepository) UpdateCotasStatus(ctx context.Context, cotaIDs []uuid.UUID, status string, selecaoID *uuid.UUID) error {
	if len(cotaIDs) == 0 {
		return nil
	}

	sql := `
		UPDATE cotas_de_presentes 
		SET status = $1, id_selecao = $2
		WHERE id = ANY($3);
	`
	_, err := r.db.Exec(ctx, sql, status, selecaoID, cotaIDs)
	if err != nil {
		return fmt.Errorf("falha ao atualizar status das cotas: %w", err)
	}
	return nil
}

func (r *PostgresPresenteRepository) FindByIDs(ctx context.Context, presenteIDs []uuid.UUID) ([]*domain.Presente, error) {
	if len(presenteIDs) == 0 {
		return []*domain.Presente{}, nil
	}

	sql := `
		SELECT
			p.id, p.id_evento, p.nome, p.descricao, p.foto_url, p.status, p.categoria, p.eh_favorito,
			p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente
		FROM presentes p
		WHERE p.id = ANY($1)
		ORDER BY p.nome ASC;
	`

	rows, err := r.db.Query(ctx, sql, presenteIDs)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar presentes por IDs: %w", err)
	}
	defer rows.Close()

	presentesMap := make(map[uuid.UUID]*domain.Presente)

	for rows.Next() {
		presente, err := r.scanPresente(rows)
		if err != nil {
			return nil, err
		}
		presentesMap[presente.ID()] = presente
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de presentes: %w", err)
	}

	// Carregar cotas para presentes fracionados
	err = r.loadCotas(ctx, presentesMap)
	if err != nil {
		return nil, err
	}

	// Converter map para slice
	presentes := make([]*domain.Presente, 0, len(presentesMap))
	for _, presente := range presentesMap {
		presentes = append(presentes, presente)
	}

	return presentes, nil
}

func (r *PostgresPresenteRepository) ListarDisponiveisPorCasamento(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Presente, error) {
	sql := `
		SELECT
			p.id, p.id_evento, p.nome, p.descricao, p.foto_url, p.status, p.categoria, p.eh_favorito,
			p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente
		FROM presentes p
		WHERE p.id_evento = $1 AND p.status IN ('DISPONIVEL', 'PARCIALMENTE_SELECIONADO')
		ORDER BY p.eh_favorito DESC, p.nome ASC;
	`

	rows, err := r.db.Query(ctx, sql, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar presentes disponíveis: %w", err)
	}
	defer rows.Close()

	presentesMap := make(map[uuid.UUID]*domain.Presente)

	for rows.Next() {
		presente, err := r.scanPresente(rows)
		if err != nil {
			return nil, err
		}
		presentesMap[presente.ID()] = presente
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas de presentes: %w", err)
	}

	// Carregar cotas para presentes fracionados
	err = r.loadCotas(ctx, presentesMap)
	if err != nil {
		return nil, err
	}

	// Converter map para slice mantendo a ordem
	presentes := make([]*domain.Presente, 0, len(presentesMap))
	for _, presente := range presentesMap {
		presentes = append(presentes, presente)
	}

	return presentes, nil
}

func (r *PostgresPresenteRepository) scanPresente(rows pgx.Rows) (*domain.Presente, error) {
	var d domain.DetalhesPresente
	var id, idCasamento uuid.UUID
	var nome, status, detalhesTipo, tipo string
	var ehFavorito bool
	var pDesc, pFotoURL, pLinkLoja, pCategoria *string
	var pValorTotal *float64

	if err := rows.Scan(
		&id, &idCasamento, &nome, &pDesc, &pFotoURL, &status, &pCategoria, &ehFavorito,
		&detalhesTipo, &pLinkLoja, &tipo, &pValorTotal,
	); err != nil {
		return nil, fmt.Errorf("falha ao escanear linha de presente: %w", err)
	}

	// Atribui valores dos ponteiros, se não forem nulos
	descricao := ""
	if pDesc != nil {
		descricao = *pDesc
	}

	fotoURL := ""
	if pFotoURL != nil {
		fotoURL = *pFotoURL
	}

	detalhesLinkLoja := ""
	if pLinkLoja != nil {
		detalhesLinkLoja = *pLinkLoja
	}

	categoria := ""
	if pCategoria != nil {
		categoria = *pCategoria
	}

	d.Tipo = detalhesTipo
	d.LinkDaLoja = detalhesLinkLoja

	return domain.HydratePresente(id, idCasamento, nome, descricao, fotoURL, status, categoria, tipo, ehFavorito, d, pValorTotal, nil), nil
}

func (r *PostgresPresenteRepository) loadCotas(ctx context.Context, presentesMap map[uuid.UUID]*domain.Presente) error {
	// Identificar presentes fracionados
	fracionadoIDs := make([]uuid.UUID, 0)
	for _, presente := range presentesMap {
		if presente.EhFracionado() {
			fracionadoIDs = append(fracionadoIDs, presente.ID())
		}
	}

	if len(fracionadoIDs) == 0 {
		return nil
	}

	// Buscar cotas
	sqlCotas := `
		SELECT id, id_presente, numero_cota, valor_cota, status, id_selecao
		FROM cotas_de_presentes
		WHERE id_presente = ANY($1)
		ORDER BY id_presente, numero_cota;
	`

	rows, err := r.db.Query(ctx, sqlCotas, fracionadoIDs)
	if err != nil {
		return fmt.Errorf("falha ao buscar cotas: %w", err)
	}
	defer rows.Close()

	cotasPorPresente := make(map[uuid.UUID][]*domain.Cota)

	for rows.Next() {
		var id, idPresente uuid.UUID
		var numeroCota int
		var valorCota float64
		var status string
		var pIDSelecao *uuid.UUID

		if err := rows.Scan(&id, &idPresente, &numeroCota, &valorCota, &status, &pIDSelecao); err != nil {
			return fmt.Errorf("falha ao escanear cota: %w", err)
		}

		cota := domain.HydrateCota(id, idPresente, numeroCota, valorCota, status, pIDSelecao)
		cotasPorPresente[idPresente] = append(cotasPorPresente[idPresente], cota)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("erro durante iteração das cotas: %w", err)
	}

	// Associar cotas aos presentes
	for presenteID, cotas := range cotasPorPresente {
		if presente, exists := presentesMap[presenteID]; exists {
			// Recriar presente com cotas carregadas
			presentesMap[presenteID] = domain.HydratePresente(
				presente.ID(),
				presente.IDCasamento(),
				presente.Nome(),
				presente.Descricao(),
				presente.FotoURL(),
				presente.Status(),
				presente.Categoria(),
				presente.Tipo(),
				presente.EhFavorito(),
				presente.Detalhes(),
				presente.ValorTotal(),
				cotas,
			)
		}
	}

	return nil
}
