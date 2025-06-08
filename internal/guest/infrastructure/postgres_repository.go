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
		INSERT INTO grupos_de_convidados (id, id_casamento, chave_de_acesso)
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

func (r *PostgresGroupRepository) FindByAccessKey(ctx context.Context, accessKey string) (*domain.GrupoDeConvidados, error) {
	// Usamos LEFT JOIN para garantir que mesmo um grupo sem convidados (caso raro) seja retornado.
	sql := `
		SELECT 
			g.id, g.id_casamento, g.chave_de_acesso,
			c.id, c.nome, c.status_rsvp
		FROM grupos_de_convidados g
		LEFT JOIN convidados c ON g.id = c.id_grupo
		WHERE g.chave_de_acesso = $1;
	`

	rows, err := r.db.Query(ctx, sql, accessKey)
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
	defer br.Close()

	// Verifica se todas as operações no lote foram bem-sucedidas.
	for i := 0; i < len(group.Convidados()); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("falha ao executar update de rsvp de convidado no lote: %w", err)
		}
	}

	return tx.Commit(ctx)
}
