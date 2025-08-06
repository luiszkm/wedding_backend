package domain

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

const (
	TipoDetalheProdutoExterno     = "PRODUTO_EXTERNO"
	TipoPresenteIntegral          = "INTEGRAL"
	TipoPresenteFracionado        = "FRACIONADO"
	StatusDisponivel              = "DISPONIVEL"
	StatusSelecionado             = "SELECIONADO"
	StatusParcialmenteSelecionado = "PARCIALMENTE_SELECIONADO"
)

var (
	ErrNomePresenteObrigatorio = errors.New("o nome do presente é obrigatório")
	ErrDetalhesInvalidos       = errors.New("os detalhes do presente são inválidos")
	ErrTipoPresenteInvalido    = errors.New("tipo de presente inválido")
	ErrValorTotalInvalido      = errors.New("valor total deve ser positivo para presentes fracionados")
	ErrNumeroCotasInvalido     = errors.New("número de cotas deve ser maior que 1 para presentes fracionados")
	ErrPresenteNaoFracionado   = errors.New("operação válida apenas para presentes fracionados")
	ErrCotasIndisponiveis      = errors.New("não há cotas suficientes disponíveis")
	ErrPresenteJaSelecionado   = errors.New("presente já foi completamente selecionado")
)

type DetalhesPresente struct {
	Tipo       string
	LinkDaLoja string
}

type Presente struct {
	id          uuid.UUID
	idCasamento uuid.UUID
	nome        string
	descricao   string
	fotoURL     string
	ehFavorito  bool
	status      string
	categoria   string
	detalhes    DetalhesPresente
	tipo        string
	valorTotal  *float64
	cotas       []*Cota
}

type ErrPresentesConflitantes struct {
	PresentesIDs []uuid.UUID
}

func (e *ErrPresentesConflitantes) Error() string {
	return fmt.Sprintf("um ou mais presentes já foram selecionados: %v", e.PresentesIDs)
}

func NewPresenteIntegral(idCasamento uuid.UUID, nome, descricao, fotoURL string,
	ehFavorito bool, categoria string, detalhes DetalhesPresente) (*Presente, error) {

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
		status:      StatusDisponivel,
		categoria:   categoria,
		detalhes:    detalhes,
		tipo:        TipoPresenteIntegral,
		valorTotal:  nil,
		cotas:       nil,
	}, nil
}

func NewPresenteFracionado(idCasamento uuid.UUID, nome, descricao, fotoURL string,
	ehFavorito bool, categoria string, detalhes DetalhesPresente,
	valorTotal float64, numeroCotas int) (*Presente, error) {

	if nome == "" {
		return nil, ErrNomePresenteObrigatorio
	}
	if detalhes.Tipo == TipoDetalheProdutoExterno && detalhes.LinkDaLoja == "" {
		return nil, ErrDetalhesInvalidos
	}
	if valorTotal <= 0 {
		return nil, ErrValorTotalInvalido
	}
	if numeroCotas <= 1 {
		return nil, ErrNumeroCotasInvalido
	}

	presente := &Presente{
		id:          uuid.New(),
		idCasamento: idCasamento,
		nome:        nome,
		descricao:   descricao,
		fotoURL:     fotoURL,
		ehFavorito:  ehFavorito,
		status:      StatusDisponivel,
		categoria:   categoria,
		detalhes:    detalhes,
		tipo:        TipoPresenteFracionado,
		valorTotal:  &valorTotal,
		cotas:       make([]*Cota, 0, numeroCotas),
	}

	valorCota := math.Round((valorTotal/float64(numeroCotas))*100) / 100

	for i := 1; i <= numeroCotas; i++ {
		cota, err := NewCota(presente.id, i, valorCota)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar cota %d: %w", i, err)
		}
		presente.cotas = append(presente.cotas, cota)
	}

	return presente, nil
}

func HydratePresente(id, idCasamento uuid.UUID, nome, descricao, fotoURL, status, categoria, tipo string,
	ehFavorito bool, detalhes DetalhesPresente, valorTotal *float64, cotas []*Cota) *Presente {

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
		tipo:        tipo,
		valorTotal:  valorTotal,
		cotas:       cotas,
	}
}

func (p *Presente) SelecionarCotas(quantidade int, idSelecao uuid.UUID) error {
	if p.tipo != TipoPresenteFracionado {
		return ErrPresenteNaoFracionado
	}

	if p.status == StatusSelecionado {
		return ErrPresenteJaSelecionado
	}

	cotasDisponiveis := p.ContarCotasDisponiveis()
	if quantidade > cotasDisponiveis {
		return ErrCotasIndisponiveis
	}

	cotasSelecionadas := 0
	for _, cota := range p.cotas {
		if cotasSelecionadas >= quantidade {
			break
		}
		if cota.EstaDisponivel() {
			if err := cota.Selecionar(idSelecao); err != nil {
				return fmt.Errorf("erro ao selecionar cota %d: %w", cota.numeroCota, err)
			}
			cotasSelecionadas++
		}
	}

	p.atualizarStatus()
	return nil
}

func (p *Presente) SelecionarIntegral(idSelecao uuid.UUID) error {
	if p.tipo != TipoPresenteIntegral {
		return errors.New("operação válida apenas para presentes integrais")
	}

	if p.status == StatusSelecionado {
		return ErrPresenteJaSelecionado
	}

	p.status = StatusSelecionado
	return nil
}

func (p *Presente) LiberarSelecao(idSelecao uuid.UUID) error {
	if p.tipo == TipoPresenteIntegral {
		p.status = StatusDisponivel
		return nil
	}

	for _, cota := range p.cotas {
		if cota.idSelecao != nil && *cota.idSelecao == idSelecao {
			if err := cota.LiberarSelecao(); err != nil {
				return fmt.Errorf("erro ao liberar cota %d: %w", cota.numeroCota, err)
			}
		}
	}

	p.atualizarStatus()
	return nil
}

func (p *Presente) atualizarStatus() {
	if p.tipo == TipoPresenteIntegral {
		return
	}

	cotasDisponiveis := p.ContarCotasDisponiveis()
	totalCotas := len(p.cotas)

	if cotasDisponiveis == 0 {
		p.status = StatusSelecionado
	} else if cotasDisponiveis == totalCotas {
		p.status = StatusDisponivel
	} else {
		p.status = StatusParcialmenteSelecionado
	}
}

func (p *Presente) ContarCotasDisponiveis() int {
	if p.tipo == TipoPresenteIntegral {
		return 0
	}

	contador := 0
	for _, cota := range p.cotas {
		if cota.EstaDisponivel() {
			contador++
		}
	}
	return contador
}

func (p *Presente) ContarCotasSelecionadas() int {
	if p.tipo == TipoPresenteIntegral {
		if p.status == StatusSelecionado {
			return 1
		}
		return 0
	}

	contador := 0
	for _, cota := range p.cotas {
		if cota.EstaSelecionada() {
			contador++
		}
	}
	return contador
}

func (p *Presente) ObterValorCota() float64 {
	if p.tipo == TipoPresenteIntegral || len(p.cotas) == 0 {
		return 0
	}
	return p.cotas[0].valorCota
}

func (p *Presente) EhFracionado() bool {
	return p.tipo == TipoPresenteFracionado
}

func (p *Presente) EhIntegral() bool {
	return p.tipo == TipoPresenteIntegral
}

func (p *Presente) ID() uuid.UUID              { return p.id }
func (p *Presente) IDCasamento() uuid.UUID     { return p.idCasamento }
func (p *Presente) Nome() string               { return p.nome }
func (p *Presente) Descricao() string          { return p.descricao }
func (p *Presente) FotoURL() string            { return p.fotoURL }
func (p *Presente) EhFavorito() bool           { return p.ehFavorito }
func (p *Presente) Status() string             { return p.status }
func (p *Presente) Detalhes() DetalhesPresente { return p.detalhes }
func (p *Presente) Categoria() string          { return p.categoria }
func (p *Presente) Tipo() string               { return p.tipo }
func (p *Presente) ValorTotal() *float64       { return p.valorTotal }
func (p *Presente) Cotas() []*Cota             { return p.cotas }
