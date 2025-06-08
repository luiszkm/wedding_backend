// file: internal/event/application/service.go
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/event/domain"
)

type EventService struct {
	repo domain.EventoRepository
}

func NewEventService(repo domain.EventoRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CriarNovoEvento(ctx context.Context, userID uuid.UUID, nome string, data time.Time, tipo, urlSlug string) (*domain.Evento, error) {
	// Aqui poderíamos adicionar uma lógica para verificar se o slug já existe,
	// usando o repo.FindBySlug que adicionamos na interface.

	tipoEvento := domain.TipoEvento(tipo)

	novoEvento, err := domain.NewEvento(userID, nome, data, tipoEvento, urlSlug)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, novoEvento); err != nil {
		return nil, fmt.Errorf("falha ao salvar novo evento: %w", err)
	}

	return novoEvento, nil
}
