// file: internal/billing/domain/assinatura.go
package domain

import (
	"errors"
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

var (
	ErrAssinaturaNaoEncontrada = errors.New("assinatura não encontrada")
	ErrAssinaturaJaExiste      = errors.New("assinatura já existe")
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

// UpdateStatus atualiza o status da assinatura.
func (a *Assinatura) Ativar(dataInicio, dataFim time.Time) {
	a.status = StatusAtiva
	a.dataInicio = dataInicio
	a.dataFim = dataFim
}
func HydrateAssinatura(id, idUsuario, idPlano uuid.UUID, dataInicio, dataFim time.Time, status StatusAssinatura) *Assinatura {
	return &Assinatura{
		id:         id,
		idUsuario:  idUsuario,
		idPlano:    idPlano,
		dataInicio: dataInicio,
		dataFim:    dataFim,
		status:     status,
	}
}

// Getters...
func (a *Assinatura) ID() uuid.UUID            { return a.id }
func (a *Assinatura) IDUsuario() uuid.UUID     { return a.idUsuario }
func (a *Assinatura) IDPlano() uuid.UUID       { return a.idPlano }
func (a *Assinatura) Status() StatusAssinatura { return a.status }
func (a *Assinatura) DataInicio() time.Time    { return a.dataInicio }
func (a *Assinatura) DataFim() time.Time       { return a.dataFim }
