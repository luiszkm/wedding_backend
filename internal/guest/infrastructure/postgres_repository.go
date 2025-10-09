// file: internal/guest/infrastructure/postgres_repository.go
package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/guest/domain"
)

type PostgresGroupRepository struct {
	db *pgxpool.Pool
}

func NewPostgresGroupRepository(db *pgxpool.Pool) domain.GroupRepository {
	return &PostgresGroupRepository{db: db}
}

func (r *PostgresGroupRepository) Save(ctx context.Context, group *domain.GrupoDeConvidados) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	sqlGrupo := `
		INSERT INTO convidados_grupos (id, id_evento, chave_de_acesso)
		VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(ctx, sqlGrupo, group.ID(), group.IDCasamento(), group.ChaveDeAcesso())
	if err != nil {
		log.Printf("!!! ERRO DO BANCO DE DADOS: %v", err)
		return fmt.Errorf("falha ao inserir grupo de convidados: %w", err)
	}

	rows := make([][]any, len(group.Convidados()))
	for i, c := range group.Convidados() {
		rows[i] = []any{c.ID(), group.ID(), c.Nome()}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"convidados"},
		[]string{"id", "id_grupo", "nome"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("falha ao inserir convidados em lote: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresGroupRepository) FindByAccessKey(ctx context.Context, eventID uuid.UUID, accessKey string) (*domain.GrupoDeConvidados, error) {
	// Usamos LEFT JOIN para garantir que mesmo um grupo sem convidados (caso raro) seja retornado.
	// Filtramos por id_evento E chave_de_acesso para evitar ambiguidade
	sql := `
		SELECT
			g.id, g.id_evento, g.chave_de_acesso,
			c.id, c.nome, c.status_rsvp
		FROM convidados_grupos g
		LEFT JOIN convidados c ON g.id = c.id_grupo
		WHERE g.id_evento = $1 AND g.chave_de_acesso = $2;
	`

	rows, err := r.db.Query(ctx, sql, eventID, accessKey)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar grupo por chave de acesso: %w", err)
	}
	defer rows.Close()

	var grupo *domain.GrupoDeConvidados
	var convidados []*domain.Convidado

	for rows.Next() {
		var grupoID, idCasamento, convidadoID uuid.UUID
		var chaveDeAcesso, nomeConvidado, statusRSVP string

		// Usamos ponteiros para os campos de convidados para detectar quando eles são NULL
		// (no caso de um grupo sem convidados).
		var pConvidadoID *uuid.UUID
		var pNomeConvidado, pStatusRSVP *string

		if err := rows.Scan(
			&grupoID, &idCasamento, &chaveDeAcesso,
			&pConvidadoID, &pNomeConvidado, &pStatusRSVP,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha da consulta: %w", err)
		}

		// Se o grupo ainda não foi criado, criamo-lo com os dados da primeira linha.
		if grupo == nil {
			grupo = domain.HydrateGroup(grupoID, idCasamento, chaveDeAcesso, nil)
		}

		// Se houver dados de convidado na linha, criamos o objeto convidado.
		if pConvidadoID != nil {
			convidadoID = *pConvidadoID
			nomeConvidado = *pNomeConvidado
			statusRSVP = *pStatusRSVP
			convidado := domain.HydrateConvidado(convidadoID, nomeConvidado, statusRSVP)
			convidados = append(convidados, convidado)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	// Se nenhuma linha foi processada, o grupo não existe.
	if grupo == nil {
		return nil, domain.ErrGrupoNaoEncontrado
	}

	// "Hidratamos" o agregado com sua lista de convidados.
	grupo = domain.HydrateGroup(grupo.ID(), grupo.IDCasamento(), grupo.ChaveDeAcesso(), convidados)

	return grupo, nil
}

func (r *PostgresGroupRepository) FindByID(ctx context.Context, userID, groupID uuid.UUID) (*domain.GrupoDeConvidados, error) {
	sql := `
		SELECT 
			g.id, g.id_evento, g.chave_de_acesso,
			c.id, c.nome, c.status_rsvp
		FROM convidados_grupos g
		JOIN eventos e ON g.id_evento = e.id
		LEFT JOIN convidados c ON g.id = c.id_grupo
		WHERE g.id = $1 AND e.id_usuario = $2;
	`
	rows, err := r.db.Query(ctx, sql, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar grupo por id: %w", err)
	}
	defer rows.Close()

	var grupo *domain.GrupoDeConvidados
	var convidados []*domain.Convidado

	for rows.Next() {
		var grupoID, idCasamento, convidadoID uuid.UUID
		var chaveDeAcesso, nomeConvidado, statusRSVP string
		var pConvidadoID *uuid.UUID
		var pNomeConvidado, pStatusRSVP *string

		if err := rows.Scan(
			&grupoID, &idCasamento, &chaveDeAcesso,
			&pConvidadoID, &pNomeConvidado, &pStatusRSVP,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha da consulta de grupo por id: %w", err)
		}

		if grupo == nil {
			grupo = domain.HydrateGroup(grupoID, idCasamento, chaveDeAcesso, nil)
		}

		if pConvidadoID != nil {
			convidadoID = *pConvidadoID
			nomeConvidado = *pNomeConvidado
			statusRSVP = *pStatusRSVP
			convidado := domain.HydrateConvidado(convidadoID, nomeConvidado, statusRSVP)
			convidados = append(convidados, convidado)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	if grupo == nil {
		return nil, domain.ErrGrupoNaoEncontrado
	}

	// "Hidratamos" o agregado com sua lista de convidados
	grupo = domain.HydrateGroup(grupo.ID(), grupo.IDCasamento(), grupo.ChaveDeAcesso(), convidados)

	return grupo, nil
}

func (r *PostgresGroupRepository) Update(ctx context.Context, userID uuid.UUID, group *domain.GrupoDeConvidados) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação para update: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Atualiza os dados do grupo principal (chave de acesso e timestamp)
	updateGroupSQL := `
		UPDATE convidados_grupos SET chave_de_acesso = $1, updated_at = $2 
		WHERE id = $3 AND id_evento IN (SELECT id FROM eventos WHERE id_usuario = $4)
	`
	cmdTag, err := tx.Exec(ctx, updateGroupSQL, group.ChaveDeAcesso(), group.UpdatedAt(), group.ID(), userID)
	if err != nil {
		return fmt.Errorf("falha ao atualizar dados do grupo: %w", err)
	}
	// Se nenhuma linha foi afetada, ou o grupo não existe ou o usuário não tem permissão.
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrGrupoNaoEncontrado
	}
	// 2. Remove TODOS os convidados antigos associados a este grupo.
	// Esta é a parte "delete" da estratégia "delete-then-insert".
	if _, err := tx.Exec(ctx, "DELETE FROM convidados WHERE id_grupo = $1", group.ID()); err != nil {
		return fmt.Errorf("falha ao deletar convidados antigos: %w", err)
	}

	// 3. Insere a NOVA lista de convidados em lote.
	// Esta é a parte "insert" da estratégia.
	if len(group.Convidados()) > 0 {
		rows := make([][]any, len(group.Convidados()))
		for i, c := range group.Convidados() {
			rows[i] = []any{c.ID(), group.ID(), c.Nome(), c.StatusRSVP()}
		}

		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"convidados"},
			[]string{"id", "id_grupo", "nome", "status_rsvp"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return fmt.Errorf("falha ao inserir nova lista de convidados: %w", err)
		}
	}

	// 4. Confirma a transação.
	return tx.Commit(ctx)
}

func (r *PostgresGroupRepository) FindAllByEventID(ctx context.Context, userID, eventID uuid.UUID, statusFilter string) ([]*domain.GrupoDeConvidados, error) {
	baseSQL := `
		SELECT 
			g.id, g.id_evento, g.chave_de_acesso,
			c.id, c.nome, c.status_rsvp
		FROM convidados_grupos g
		JOIN eventos e ON g.id_evento = e.id
		LEFT JOIN convidados c ON g.id = c.id_grupo
		WHERE g.id_evento = $1 AND e.id_usuario = $2`

	args := []interface{}{eventID, userID}

	if statusFilter != "" {
		baseSQL += " AND c.status_rsvp = $3"
		args = append(args, statusFilter)
	}

	baseSQL += " ORDER BY g.created_at DESC"

	rows, err := r.db.Query(ctx, baseSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar grupos por evento: %w", err)
	}
	defer rows.Close()

	gruposMap := make(map[uuid.UUID]*domain.GrupoDeConvidados)
	var gruposOrdenados []*domain.GrupoDeConvidados

	for rows.Next() {
		var grupoID, idEvento, convidadoID uuid.UUID
		var chaveDeAcesso, nomeConvidado, statusRSVP string
		var pConvidadoID *uuid.UUID
		var pNomeConvidado, pStatusRSVP *string

		if err := rows.Scan(
			&grupoID, &idEvento, &chaveDeAcesso,
			&pConvidadoID, &pNomeConvidado, &pStatusRSVP,
		); err != nil {
			return nil, fmt.Errorf("falha ao escanear linha da consulta por evento: %w", err)
		}

		grupo, existe := gruposMap[grupoID]
		if !existe {
			grupo = domain.HydrateGroup(grupoID, idEvento, chaveDeAcesso, nil)
			gruposMap[grupoID] = grupo
			gruposOrdenados = append(gruposOrdenados, grupo)
		}

		if pConvidadoID != nil {
			convidadoID = *pConvidadoID
			nomeConvidado = *pNomeConvidado
			statusRSVP = *pStatusRSVP
			convidado := domain.HydrateConvidado(convidadoID, nomeConvidado, statusRSVP)

			// Precisa recriar o grupo com os convidados atualizados
			convidadosAtuais := grupo.Convidados()
			convidadosAtualizados := append(convidadosAtuais, convidado)
			grupoAtualizado := domain.HydrateGroup(grupo.ID(), grupo.IDCasamento(), grupo.ChaveDeAcesso(), convidadosAtualizados)
			gruposMap[grupoID] = grupoAtualizado

			// Atualizar na lista ordenada também
			for i, g := range gruposOrdenados {
				if g.ID() == grupoID {
					gruposOrdenados[i] = grupoAtualizado
					break
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	return gruposOrdenados, nil
}

func (r *PostgresGroupRepository) Delete(ctx context.Context, userID, groupID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação para delete: %w", err)
	}
	defer tx.Rollback(ctx)

	// Primeiro, verificar se o grupo pertence ao usuário
	var exists bool
	checkSQL := `
		SELECT EXISTS(
			SELECT 1 FROM convidados_grupos g
			JOIN eventos e ON g.id_evento = e.id
			WHERE g.id = $1 AND e.id_usuario = $2
		)
	`
	if err := tx.QueryRow(ctx, checkSQL, groupID, userID).Scan(&exists); err != nil {
		return fmt.Errorf("falha ao verificar propriedade do grupo: %w", err)
	}

	if !exists {
		return domain.ErrGrupoNaoEncontrado
	}

	// Deletar convidados primeiro (FK constraint)
	if _, err := tx.Exec(ctx, "DELETE FROM convidados WHERE id_grupo = $1", groupID); err != nil {
		return fmt.Errorf("falha ao deletar convidados: %w", err)
	}

	// Deletar o grupo
	cmdTag, err := tx.Exec(ctx, "DELETE FROM convidados_grupos WHERE id = $1", groupID)
	if err != nil {
		return fmt.Errorf("falha ao deletar grupo: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrGrupoNaoEncontrado
	}

	return tx.Commit(ctx)
}

func (r *PostgresGroupRepository) GetRSVPStats(ctx context.Context, userID, eventID uuid.UUID) (*domain.RSVPStats, error) {
	sql := `
		SELECT 
			COUNT(DISTINCT g.id) as total_grupos,
			COUNT(c.id) as total_convidados,
			COUNT(CASE WHEN c.status_rsvp = 'CONFIRMADO' THEN 1 END) as confirmados,
			COUNT(CASE WHEN c.status_rsvp = 'RECUSADO' THEN 1 END) as recusados,
			COUNT(CASE WHEN c.status_rsvp = 'PENDENTE' THEN 1 END) as pendentes
		FROM convidados_grupos g
		JOIN eventos e ON g.id_evento = e.id
		LEFT JOIN convidados c ON g.id = c.id_grupo
		WHERE g.id_evento = $1 AND e.id_usuario = $2
	`

	var totalGrupos, totalConvidados, confirmados, recusados, pendentes int

	err := r.db.QueryRow(ctx, sql, eventID, userID).Scan(
		&totalGrupos, &totalConvidados, &confirmados, &recusados, &pendentes,
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter estatísticas RSVP: %w", err)
	}

	stats := &domain.RSVPStats{
		TotalGrupos:           totalGrupos,
		TotalConvidados:       totalConvidados,
		ConvidadosConfirmados: confirmados,
		ConvidadosRecusados:   recusados,
		ConvidadosPendentes:   pendentes,
	}

	// Calcular percentuais
	if totalConvidados > 0 {
		stats.PercentualConfirmado = float64(confirmados) / float64(totalConvidados) * 100
		stats.PercentualRecusado = float64(recusados) / float64(totalConvidados) * 100
		stats.PercentualPendente = float64(pendentes) / float64(totalConvidados) * 100
	}

	return stats, nil
}

func (r *PostgresGroupRepository) UpdateRSVP(ctx context.Context, group *domain.GrupoDeConvidados) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação para update de rsvp: %w", err)
	}
	defer tx.Rollback(ctx)

	// Prepara um lote para atualizar todos os convidados de uma vez.
	batch := &pgx.Batch{}
	updateGuestSQL := "UPDATE convidados SET status_rsvp = $1 WHERE id = $2 AND id_grupo = $3"
	for _, convidado := range group.Convidados() {
		batch.Queue(updateGuestSQL, convidado.StatusRSVP(), convidado.ID(), group.ID())
	}

	br := tx.SendBatch(ctx, batch)

	// Verifica se todas as operações no lote foram bem-sucedidas.
	for i := 0; i < len(group.Convidados()); i++ {
		cmdTag, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("falha ao executar update de rsvp de convidado no lote: %w", err)
		}
		// Verifica se a linha foi realmente atualizada
		if cmdTag.RowsAffected() == 0 {
			br.Close()
			return fmt.Errorf("convidado não foi atualizado: possivelmente não pertence ao grupo especificado")
		}
	}

	// Fechar o batch reader ANTES de fazer commit
	if err := br.Close(); err != nil {
		return fmt.Errorf("falha ao fechar batch reader: %w", err)
	}

	return tx.Commit(ctx)
}
