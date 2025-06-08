// file: internal/messageboard/domain/recado.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNomeAutorObrigatorio        = errors.New("o nome do autor é obrigatório")
	ErrTextoObrigatorio            = errors.New("o texto do recado é obrigatório")
	ErrStatusInvalidoParaModeracao = errors.New("status de moderação inválido")
	ErrRecadoNaoEncontrado         = errors.New("recado não encontrado")
)

const (
	StatusPendente  = "PENDENTE"
	StatusAprovado  = "APROVADO"
	StatusRejeitado = "REJEITADO"
)

// Recado é o agregado raiz do contexto de Mural de Recados.
type Recado struct {
	id                  uuid.UUID
	idCasamento         uuid.UUID
	idGrupoDeConvidados uuid.UUID
	nomeDoAutor         string
	texto               string
	status              string
	ehFavorito          bool
	dataDeCriacao       time.Time
}

// NewRecado é a fábrica para criar um recado em estado válido.
func NewRecado(idCasamento, idGrupo uuid.UUID, nomeAutor, texto string) (*Recado, error) {
	if nomeAutor == "" {
		return nil, ErrNomeAutorObrigatorio
	}
	if texto == "" {
		return nil, ErrTextoObrigatorio
	}

	return &Recado{
		id:                  uuid.New(),
		idCasamento:         idCasamento,
		idGrupoDeConvidados: idGrupo,
		nomeDoAutor:         nomeAutor,
		texto:               texto,
		status:              StatusPendente, // Todo novo recado começa como pendente.
		ehFavorito:          false,
		dataDeCriacao:       time.Now(),
	}, nil
}
func HydrateRecado(id, idCasamento, idGrupo uuid.UUID, nomeAutor, texto, status string, ehFavorito bool, dataCriacao time.Time) *Recado {
	return &Recado{
		id:                  id,
		idCasamento:         idCasamento,
		idGrupoDeConvidados: idGrupo,
		nomeDoAutor:         nomeAutor,
		texto:               texto,
		status:              status,
		ehFavorito:          ehFavorito,
		dataDeCriacao:       dataCriacao,
	}
}
func (r *Recado) Aprovar() {
	r.status = StatusAprovado
}
func (r *Recado) Rejeitar() {
	r.status = StatusRejeitado
}
func (r *Recado) DefinirFavorito(favorito bool) {
	r.ehFavorito = favorito
}

// Getters para acesso controlado.
func (r *Recado) ID() uuid.UUID                  { return r.id }
func (r *Recado) IDCasamento() uuid.UUID         { return r.idCasamento }
func (r *Recado) IDGrupoDeConvidados() uuid.UUID { return r.idGrupoDeConvidados }
func (r *Recado) NomeDoAutor() string            { return r.nomeDoAutor }
func (r *Recado) Texto() string                  { return r.texto }
func (r *Recado) Status() string                 { return r.status }
func (r *Recado) EhFavorito() bool               { return r.ehFavorito }
func (r *Recado) DataDeCriacao() time.Time       { return r.dataDeCriacao }
