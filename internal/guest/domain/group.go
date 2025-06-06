// file: internal/guest/domain/group.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrChaveDeAcessoObrigatoria      = errors.New("a chave de acesso é obrigatória")
	ErrPeloMenosUmConvidado          = errors.New("o grupo deve ter pelo menos um convidado")
	ErrGrupoNaoEncontrado            = errors.New("grupo de convidados não encontrado")
	ErrConvidadoNaoEncontradoNoGrupo = errors.New("um ou mais convidados não pertencem a este grupo")
	ErrStatusRSVPInvalido            = errors.New("status de rsvp inválido")
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
type RespostaRSVP struct {
	ConvidadoID uuid.UUID
	Status      string
}

// Convidado é uma entidade interna do agregado.
type Convidado struct {
	id         uuid.UUID
	nome       string
	statusRSVP string
}

type ConvidadoParaRevisao struct {
	ID   uuid.UUID // Pode ser uuid.Nil se for um novo convidado
	Nome string
}

const (
	StatusRSVPConfirmado = "CONFIRMADO"
	StatusRSVPRecusado   = "RECUSADO"
	StatusRSVPPendente   = "PENDENTE"
)

func (g *GrupoDeConvidados) UpdatedAt() time.Time {
	return g.updatedAt
}

// Revisar atualiza o estado do agregado com base em novos dados.
func (g *GrupoDeConvidados) Revisar(novaChaveDeAcesso string, convidadosParaRevisao []ConvidadoParaRevisao) error {
	if novaChaveDeAcesso == "" {
		return ErrChaveDeAcessoObrigatoria
	}
	if len(convidadosParaRevisao) == 0 {
		return ErrPeloMenosUmConvidado
	}

	g.chaveDeAcesso = novaChaveDeAcesso

	convidadosAtuaisMap := make(map[uuid.UUID]*Convidado)
	for _, c := range g.convidados {
		convidadosAtuaisMap[c.id] = c
	}

	var convidadosFinais []*Convidado
	idsProcessados := make(map[uuid.UUID]bool)

	for _, cRevisao := range convidadosParaRevisao {
		// Se o ID for zero, é um novo convidado
		if cRevisao.ID == uuid.Nil {
			novoConvidado := &Convidado{
				id:         uuid.New(),
				nome:       cRevisao.Nome,
				statusRSVP: StatusRSVPPendente,
			}
			convidadosFinais = append(convidadosFinais, novoConvidado)
		} else {
			// Se o ID existe, é um convidado existente (potencialmente renomeado)
			convidadoExistente, ok := convidadosAtuaisMap[cRevisao.ID]
			if !ok {
				// Tentou editar um convidado que não pertencia a este grupo
				return ErrConvidadoNaoEncontradoNoGrupo
			}
			convidadoExistente.nome = cRevisao.Nome // Atualiza o nome
			convidadosFinais = append(convidadosFinais, convidadoExistente)
			idsProcessados[cRevisao.ID] = true
		}
	}

	// Substitui a lista de convidados antiga pela nova
	g.convidados = convidadosFinais
	g.updatedAt = time.Now()

	return nil
}

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

func (g *GrupoDeConvidados) ConfirmarPresenca(respostas []RespostaRSVP) error {
	// Cria um mapa para busca rápida dos convidados do grupo.
	convidadosDoGrupo := make(map[uuid.UUID]*Convidado)
	for _, c := range g.convidados {
		convidadosDoGrupo[c.id] = c
	}

	// Primeira passagem: validação. Garante que a operação seja atômica.
	for _, resposta := range respostas {
		// Regra 1: O status deve ser válido.
		if resposta.Status != StatusRSVPConfirmado && resposta.Status != StatusRSVPRecusado {
			return ErrStatusRSVPInvalido
		}
		// Regra 2: O convidado da resposta deve pertencer ao grupo.
		if _, ok := convidadosDoGrupo[resposta.ConvidadoID]; !ok {
			return ErrConvidadoNaoEncontradoNoGrupo
		}
	}

	// Segunda passagem: atualização. Ocorre apenas se toda a validação passou.
	for _, resposta := range respostas {
		convidado := convidadosDoGrupo[resposta.ConvidadoID]
		convidado.statusRSVP = resposta.Status
	}

	g.updatedAt = time.Now()
	return nil
}

// Getters para expor campos privados de forma controlada
func (g *GrupoDeConvidados) ID() uuid.UUID            { return g.id }
func (g *GrupoDeConvidados) IDCasamento() uuid.UUID   { return g.idCasamento }
func (g *GrupoDeConvidados) ChaveDeAcesso() string    { return g.chaveDeAcesso }
func (g *GrupoDeConvidados) Convidados() []*Convidado { return g.convidados }
func (c *Convidado) ID() uuid.UUID                    { return c.id }
func (c *Convidado) Nome() string                     { return c.nome }
func (c *Convidado) StatusRSVP() string               { return c.statusRSVP }
