// file: internal/gift/domain/selecao.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Selecao representa o agregado de uma seleção de presentes feita por um grupo.
type Selecao struct {
	id                   uuid.UUID
	idCasamento          uuid.UUID
	idGrupoDeConvidados  uuid.UUID
	dataDaSelecao        time.Time
	presentesConfirmados []PresenteConfirmado
}

// PresenteConfirmado é uma entidade de valor dentro da seleção.
type PresenteConfirmado struct {
	ID   uuid.UUID
	Nome string
}

// NewSelecao é a fábrica para o agregado.
func NewSelecao(idCasamento, idGrupo uuid.UUID, presentes []PresenteConfirmado) *Selecao {
	return &Selecao{
		id:                   uuid.New(),
		idCasamento:          idCasamento,
		idGrupoDeConvidados:  idGrupo,
		dataDaSelecao:        time.Now(),
		presentesConfirmados: presentes,
	}
}
func HydrateSelecao(id, idCasamento, idGrupo uuid.UUID, presentes []PresenteConfirmado, dataDaSelecao time.Time) *Selecao {
	return &Selecao{
		id:                   id,
		idCasamento:          idCasamento,
		idGrupoDeConvidados:  idGrupo,
		presentesConfirmados: presentes,
		dataDaSelecao:        dataDaSelecao,
	}
}

// Getters
func (s *Selecao) ID() uuid.UUID                              { return s.id }
func (s *Selecao) PresentesConfirmados() []PresenteConfirmado { return s.presentesConfirmados }
func (s *Selecao) DataDaSelecao() time.Time                   { return s.dataDaSelecao }
