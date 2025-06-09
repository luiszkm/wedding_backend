// file: internal/billing/domain/plano.go
package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrPlanoNaoEncontrado = errors.New("plano não encontrado")

type Plano struct {
	id                  uuid.UUID
	nome                string
	precoEmCentavos     int
	numeroMaximoEventos int
	duracaoEmDias       int
}

// Getters
func (p *Plano) ID() uuid.UUID            { return p.id }
func (p *Plano) Nome() string             { return p.nome }
func (p *Plano) PrecoEmCentavos() int     { return p.precoEmCentavos }
func (p *Plano) NumeroMaximoEventos() int { return p.numeroMaximoEventos }
func (p *Plano) DuracaoEmDias() int       { return p.duracaoEmDias }

// HydratePlano para reconstruir o objeto a partir do banco
func HydratePlano(id uuid.UUID, nome string, preco, eventos, dias int) *Plano {
	return &Plano{
		id:                  id,
		nome:                nome,
		precoEmCentavos:     preco,
		numeroMaximoEventos: eventos,
		duracaoEmDias:       dias,
	}
}
