// file: internal/event/domain/evento.go
package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TipoEvento string

const (
	TipoCasamento   TipoEvento = "CASAMENTO"
	TipoAniversario TipoEvento = "ANIVERSARIO"
	TipoChaDeBebe   TipoEvento = "CHA_DE_BEBE"
	TipoOutro       TipoEvento = "OUTRO"
)

var (
	ErrEventoNaoEncontrado = errors.New("evento não encontrado")
	ErrTipoEventoInvalido  = errors.New("tipo de evento inválido")
	ErrSlugEmUso           = errors.New("a URL amigável (slug) já está em uso")
	ErrEventoJaExiste      = errors.New("evento já existe")
	ErrTemplateInvalido    = errors.New("template inválido")
	ErrPaletaCoresInvalida = errors.New("paleta de cores inválida")
)

func (t TipoEvento) IsValid() bool {
	switch t {
	case TipoCasamento, TipoAniversario, TipoChaDeBebe, TipoOutro:
		return true
	}
	return false
}

type PaletaCores map[string]string

type Evento struct {
	id                uuid.UUID
	idUsuario         uuid.UUID
	nome              string
	data              time.Time
	tipo              TipoEvento
	urlSlug           string
	idTemplate        string
	idTemplateArquivo *string
	paletaCores       PaletaCores
}

func NewEvento(idUsuario uuid.UUID, nome string, data time.Time, tipo TipoEvento, urlSlug string) (*Evento, error) {
	if nome == "" || urlSlug == "" {
		return nil, errors.New("nome e urlSlug são obrigatórios")
	}
	if !tipo.IsValid() {
		return nil, errors.New("tipo de evento inválido")
	}
	
	// Paleta de cores padrão
	paletaDefault := PaletaCores{
		"primary":    "#2563eb",
		"secondary":  "#f1f5f9",
		"accent":     "#10b981",
		"background": "#ffffff",
		"text":       "#1f2937",
	}
	
	return &Evento{
		id:                uuid.New(),
		idUsuario:         idUsuario,
		nome:              strings.TrimSpace(nome),
		data:              data,
		tipo:              tipo,
		urlSlug:           strings.TrimSpace(urlSlug),
		idTemplate:        "template_moderno",
		idTemplateArquivo: nil,
		paletaCores:       paletaDefault,
	}, nil
}

// HydrateEvento cria uma nova instância de Evento a partir dos dados fornecidos.
func HydrateEvento(id, idUsuario uuid.UUID, nome string, data time.Time, tipo TipoEvento, urlSlug string, idTemplate string, idTemplateArquivo *string, paletaCores PaletaCores) *Evento {
	if nome == "" || urlSlug == "" {
		return nil // ou retornar um erro, dependendo da lógica de negócio
	}
	if !tipo.IsValid() {
		return nil // ou retornar um erro, dependendo da lógica de negócio
	}
	
	// Se paleta de cores não foi fornecida, usar padrão
	if paletaCores == nil {
		paletaCores = PaletaCores{
			"primary":    "#2563eb",
			"secondary":  "#f1f5f9",
			"accent":     "#10b981",
			"background": "#ffffff",
			"text":       "#1f2937",
		}
	}
	
	// Se template não foi especificado, usar padrão
	if idTemplate == "" {
		idTemplate = "template_moderno"
	}
	
	return &Evento{
		id:                id,
		idUsuario:         idUsuario,
		nome:              strings.TrimSpace(nome),
		data:              data,
		tipo:              tipo,
		urlSlug:           strings.TrimSpace(urlSlug),
		idTemplate:        idTemplate,
		idTemplateArquivo: idTemplateArquivo,
		paletaCores:       paletaCores,
	}
}

// Getters
func (e *Evento) ID() uuid.UUID        { return e.id }
func (e *Evento) IDUsuario() uuid.UUID { return e.idUsuario }
func (e *Evento) Nome() string         { return e.nome }
func (e *Evento) Data() time.Time      { return e.data }
func (e *Evento) Tipo() TipoEvento     { return e.tipo }
func (e *Evento) UrlSlug() string      { return e.urlSlug }
func (e *Evento) IDTemplate() string   { return e.idTemplate }
func (e *Evento) IDTemplateArquivo() *string { return e.idTemplateArquivo }
func (e *Evento) PaletaCores() PaletaCores   { return e.paletaCores }

// Métodos de negócio para templates

// UsaTemplateBespoke retorna true se o evento usa um template personalizado
func (e *Evento) UsaTemplateBespoke() bool {
	return e.idTemplateArquivo != nil && *e.idTemplateArquivo != ""
}

// GetTemplateAtivo retorna o template que deve ser usado (bespoke tem precedência)
func (e *Evento) GetTemplateAtivo() string {
	if e.UsaTemplateBespoke() {
		return *e.idTemplateArquivo
	}
	return e.idTemplate
}

// DefinirTemplate altera o template padrão do evento
func (e *Evento) DefinirTemplate(idTemplate string) error {
	if idTemplate == "" {
		return ErrTemplateInvalido
	}
	
	// Validar se é um template padrão válido
	templatesValidos := []string{"template_moderno", "template_classico", "template_elegante"}
	valido := false
	for _, t := range templatesValidos {
		if t == idTemplate {
			valido = true
			break
		}
	}
	
	if !valido {
		return ErrTemplateInvalido
	}
	
	e.idTemplate = idTemplate
	// Limpar template bespoke se existir
	e.idTemplateArquivo = nil
	return nil
}

// DefinirTemplateBespoke define um template personalizado
func (e *Evento) DefinirTemplateBespoke(nomeArquivo string) error {
	if nomeArquivo == "" {
		return ErrTemplateInvalido
	}
	
	// Validar formato do nome do arquivo
	if !strings.HasSuffix(nomeArquivo, ".html") {
		nomeArquivo += ".html"
	}
	
	e.idTemplateArquivo = &nomeArquivo
	return nil
}

// DefinirPaletaCores atualiza a paleta de cores do evento
func (e *Evento) DefinirPaletaCores(paleta PaletaCores) error {
	if paleta == nil {
		return ErrPaletaCoresInvalida
	}
	
	// Validar cores obrigatórias
	coresObrigatorias := []string{"primary", "secondary", "background", "text"}
	for _, cor := range coresObrigatorias {
		if _, existe := paleta[cor]; !existe {
			return ErrPaletaCoresInvalida
		}
	}
	
	e.paletaCores = paleta
	return nil
}

// PaletaCoresJSON retorna a paleta de cores em formato JSON para o banco
func (e *Evento) PaletaCoresJSON() ([]byte, error) {
	return json.Marshal(e.paletaCores)
}

// SetPaletaCoresFromJSON define a paleta a partir de JSON
func (e *Evento) SetPaletaCoresFromJSON(jsonData []byte) error {
	var paleta PaletaCores
	if err := json.Unmarshal(jsonData, &paleta); err != nil {
		return err
	}
	return e.DefinirPaletaCores(paleta)
}
