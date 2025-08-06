package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/communication/domain"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
)

type CommunicationService struct {
	comunicadoRepo domain.ComunicadoRepository
	eventRepo      eventDomain.EventoRepository
}

func NewCommunicationService(comunicadoRepo domain.ComunicadoRepository, eventRepo eventDomain.EventoRepository) *CommunicationService {
	return &CommunicationService{
		comunicadoRepo: comunicadoRepo,
		eventRepo:      eventRepo,
	}
}

func (s *CommunicationService) CriarComunicado(ctx context.Context, userID, idEvento uuid.UUID, titulo, mensagem string) (*domain.Comunicado, error) {
	if err := s.verificarPermissaoEvento(ctx, userID, idEvento); err != nil {
		return nil, err
	}

	novoComunicado, err := domain.NewComunicado(idEvento, titulo, mensagem)
	if err != nil {
		return nil, err
	}

	if err := s.comunicadoRepo.Criar(ctx, novoComunicado); err != nil {
		return nil, fmt.Errorf("falha ao salvar comunicado: %w", err)
	}

	return novoComunicado, nil
}

func (s *CommunicationService) EditarComunicado(ctx context.Context, userID, idComunicado uuid.UUID, novoTitulo, novaMensagem string) error {
	comunicado, err := s.comunicadoRepo.BuscarPorID(ctx, idComunicado)
	if err != nil {
		return err
	}

	if err := s.verificarPermissaoEvento(ctx, userID, comunicado.IDEvento()); err != nil {
		return err
	}

	if err := comunicado.Editar(novoTitulo, novaMensagem); err != nil {
		return err
	}

	if err := s.comunicadoRepo.Editar(ctx, comunicado); err != nil {
		return fmt.Errorf("falha ao atualizar comunicado: %w", err)
	}

	return nil
}

func (s *CommunicationService) DeletarComunicado(ctx context.Context, userID, idComunicado uuid.UUID) error {
	comunicado, err := s.comunicadoRepo.BuscarPorID(ctx, idComunicado)
	if err != nil {
		return err
	}

	if err := s.verificarPermissaoEvento(ctx, userID, comunicado.IDEvento()); err != nil {
		return err
	}

	if err := s.comunicadoRepo.Deletar(ctx, idComunicado); err != nil {
		return fmt.Errorf("falha ao deletar comunicado: %w", err)
	}

	return nil
}

func (s *CommunicationService) ListarComunicadosPorEvento(ctx context.Context, idEvento uuid.UUID) ([]*domain.Comunicado, error) {
	comunicados, err := s.comunicadoRepo.BuscarPorEvento(ctx, idEvento)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar comunicados: %w", err)
	}

	return comunicados, nil
}

func (s *CommunicationService) verificarPermissaoEvento(ctx context.Context, userID, idEvento uuid.UUID) error {
	_, err := s.eventRepo.FindByID(ctx, userID, idEvento)
	if err != nil {
		return fmt.Errorf("evento não encontrado ou usuário não tem permissão: %w", err)
	}

	return nil
}
