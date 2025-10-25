package infrastructure

import (
	"context"
	"fmt"
	"time"

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
		SET nome = $1, descricao = $2, eh_favorito = $3, foto_url = $4,
		    categoria = $5, detalhes_tipo = $6, detalhes_link_loja = $7
		WHERE id = $8;
	`
	detalhes := presente.Detalhes()
	_, err := r.db.Exec(ctx, sql,
		presente.Nome(),
		presente.Descricao(),
		presente.EhFavorito(),
		presente.FotoURL(),
		presente.Categoria(),
		detalhes.Tipo,
		detalhes.LinkDaLoja,
		presente.ID(),
	)
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

func (r *PostgresPresenteRepository) ListarTodosPorEvento(ctx context.Context, eventoID uuid.UUID) ([]*domain.PresenteComSelecao, error) {
	// Query unificada que retorna presentes com suas confirmações
	// Presentes fracionados aparecem múltiplas vezes (uma para cada palavra mágica diferente)
	sql := `
		WITH presentes_integrais AS (
			-- Presentes integrais: usa id_selecao diretamente da tabela presentes
			SELECT
				p.id, p.id_evento, p.nome, p.descricao, p.foto_url,
				p.status, p.categoria, p.eh_favorito,
				p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente,
				cg.chave_de_acesso,
				ps.data_da_selecao,
				CASE WHEN ps.id IS NOT NULL THEN 1 ELSE 0 END as quantidade_cotas
			FROM presentes p
			LEFT JOIN presentes_selecoes ps ON ps.id = p.id_selecao
			LEFT JOIN convidados_grupos cg ON cg.id = ps.id_grupo_de_convidados
			WHERE p.id_evento = $1 AND p.tipo = 'INTEGRAL'
		),
		presentes_fracionados_com_selecao AS (
			-- Presentes fracionados com pelo menos 1 cota selecionada
			-- Agrupa por palavra mágica (pode retornar múltiplas linhas por presente)
			SELECT
				p.id, p.id_evento, p.nome, p.descricao, p.foto_url,
				p.status, p.categoria, p.eh_favorito,
				p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente,
				cg.chave_de_acesso,
				ps.data_da_selecao,
				COUNT(cp.id)::int as quantidade_cotas
			FROM presentes p
			INNER JOIN cotas_de_presentes cp ON cp.id_presente = p.id AND cp.status != 'DISPONIVEL'
			INNER JOIN presentes_selecoes ps ON ps.id = cp.id_selecao
			INNER JOIN convidados_grupos cg ON cg.id = ps.id_grupo_de_convidados
			WHERE p.id_evento = $1 AND p.tipo = 'FRACIONADO'
			GROUP BY p.id, p.id_evento, p.nome, p.descricao, p.foto_url,
					p.status, p.categoria, p.eh_favorito,
					p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente,
					cg.chave_de_acesso, ps.data_da_selecao
		),
		presentes_fracionados_disponiveis AS (
			-- Presentes fracionados que ainda têm cotas disponíveis (aparecem 1x com selecao=null)
			SELECT
				p.id, p.id_evento, p.nome, p.descricao, p.foto_url,
				p.status, p.categoria, p.eh_favorito,
				p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente,
				NULL::VARCHAR as chave_de_acesso,
				NULL::TIMESTAMP WITH TIME ZONE as data_da_selecao,
				0 as quantidade_cotas
			FROM presentes p
			WHERE p.id_evento = $1
				AND p.tipo = 'FRACIONADO'
				AND p.status IN ('DISPONIVEL', 'PARCIALMENTE_SELECIONADO')
		)
		SELECT * FROM presentes_integrais
		UNION ALL
		SELECT * FROM presentes_fracionados_com_selecao
		UNION ALL
		SELECT * FROM presentes_fracionados_disponiveis
		ORDER BY nome ASC, data_da_selecao ASC NULLS FIRST;
	`

	rows, err := r.db.Query(ctx, sql, eventoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar todos os presentes com seleções: %w", err)
	}
	defer rows.Close()

	var resultado []*domain.PresenteComSelecao
	presentesMap := make(map[uuid.UUID]*domain.Presente)

	for rows.Next() {
		var d domain.DetalhesPresente
		var id, idCasamento uuid.UUID
		var nome, status, detalhesTipo, tipo string
		var ehFavorito bool
		var pDesc, pFotoURL, pLinkLoja, pCategoria, pChaveDeAcesso *string
		var pValorTotal *float64
		var pDataSelecao *time.Time
		var quantidadeCotas int

		if err := rows.Scan(
			&id, &idCasamento, &nome, &pDesc, &pFotoURL, &status, &pCategoria, &ehFavorito,
			&detalhesTipo, &pLinkLoja, &tipo, &pValorTotal,
			&pChaveDeAcesso, &pDataSelecao, &quantidadeCotas,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha de presente com seleção: %w", err)
		}

		// Construir presente
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

		// Verificar se já temos este presente no mapa
		presente, exists := presentesMap[id]
		if !exists {
			presente = domain.HydratePresente(id, idCasamento, nome, descricao, fotoURL, status, categoria, tipo, ehFavorito, d, pValorTotal, nil)
			presentesMap[id] = presente
		}

		// Criar registro de PresenteComSelecao
		pcs := &domain.PresenteComSelecao{
			Presente:        presente,
			ChaveDeAcesso:   pChaveDeAcesso,
			QuantidadeCotas: quantidadeCotas,
			DataSelecao:     pDataSelecao,
		}

		resultado = append(resultado, pcs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	// Carregar cotas para presentes fracionados
	if len(presentesMap) > 0 {
		err = r.loadCotas(ctx, presentesMap)
		if err != nil {
			return nil, err
		}
	}

	return resultado, nil
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

func (r *PostgresPresenteRepository) FindByID(ctx context.Context, userID, presenteID uuid.UUID) (*domain.Presente, error) {
	sql := `
		SELECT
			p.id, p.id_evento, p.nome, p.descricao, p.foto_url, p.status, p.categoria, p.eh_favorito,
			p.detalhes_tipo, p.detalhes_link_loja, p.tipo, p.valor_total_presente
		FROM presentes p
		JOIN eventos e ON p.id_evento = e.id
		WHERE p.id = $1 AND e.id_usuario = $2;
	`

	rows, err := r.db.Query(ctx, sql, presenteID, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar presente por ID: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, domain.ErrPresenteNaoEncontrado
	}

	presente, err := r.scanPresente(rows)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante leitura do presente: %w", err)
	}

	// Carregar cotas se for fracionado
	if presente.EhFracionado() {
		presentesMap := map[uuid.UUID]*domain.Presente{presente.ID(): presente}
		err = r.loadCotas(ctx, presentesMap)
		if err != nil {
			return nil, err
		}
		presente = presentesMap[presente.ID()]
	}

	return presente, nil
}

func (r *PostgresPresenteRepository) Delete(ctx context.Context, userID, presenteID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação para delete: %w", err)
	}
	defer tx.Rollback(ctx)

	// Verificar se o presente pertence ao usuário
	var exists bool
	checkSQL := `
		SELECT EXISTS(
			SELECT 1 FROM presentes p
			JOIN eventos e ON p.id_evento = e.id
			WHERE p.id = $1 AND e.id_usuario = $2
		)
	`
	if err := tx.QueryRow(ctx, checkSQL, presenteID, userID).Scan(&exists); err != nil {
		return fmt.Errorf("falha ao verificar propriedade do presente: %w", err)
	}

	if !exists {
		return domain.ErrPresenteNaoEncontrado
	}

	// Deletar cotas primeiro (FK constraint)
	if _, err := tx.Exec(ctx, "DELETE FROM cotas_de_presentes WHERE id_presente = $1", presenteID); err != nil {
		return fmt.Errorf("falha ao deletar cotas: %w", err)
	}

	// Deletar o presente
	cmdTag, err := tx.Exec(ctx, "DELETE FROM presentes WHERE id = $1", presenteID)
	if err != nil {
		return fmt.Errorf("falha ao deletar presente: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrPresenteNaoEncontrado
	}

	return tx.Commit(ctx)
}
