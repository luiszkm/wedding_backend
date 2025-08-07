package domain

import (
	"context"

	"github.com/google/uuid"
)

type ItineraryRepository interface {
	// Save persiste um novo item do roteiro
	Save(ctx context.Context, userID uuid.UUID, item *ItineraryItem) error

	// FindByEventID retorna todos os itens do roteiro de um evento, ordenados por horário
	FindByEventID(ctx context.Context, eventID uuid.UUID) ([]*ItineraryItem, error)

	// FindByID busca um item específico do roteiro
	FindByID(ctx context.Context, userID, itemID uuid.UUID) (*ItineraryItem, error)

	// Update atualiza um item existente do roteiro
	Update(ctx context.Context, userID uuid.UUID, item *ItineraryItem) error

	// Delete remove um item do roteiro
	Delete(ctx context.Context, userID, itemID uuid.UUID) error

	// FindByEventIDAndUserID retorna itens do roteiro de um evento específico do usuário (para operações autenticadas)
	FindByEventIDAndUserID(ctx context.Context, eventID, userID uuid.UUID) ([]*ItineraryItem, error)
}