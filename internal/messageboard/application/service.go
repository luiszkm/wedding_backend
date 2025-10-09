// file: internal/messageboard/application/service.go
package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	guestDomain "github.com/luiszkm/wedding_backend/internal/guest/domain"
	"github.com/luiszkm/wedding_backend/internal/messageboard/domain"
)

type MessageBoardService struct {
	recadoRepo domain.RecadoRepository     // Dependência do repositório de recados
	guestRepo  guestDomain.GroupRepository // Dependência do repositório de outro contexto
	eventRepo  eventDomain.EventoRepository
}

type ModeracaoCommand struct {
	Status     *string
	EhFavorito *bool
}

// NewMessageBoardService cria uma nova instância do serviço.
func NewMessageBoardService(recadoRepo domain.RecadoRepository, guestRepo guestDomain.GroupRepository, eventRepo eventDomain.EventoRepository) *MessageBoardService {
	return &MessageBoardService{recadoRepo: recadoRepo, guestRepo: guestRepo, eventRepo: eventRepo}
}

// DeixarNovoRecado orquestra a criação de um novo recado.
func (s *MessageBoardService) DeixarNovoRecado(ctx context.Context, eventID uuid.UUID, chaveDeAcesso, nomeAutor, texto string) error {
	// 1. Valida a chave de acesso buscando pelo grupo de convidados.
	grupo, err := s.guestRepo.FindByAccessKey(ctx, eventID, chaveDeAcesso)
	if err != nil {
		// Retorna o erro original (pode ser "não encontrado" ou outro erro técnico).
		return fmt.Errorf("falha ao validar chave de acesso: %w", err)
	}

	// 2. Usa a fábrica do domínio para criar o novo recado.
	novoRecado, err := domain.NewRecado(grupo.IDCasamento(), grupo.ID(), nomeAutor, texto)
	if err != nil {
		return err // Retorna erros de validação de negócio (nome/texto vazio)
	}

	// 3. Persiste o novo recado.
	if err := s.recadoRepo.Save(ctx, novoRecado); err != nil {
		return fmt.Errorf("falha ao salvar novo recado: %w", err)
	}

	return nil
}

func (s *MessageBoardService) ListarRecadosParaAdmin(ctx context.Context, userID, idEvento uuid.UUID) ([]*domain.Recado, error) {
	// 1. AUTORIZAÇÃO: Verifica se o usuário é o dono do evento.
	if _, err := s.eventRepo.FindByID(ctx, userID, idEvento); err != nil {
		return nil, fmt.Errorf("permissão negada ou evento não encontrado: %w", err)
	}

	// 2. Se a permissão for válida, busca os recados.
	recados, err := s.recadoRepo.ListarPorEvento(ctx, idEvento) // Renomeie ListarPorCasamento
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista de recados para admin: %w", err)
	}
	return recados, nil
}
func (s *MessageBoardService) ModerarRecado(ctx context.Context, userID, recadoID uuid.UUID, cmd ModeracaoCommand) (*domain.Recado, error) {
	recado, err := s.recadoRepo.FindByID(ctx, userID, recadoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar recado para moderação: %w", err)
	}

	if cmd.Status != nil {
		status := *cmd.Status
		if status == domain.StatusAprovado {
			recado.Aprovar()
		} else if status == domain.StatusRejeitado {
			recado.Rejeitar()
		} else {
			return nil, domain.ErrStatusInvalidoParaModeracao
		}
	}

	if cmd.EhFavorito != nil {
		recado.DefinirFavorito(*cmd.EhFavorito)
	}

	if err := s.recadoRepo.Update(ctx, userID, recado); err != nil {
		return nil, fmt.Errorf("falha ao salvar moderação do recado: %w", err)
	}

	return recado, nil
}
func (s *MessageBoardService) ListarRecadosPublicos(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Recado, error) {
	recados, err := s.recadoRepo.ListarAprovadosPorCasamento(ctx, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista pública de recados: %w", err)
	}
	return recados, nil
}
