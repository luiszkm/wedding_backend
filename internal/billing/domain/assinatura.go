// file: internal/billing/domain/assinatura.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrAssinaturaNaoEncontrada = errors.New("assinatura não encontrada")

type StatusAssinatura string

const (
	StatusPendente     StatusAssinatura = "PENDENTE"
	StatusAtiva        StatusAssinatura = "ATIVA"
	StatusExpirada     StatusAssinatura = "EXPIRADA"
	StatusCancelada    StatusAssinatura = "CANCELADA"
	StatusInadimplente StatusAssinatura = "INADIMPLENTE"
)

type Assinatura struct {
	id                   uuid.UUID
	idUsuario            uuid.UUID
	idPlano              uuid.UUID
	idStripeSubscription string
	dataInicio           time.Time
	dataFim              time.Time
	status               StatusAssinatura
}

func NewAssinatura(idUsuario, idPlano uuid.UUID) *Assinatura {
	return &Assinatura{
		id:        uuid.New(),
		idUsuario: idUsuario,
		idPlano:   idPlano,
		status:    StatusPendente,
	}
}
func (a *Assinatura) Renovar(novaDataFim time.Time) {
	a.status = StatusAtiva // Garante que o status volte para ativo caso estivesse inadimplente
	a.dataFim = novaDataFim
}

func HydrateAssinatura(id, idUsuario, idPlano uuid.UUID, idStripeSub string, dataInicio, dataFim time.Time, status StatusAssinatura) *Assinatura {
	return &Assinatura{
		id:                   id,
		idUsuario:            idUsuario,
		idPlano:              idPlano,
		idStripeSubscription: idStripeSub,
		dataInicio:           dataInicio,
		dataFim:              dataFim,
		status:               status,
	}
}

func (a *Assinatura) Ativar(idStripeSub string, dataInicio, dataFim time.Time) {
	a.idStripeSubscription = idStripeSub
	a.status = StatusAtiva
	a.dataInicio = dataInicio
	a.dataFim = dataFim
}

func (a *Assinatura) Cancelar() {
	a.status = StatusCancelada
}

// Getters...
func (a *Assinatura) ID() uuid.UUID                { return a.id }
func (a *Assinatura) IDUsuario() uuid.UUID         { return a.idUsuario }
func (a *Assinatura) IDPlano() uuid.UUID           { return a.idPlano }
func (a *Assinatura) IDStripeSubscription() string { return a.idStripeSubscription }
func (a *Assinatura) Status() StatusAssinatura     { return a.status }
func (a *Assinatura) DataInicio() time.Time        { return a.dataInicio }
func (a *Assinatura) DataFim() time.Time           { return a.dataFim }
