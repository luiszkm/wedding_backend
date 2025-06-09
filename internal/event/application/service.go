// file: internal/event/application/service.go
package application

import (
	"context"
	"errors"
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

func (s *EventService) CriarNovoEvento(ctx context.Context, userID uuid.UUID, nome string, data time.Time, tipo string, urlSlug string) (*domain.Evento, error) {
	// Verificação de negócio: o slug da URL deve ser único
	_, err := s.repo.FindBySlug(ctx, urlSlug)
	if err == nil {
		return nil, domain.ErrSlugEmUso
	}
	if !errors.Is(err, domain.ErrEventoNaoEncontrado) {
		return nil, fmt.Errorf("erro ao verificar slug: %w", err)
	}

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
