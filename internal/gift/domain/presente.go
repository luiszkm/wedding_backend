// file: internal/gift/domain/presente.go
package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	TipoDetalheProdutoExterno = "PRODUTO_EXTERNO"
)

var (
	ErrNomePresenteObrigatorio = errors.New("o nome do presente é obrigatório")
	ErrDetalhesInvalidos       = errors.New("os detalhes do presente são inválidos")
)

// DetalhesPresente representa a parte polimórfica de um presente.
type DetalhesPresente struct {
	Tipo       string
	LinkDaLoja string
	// outros campos como ChavePix viriam aqui
}

// Presente é o agregado raiz do contexto.
type Presente struct {
	id          uuid.UUID
	idCasamento uuid.UUID
	nome        string
	descricao   string
	fotoURL     string
	ehFavorito  bool
	status      string // e.g., "DISPONIVEL"
	categoria   string
	detalhes    DetalhesPresente
}
type ErrPresentesConflitantes struct {
	PresentesIDs []uuid.UUID
}

// NewPresente é a fábrica para criar um presente em estado válido.
func NewPresente(idCasamento uuid.UUID, nome, descricao,
	fotoURL string, ehFavorito bool, categoria string,
	detalhes DetalhesPresente) (*Presente, error) {
	if nome == "" {
		return nil, ErrNomePresenteObrigatorio
	}
	if detalhes.Tipo == TipoDetalheProdutoExterno && detalhes.LinkDaLoja == "" {
		return nil, ErrDetalhesInvalidos
	}

	return &Presente{
		id:          uuid.New(),
		idCasamento: idCasamento,
		nome:        nome,
		descricao:   descricao,
		fotoURL:     fotoURL,
		ehFavorito:  ehFavorito,
		status:      "DISPONIVEL",
		categoria:   categoria, // Exemplo de categoria, poderia ser um campo adicional
		detalhes:    detalhes,
	}, nil
}

func HydratePresente(id, idCasamento uuid.UUID, nome, descricao, fotoURL, status, categoria string, ehFavorito bool, detalhes DetalhesPresente) *Presente {
	return &Presente{
		id:          id,
		idCasamento: idCasamento,
		nome:        nome,
		descricao:   descricao,
		fotoURL:     fotoURL,
		status:      status,
		categoria:   categoria,
		ehFavorito:  ehFavorito,
		detalhes:    detalhes,
	}
}
func (e *ErrPresentesConflitantes) Error() string {
	return fmt.Sprintf("um ou mais presentes já foram selecionados: %v", e.PresentesIDs)
}

// Getters para acesso controlado aos campos.
func (p *Presente) ID() uuid.UUID              { return p.id }
func (p *Presente) IDCasamento() uuid.UUID     { return p.idCasamento }
func (p *Presente) Nome() string               { return p.nome }
func (p *Presente) Descricao() string          { return p.descricao }
func (p *Presente) FotoURL() string            { return p.fotoURL }
func (p *Presente) EhFavorito() bool           { return p.ehFavorito }
func (p *Presente) Status() string             { return p.status }
func (p *Presente) Detalhes() DetalhesPresente { return p.detalhes }
func (p *Presente) Categoria() string          { return p.categoria }
