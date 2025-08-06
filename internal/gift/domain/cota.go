package domain

import (
	"errors"

	"github.com/google/uuid"
)

const (
	StatusCotaDisponivel  = "DISPONIVEL"
	StatusCotaSelecionada = "SELECIONADO"
)

var (
	ErrCotaJaSelecionada  = errors.New("cota já foi selecionada")
	ErrCotaJaDisponivel   = errors.New("cota já está disponível")
	ErrValorCotaInvalido  = errors.New("valor da cota deve ser positivo")
	ErrNumeroCotaInvalido = errors.New("número da cota deve ser positivo")
)

type Cota struct {
	id         uuid.UUID
	idPresente uuid.UUID
	numeroCota int
	valorCota  float64
	status     string
	idSelecao  *uuid.UUID
}

func NewCota(idPresente uuid.UUID, numeroCota int, valorCota float64) (*Cota, error) {
	if numeroCota <= 0 {
		return nil, ErrNumeroCotaInvalido
	}
	if valorCota <= 0 {
		return nil, ErrValorCotaInvalido
	}

	return &Cota{
		id:         uuid.New(),
		idPresente: idPresente,
		numeroCota: numeroCota,
		valorCota:  valorCota,
		status:     StatusCotaDisponivel,
		idSelecao:  nil,
	}, nil
}

func HydrateCota(id, idPresente uuid.UUID, numeroCota int, valorCota float64, status string, idSelecao *uuid.UUID) *Cota {
	return &Cota{
		id:         id,
		idPresente: idPresente,
		numeroCota: numeroCota,
		valorCota:  valorCota,
		status:     status,
		idSelecao:  idSelecao,
	}
}

func (c *Cota) Selecionar(idSelecao uuid.UUID) error {
	if c.status == StatusCotaSelecionada {
		return ErrCotaJaSelecionada
	}

	c.status = StatusCotaSelecionada
	c.idSelecao = &idSelecao
	return nil
}

func (c *Cota) LiberarSelecao() error {
	if c.status == StatusCotaDisponivel {
		return ErrCotaJaDisponivel
	}

	c.status = StatusCotaDisponivel
	c.idSelecao = nil
	return nil
}

func (c *Cota) EstaDisponivel() bool {
	return c.status == StatusCotaDisponivel
}

func (c *Cota) EstaSelecionada() bool {
	return c.status == StatusCotaSelecionada
}

func (c *Cota) ID() uuid.UUID         { return c.id }
func (c *Cota) IDPresente() uuid.UUID { return c.idPresente }
func (c *Cota) NumeroCota() int       { return c.numeroCota }
func (c *Cota) ValorCota() float64    { return c.valorCota }
func (c *Cota) Status() string        { return c.status }
func (c *Cota) IDSelecao() *uuid.UUID { return c.idSelecao }
