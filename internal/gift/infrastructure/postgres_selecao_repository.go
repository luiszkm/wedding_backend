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

func (r *PostgresSelecaoRepository) SalvarSelecao(ctx context.Context, chaveDeAcesso string, quantidades map[uuid.UUID]int) (*domain.Selecao, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Obter o ID do grupo de convidados a partir da chave de acesso.
	var grupoID uuid.UUID
	err = tx.QueryRow(ctx, "SELECT id FROM convidados_grupos WHERE chave_de_acesso = $1", chaveDeAcesso).Scan(&grupoID)
	if err != nil {
		// Se a chave não existe, retorna o erro de grupo não encontrado.
		return nil, guestDomain.ErrGrupoNaoEncontrado
	}

	// Extrair IDs dos presentes
	presenteIDs := make([]uuid.UUID, 0, len(quantidades))
	for id := range quantidades {
		presenteIDs = append(presenteIDs, id)
	}

	// 2. Verificar os presentes e seus tipos com bloqueio de linha (FOR UPDATE).
	sqlBuscarPresentes := `
		SELECT id, nome, status, tipo, valor_total_presente, id_evento
		FROM presentes
		WHERE id = ANY($1)
		FOR UPDATE
	`
	rows, err := tx.Query(ctx, sqlBuscarPresentes, presenteIDs)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar presentes para seleção: %w", err)
	}
	defer rows.Close()

	type presenteInfo struct {
		id          uuid.UUID
		nome        string
		status      string
		tipo        string
		valorTotal  *float64
		idCasamento uuid.UUID
	}

	presentes := make(map[uuid.UUID]presenteInfo)
	var idCasamento uuid.UUID

	for rows.Next() {
		var p presenteInfo
		if err := rows.Scan(&p.id, &p.nome, &p.status, &p.tipo, &p.valorTotal, &p.idCasamento); err != nil {
			return nil, fmt.Errorf("falha ao escanear presente: %w", err)
		}
		presentes[p.id] = p
		if idCasamento == uuid.Nil {
			idCasamento = p.idCasamento
		}
	}

	// Verificar se todos os presentes foram encontrados
	if len(presentes) != len(presenteIDs) {
		return nil, errors.New("um ou mais IDs de presentes são inválidos")
	}

	// 3. Criar o registro da seleção.
	selecaoID := uuid.New()
	sqlInsertSelecao := "INSERT INTO presentes_selecoes (id, id_evento, id_grupo_de_convidados) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, sqlInsertSelecao, selecaoID, idCasamento, grupoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar registro de seleção: %w", err)
	}

	// 4. Processar cada presente baseado no tipo
	var presentesConflitantes []uuid.UUID
	presentesConfirmados := make([]domain.PresenteConfirmado, 0, len(presentes))

	for presenteID, quantidade := range quantidades {
		presente := presentes[presenteID]

		if presente.tipo == "INTEGRAL" {
			// Presente integral: deve estar disponível e quantidade = 1
			if quantidade != 1 {
				return nil, fmt.Errorf("presente integral %s deve ter quantidade 1", presente.nome)
			}
			if presente.status != "DISPONIVEL" {
				presentesConflitantes = append(presentesConflitantes, presenteID)
				continue
			}

			// Atualizar presente integral
			sqlUpdateIntegral := "UPDATE presentes SET status = 'SELECIONADO', id_selecao = $1 WHERE id = $2"
			if _, err := tx.Exec(ctx, sqlUpdateIntegral, selecaoID, presenteID); err != nil {
				return nil, fmt.Errorf("falha ao atualizar presente integral: %w", err)
			}

			presentesConfirmados = append(presentesConfirmados, domain.PresenteConfirmado{
				ID:         presenteID,
				Nome:       presente.nome,
				Quantidade: 1,
				ValorCota:  nil,
			})

		} else if presente.tipo == "FRACIONADO" {
			// Presente fracionado: selecionar cotas disponíveis
			sqlBuscarCotas := `
				SELECT id, valor_cota
				FROM cotas_de_presentes
				WHERE id_presente = $1 AND status = 'DISPONIVEL'
				ORDER BY numero_cota
				LIMIT $2
				FOR UPDATE
			`
			cotasRows, err := tx.Query(ctx, sqlBuscarCotas, presenteID, quantidade)
			if err != nil {
				return nil, fmt.Errorf("falha ao buscar cotas: %w", err)
			}

			cotasIDs := make([]uuid.UUID, 0, quantidade)
			var valorCota float64

			for cotasRows.Next() {
				var cotaID uuid.UUID
				if err := cotasRows.Scan(&cotaID, &valorCota); err != nil {
					cotasRows.Close()
					return nil, fmt.Errorf("falha ao escanear cota: %w", err)
				}
				cotasIDs = append(cotasIDs, cotaID)
			}
			cotasRows.Close()

			// Verificar se há cotas suficientes
			if len(cotasIDs) < quantidade {
				return nil, fmt.Errorf("presente %s tem apenas %d cotas disponíveis, solicitado %d", presente.nome, len(cotasIDs), quantidade)
			}

			// Atualizar cotas para SELECIONADO
			sqlUpdateCotas := "UPDATE cotas_de_presentes SET status = 'SELECIONADO', id_selecao = $1 WHERE id = ANY($2)"
			if _, err := tx.Exec(ctx, sqlUpdateCotas, selecaoID, cotasIDs); err != nil {
				return nil, fmt.Errorf("falha ao atualizar cotas: %w", err)
			}

			// Verificar quantas cotas ainda estão disponíveis
			var cotasRestantes int
			sqlContarCotas := "SELECT COUNT(*) FROM cotas_de_presentes WHERE id_presente = $1 AND status = 'DISPONIVEL'"
			if err := tx.QueryRow(ctx, sqlContarCotas, presenteID).Scan(&cotasRestantes); err != nil {
				return nil, fmt.Errorf("falha ao contar cotas restantes: %w", err)
			}

			// Atualizar status do presente
			var novoStatus string
			if cotasRestantes == 0 {
				novoStatus = "SELECIONADO"
			} else {
				novoStatus = "PARCIALMENTE_SELECIONADO"
			}

			sqlUpdatePresenteFracionado := "UPDATE presentes SET status = $1 WHERE id = $2"
			if _, err := tx.Exec(ctx, sqlUpdatePresenteFracionado, novoStatus, presenteID); err != nil {
				return nil, fmt.Errorf("falha ao atualizar status do presente fracionado: %w", err)
			}

			presentesConfirmados = append(presentesConfirmados, domain.PresenteConfirmado{
				ID:         presenteID,
				Nome:       presente.nome,
				Quantidade: quantidade,
				ValorCota:  &valorCota,
			})
		}
	}

	// Se encontramos presentes conflitantes, retornar erro
	if len(presentesConflitantes) > 0 {
		return nil, &domain.ErrPresentesConflitantes{PresentesIDs: presentesConflitantes}
	}

	// 5. Se tudo deu certo, confirma a transação.
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("falha ao commitar transação: %w", err)
	}

	selecao := domain.NewSelecao(idCasamento, grupoID, presentesConfirmados)
	// Sobrescreve o ID gerado com o do banco de dados para consistência.
	selecao = domain.HydrateSelecao(selecaoID, idCasamento, grupoID, presentesConfirmados, selecao.DataDaSelecao())

	return selecao, nil
}
