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
	idStripePrice       string // <-- NOVO CAMPO

}

// Getters
func (p *Plano) ID() uuid.UUID            { return p.id }
func (p *Plano) Nome() string             { return p.nome }
func (p *Plano) PrecoEmCentavos() int     { return p.precoEmCentavos }
func (p *Plano) NumeroMaximoEventos() int { return p.numeroMaximoEventos }
func (p *Plano) DuracaoEmDias() int       { return p.duracaoEmDias }

// IDStripe retorna o ID do plano na Stripe, que é o mesmo que o ID do plano no sistema
func (p *Plano) IDStripePrice() string {
	return p.idStripePrice
}

// HydratePlano para reconstruir o objeto a partir do banco
func HydratePlano(id uuid.UUID, nome string, idStripePrice string, preco, eventos, dias int) *Plano {
	return &Plano{
		id:                  id,
		nome:                nome,
		precoEmCentavos:     preco,
		numeroMaximoEventos: eventos,
		duracaoEmDias:       dias,
		idStripePrice:       idStripePrice,
	}
}
