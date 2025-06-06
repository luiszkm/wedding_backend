// file: internal/guest/domain/group.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrChaveDeAcessoObrigatoria = errors.New("a chave de acesso é obrigatória")
	ErrPeloMenosUmConvidado     = errors.New("o grupo deve ter pelo menos um convidado")
	ErrGrupoNaoEncontrado       = errors.New("grupo de convidados não encontrado")
)

// GrupoDeConvidados é o agregado raiz para o contexto de RSVP.
// Ele é a fronteira de consistência para as operações.
type GrupoDeConvidados struct {
	id            uuid.UUID
	idCasamento   uuid.UUID
	chaveDeAcesso string
	convidados    []*Convidado
	createdAt     time.Time
	updatedAt     time.Time
}

// Convidado é uma entidade interna do agregado.
type Convidado struct {
	id         uuid.UUID
	nome       string
	statusRSVP string
}

const (
	StatusRSVPConfirmado = "CONFIRMADO"
	StatusRSVPRecusado   = "RECUSADO"
)

// NewGrupoDeConvidados é a fábrica para nosso agregado. Garante que ele seja criado em um estado válido.
func NewGrupoDeConvidados(idCasamento uuid.UUID, chaveDeAcesso string, nomesDosConvidados []string) (*GrupoDeConvidados, error) {
	if chaveDeAcesso == "" {
		return nil, ErrChaveDeAcessoObrigatoria
	}
	if len(nomesDosConvidados) == 0 {
		return nil, ErrPeloMenosUmConvidado
	}

	convidados := make([]*Convidado, len(nomesDosConvidados))
	for i, nome := range nomesDosConvidados {
		convidados[i] = &Convidado{
			id:   uuid.New(),
			nome: nome,
		}
	}

	return &GrupoDeConvidados{
		id:            uuid.New(),
		idCasamento:   idCasamento,
		chaveDeAcesso: chaveDeAcesso,
		convidados:    convidados,
	}, nil
}

func HydrateGroup(id, idCasamento uuid.UUID, chaveDeAcesso string, convidados []*Convidado) *GrupoDeConvidados {
	return &GrupoDeConvidados{
		id:            id,
		idCasamento:   idCasamento,
		chaveDeAcesso: chaveDeAcesso,
		convidados:    convidados,
	}
}

func HydrateConvidado(id uuid.UUID, nome, statusRSVP string) *Convidado {
	return &Convidado{
		id:         id,
		nome:       nome,
		statusRSVP: statusRSVP,
	}
}

// Getters para expor campos privados de forma controlada
func (g *GrupoDeConvidados) ID() uuid.UUID            { return g.id }
func (g *GrupoDeConvidados) IDCasamento() uuid.UUID   { return g.idCasamento }
func (g *GrupoDeConvidados) ChaveDeAcesso() string    { return g.chaveDeAcesso }
func (g *GrupoDeConvidados) Convidados() []*Convidado { return g.convidados }
func (c *Convidado) ID() uuid.UUID                    { return c.id }
func (c *Convidado) Nome() string                     { return c.nome }
func (c *Convidado) StatusRSVP() string               { return c.statusRSVP }
