// file: internal/billing/domain/assinatura.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type StatusAssinatura string

const (
	StatusPendente  StatusAssinatura = "PENDENTE"
	StatusAtiva     StatusAssinatura = "ATIVA"
	StatusExpirada  StatusAssinatura = "EXPIRADA"
	StatusCancelada StatusAssinatura = "CANCELADA"
)

type Assinatura struct {
	id         uuid.UUID
	idUsuario  uuid.UUID
	idPlano    uuid.UUID
	dataInicio time.Time
	dataFim    time.Time
	status     StatusAssinatura
}

// NewAssinatura cria uma nova assinatura no estado PENDENTE.
func NewAssinatura(idUsuario, idPlano uuid.UUID) *Assinatura {
	return &Assinatura{
		id:        uuid.New(),
		idUsuario: idUsuario,
		idPlano:   idPlano,
		status:    StatusPendente, // Toda nova assinatura começa como pendente
	}
}

// Getters...
func (a *Assinatura) ID() uuid.UUID            { return a.id }
func (a *Assinatura) IDUsuario() uuid.UUID     { return a.idUsuario }
func (a *Assinatura) IDPlano() uuid.UUID       { return a.idPlano }
func (a *Assinatura) Status() StatusAssinatura { return a.status }
