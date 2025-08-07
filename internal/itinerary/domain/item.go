package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTituloAtividadeObrigatorio  = errors.New("título da atividade é obrigatório")
	ErrTituloAtividadeMuitoLongo   = errors.New("título da atividade não pode exceder 255 caracteres")
	ErrHorarioObrigatorio         = errors.New("horário da atividade é obrigatório")
	ErrEventoObrigatorio          = errors.New("ID do evento é obrigatório")
	ErrItemRoteiroNaoEncontrado   = errors.New("item do roteiro não encontrado")
)

// ItineraryItem representa um item individual no roteiro do evento
type ItineraryItem struct {
	id                uuid.UUID
	idEvento          uuid.UUID
	horario           time.Time
	tituloAtividade   string
	descricaoAtividade *string
	createdAt         time.Time
	updatedAt         time.Time
}

// NewItineraryItem é a função factory para criar um novo item do roteiro
func NewItineraryItem(idEvento uuid.UUID, horario time.Time, tituloAtividade string, descricaoAtividade *string) (*ItineraryItem, error) {
	if idEvento == uuid.Nil {
		return nil, ErrEventoObrigatorio
	}

	if horario.IsZero() {
		return nil, ErrHorarioObrigatorio
	}

	titulo := strings.TrimSpace(tituloAtividade)
	if titulo == "" {
		return nil, ErrTituloAtividadeObrigatorio
	}
	if len(titulo) > 255 {
		return nil, ErrTituloAtividadeMuitoLongo
	}

	now := time.Now()
	return &ItineraryItem{
		id:                uuid.New(),
		idEvento:          idEvento,
		horario:           horario,
		tituloAtividade:   titulo,
		descricaoAtividade: descricaoAtividade,
		createdAt:         now,
		updatedAt:         now,
	}, nil
}

// HydrateItineraryItem reconstrói um item do roteiro a partir de dados persistidos
func HydrateItineraryItem(id, idEvento uuid.UUID, horario time.Time, tituloAtividade string, descricaoAtividade *string, createdAt, updatedAt time.Time) *ItineraryItem {
	return &ItineraryItem{
		id:                id,
		idEvento:          idEvento,
		horario:           horario,
		tituloAtividade:   tituloAtividade,
		descricaoAtividade: descricaoAtividade,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// Update atualiza os dados do item do roteiro
func (i *ItineraryItem) Update(horario time.Time, tituloAtividade string, descricaoAtividade *string) error {
	if horario.IsZero() {
		return ErrHorarioObrigatorio
	}

	titulo := strings.TrimSpace(tituloAtividade)
	if titulo == "" {
		return ErrTituloAtividadeObrigatorio
	}
	if len(titulo) > 255 {
		return ErrTituloAtividadeMuitoLongo
	}

	i.horario = horario
	i.tituloAtividade = titulo
	i.descricaoAtividade = descricaoAtividade
	i.updatedAt = time.Now()

	return nil
}

// Getters para expor campos privados de forma controlada
func (i *ItineraryItem) ID() uuid.UUID               { return i.id }
func (i *ItineraryItem) IDEvento() uuid.UUID         { return i.idEvento }
func (i *ItineraryItem) Horario() time.Time          { return i.horario }
func (i *ItineraryItem) TituloAtividade() string     { return i.tituloAtividade }
func (i *ItineraryItem) DescricaoAtividade() *string { return i.descricaoAtividade }
func (i *ItineraryItem) CreatedAt() time.Time        { return i.createdAt }
func (i *ItineraryItem) UpdatedAt() time.Time        { return i.updatedAt }