// file: internal/event/domain/evento.go
package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TipoEvento string

const (
	TipoCasamento   TipoEvento = "CASAMENTO"
	TipoAniversario TipoEvento = "ANIVERSARIO"
	TipoChaDeBebe   TipoEvento = "CHA_DE_BEBE"
	TipoOutro       TipoEvento = "OUTRO"
)

var (
	ErrEventoNaoEncontrado = errors.New("evento não encontrado")
	ErrTipoEventoInvalido  = errors.New("tipo de evento inválido")
	ErrSlugEmUso           = errors.New("a URL amigável (slug) já está em uso")
	ErrEventoJaExiste      = errors.New("evento já existe")
)

func (t TipoEvento) IsValid() bool {
	switch t {
	case TipoCasamento, TipoAniversario, TipoChaDeBebe, TipoOutro:
		return true
	}
	return false
}

type Evento struct {
	id        uuid.UUID
	idUsuario uuid.UUID
	nome      string
	data      time.Time
	tipo      TipoEvento
	urlSlug   string
}

func NewEvento(idUsuario uuid.UUID, nome string, data time.Time, tipo TipoEvento, urlSlug string) (*Evento, error) {
	if nome == "" || urlSlug == "" {
		return nil, errors.New("nome e urlSlug são obrigatórios")
	}
	if !tipo.IsValid() {
		return nil, errors.New("tipo de evento inválido")
	}
	return &Evento{
		id:        uuid.New(),
		idUsuario: idUsuario,
		nome:      strings.TrimSpace(nome),
		data:      data,
		tipo:      tipo,
		urlSlug:   strings.TrimSpace(urlSlug),
	}, nil
}

// HydrateEvento cria uma nova instância de Evento a partir dos dados fornecidos.
func HydrateEvento(id, idUsuario uuid.UUID, nome string, data time.Time, tipo TipoEvento, urlSlug string) *Evento {
	if nome == "" || urlSlug == "" {
		return nil // ou retornar um erro, dependendo da lógica de negócio
	}
	if !tipo.IsValid() {
		return nil // ou retornar um erro, dependendo da lógica de negócio
	}
	return &Evento{
		id:        id,
		idUsuario: idUsuario,
		nome:      strings.TrimSpace(nome),
		data:      data,
		tipo:      tipo,
		urlSlug:   strings.TrimSpace(urlSlug),
	}
}

// Getters
func (e *Evento) ID() uuid.UUID        { return e.id }
func (e *Evento) IDUsuario() uuid.UUID { return e.idUsuario }
func (e *Evento) Nome() string         { return e.nome }
func (e *Evento) Data() time.Time      { return e.data }
func (e *Evento) Tipo() TipoEvento     { return e.tipo }
func (e *Evento) UrlSlug() string      { return e.urlSlug }
