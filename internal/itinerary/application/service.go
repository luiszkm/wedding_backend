package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/itinerary/domain"
)

type ItineraryService struct {
	repo domain.ItineraryRepository
}

func NewItineraryService(repo domain.ItineraryRepository) *ItineraryService {
	return &ItineraryService{repo: repo}
}

// CreateItineraryItem cria um novo item do roteiro
func (s *ItineraryService) CreateItineraryItem(ctx context.Context, userID, eventID uuid.UUID, horario time.Time, titulo string, descricao *string) (uuid.UUID, error) {
	// 1. Usa a factory do domínio para criar o item
	item, err := domain.NewItineraryItem(eventID, horario, titulo, descricao)
	if err != nil {
		return uuid.Nil, fmt.Errorf("falha ao criar item do roteiro: %w", err)
	}

	// 2. Persiste o item (com verificação de propriedade do evento)
	if err := s.repo.Save(ctx, userID, item); err != nil {
		return uuid.Nil, fmt.Errorf("falha ao salvar item do roteiro: %w", err)
	}

	return item.ID(), nil
}

// GetItineraryByEventID retorna todos os itens do roteiro de um evento (acesso público)
func (s *ItineraryService) GetItineraryByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.ItineraryItem, error) {
	items, err := s.repo.FindByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter itens do roteiro: %w", err)
	}
	return items, nil
}

// GetItineraryByEventIDAndUserID retorna itens do roteiro para o proprietário do evento
func (s *ItineraryService) GetItineraryByEventIDAndUserID(ctx context.Context, eventID, userID uuid.UUID) ([]*domain.ItineraryItem, error) {
	items, err := s.repo.FindByEventIDAndUserID(ctx, eventID, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter itens do roteiro: %w", err)
	}
	return items, nil
}

// UpdateItineraryItem atualiza um item do roteiro
func (s *ItineraryService) UpdateItineraryItem(ctx context.Context, userID, itemID uuid.UUID, horario time.Time, titulo string, descricao *string) error {
	// 1. Carrega o item existente (com verificação de propriedade)
	item, err := s.repo.FindByID(ctx, userID, itemID)
	if err != nil {
		return fmt.Errorf("falha ao buscar item do roteiro: %w", err)
	}

	// 2. Aplica as mudanças usando a lógica de domínio
	if err := item.Update(horario, titulo, descricao); err != nil {
		return err
	}

	// 3. Persiste as alterações
	if err := s.repo.Update(ctx, userID, item); err != nil {
		return fmt.Errorf("falha ao atualizar item do roteiro: %w", err)
	}

	return nil
}

// DeleteItineraryItem remove um item do roteiro
func (s *ItineraryService) DeleteItineraryItem(ctx context.Context, userID, itemID uuid.UUID) error {
	if err := s.repo.Delete(ctx, userID, itemID); err != nil {
		return fmt.Errorf("falha ao deletar item do roteiro: %w", err)
	}
	return nil
}

// GetItineraryItemByID retorna um item específico do roteiro
func (s *ItineraryService) GetItineraryItemByID(ctx context.Context, userID, itemID uuid.UUID) (*domain.ItineraryItem, error) {
	item, err := s.repo.FindByID(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter item do roteiro: %w", err)
	}
	return item, nil
}