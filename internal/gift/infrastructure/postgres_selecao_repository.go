// file: internal/gift/infrastructure/postgres_selecao_repository.go
package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/wedding_backend/internal/gift/domain"
	guestDomain "github.com/luiszkm/wedding_backend/internal/guest/domain"
)

type PostgresSelecaoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSelecaoRepository(db *pgxpool.Pool) domain.SelecaoRepository {
	return &PostgresSelecaoRepository{db: db}
}

func (r *PostgresSelecaoRepository) SalvarSelecao(ctx context.Context, chaveDeAcesso string, idsDosPresentes []uuid.UUID) (*domain.Selecao, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Obter o ID do grupo de convidados a partir da chave de acesso.
	var grupoID uuid.UUID
	err = tx.QueryRow(ctx, "SELECT id FROM grupos_de_convidados WHERE chave_de_acesso = $1", chaveDeAcesso).Scan(&grupoID)
	if err != nil {
		// Se a chave não existe, retorna o erro de grupo não encontrado.
		return nil, guestDomain.ErrGrupoNaoEncontrado
	}

	// 2. Verificar o status dos presentes desejados com um bloqueio de linha (FOR UPDATE).
	// Isso impede que outra transação modifique estas linhas até a nossa terminar.
	rows, err := tx.Query(ctx, "SELECT id, nome, status FROM presentes WHERE id = ANY($1) FOR UPDATE", idsDosPresentes)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar presentes para seleção: %w", err)
	}
	defer rows.Close()

	var presentesConflitantes []uuid.UUID
	var presentesDisponiveis []domain.PresenteConfirmado
	presentesVerificados := make(map[uuid.UUID]bool)

	for rows.Next() {
		var id uuid.UUID
		var nome, status string
		if err := rows.Scan(&id, &nome, &status); err != nil {
			return nil, fmt.Errorf("falha ao escanear presente: %w", err)
		}
		if status != "DISPONIVEL" {
			presentesConflitantes = append(presentesConflitantes, id)
		}
		presentesDisponiveis = append(presentesDisponiveis, domain.PresenteConfirmado{ID: id, Nome: nome})
		presentesVerificados[id] = true
	}

	// Se algum dos presentes não foi encontrado no banco, é um erro.
	if len(presentesVerificados) != len(idsDosPresentes) {
		return nil, errors.New("um ou mais IDs de presentes são inválidos")
	}

	// Se encontramos presentes já selecionados, retornamos o erro de conflito.
	if len(presentesConflitantes) > 0 {
		return nil, &domain.ErrPresentesConflitantes{PresentesIDs: presentesConflitantes}
	}

	// 3. Criar o registro da seleção.
	var idCasamento, selecaoID uuid.UUID
	// Assume-se que todos os presentes são do mesmo casamento. Pega o ID do primeiro.
	err = tx.QueryRow(ctx, "SELECT id_casamento FROM presentes WHERE id = $1", idsDosPresentes[0]).Scan(&idCasamento)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter id do casamento: %w", err)
	}

	sqlInsertSelecao := "INSERT INTO selecoes_de_presentes (id_casamento, id_grupo_de_convidados) VALUES ($1, $2) RETURNING id"
	err = tx.QueryRow(ctx, sqlInsertSelecao, idCasamento, grupoID).Scan(&selecaoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar registro de seleção: %w", err)
	}

	// 4. Atualizar os presentes para o status 'SELECIONADO' e ligá-los à seleção.
	sqlUpdatePresentes := "UPDATE presentes SET status = 'SELECIONADO', id_selecao = $1 WHERE id = ANY($2)"
	if _, err := tx.Exec(ctx, sqlUpdatePresentes, selecaoID, idsDosPresentes); err != nil {
		return nil, fmt.Errorf("falha ao atualizar status dos presentes: %w", err)
	}

	// 5. Se tudo deu certo, confirma a transação.
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("falha ao commitar transação: %w", err)
	}

	selecao := domain.NewSelecao(idCasamento, grupoID, presentesDisponiveis)
	// Sobrescreve o ID gerado com o do banco de dados para consistência.
	selecao = domain.HydrateSelecao(selecaoID, idCasamento, grupoID, presentesDisponiveis, selecao.DataDaSelecao())

	return selecao, nil
}
