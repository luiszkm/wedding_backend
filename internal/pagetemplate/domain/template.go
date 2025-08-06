// file: internal/pagetemplate/domain/template.go
package domain

import (
	"errors"
	"html/template"
	"strings"
	"time"

	"github.com/google/uuid"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	giftDomain "github.com/luiszkm/wedding_backend/internal/gift/domain"
	galleryDomain "github.com/luiszkm/wedding_backend/internal/gallery/domain"
	guestDomain "github.com/luiszkm/wedding_backend/internal/guest/domain"
	mbDomain "github.com/luiszkm/wedding_backend/internal/messageboard/domain"
)

type TipoTemplate string

const (
	TipoStandard TipoTemplate = "STANDARD"
	TipoBespoke  TipoTemplate = "BESPOKE"
)

var (
	ErrTemplateNaoEncontrado = errors.New("template não encontrado")
	ErrTemplateMalformado    = errors.New("template malformado")
	ErrDadosIncompletos      = errors.New("dados incompletos para renderização")
)

// TemplateMetadata contém os metadados de um template
type TemplateMetadata struct {
	ID               string
	Nome             string
	Descricao        string
	Tipo             TipoTemplate
	CaminhoArquivo   string
	PaletaDefault    eventDomain.PaletaCores
	SuportaGifts     bool
	SuportaGallery   bool
	SuportaMessages  bool
	SuportaRSVP      bool
	CriadoEm         time.Time
}

// EventPageData contém todos os dados necessários para renderizar uma página de evento
type EventPageData struct {
	Event        *eventDomain.Evento
	GuestGroups  []*guestDomain.GrupoDeConvidados
	Gifts        []*giftDomain.Presente
	Messages     []*mbDomain.Recado
	Photos       []*galleryDomain.Foto
	PaletaCores  eventDomain.PaletaCores
	
	// Flags de controle de exibição
	ShowGifts    bool
	ShowGallery  bool
	ShowMessages bool
	ShowRSVP     bool
	
	// Dados adicionais
	Contact      *ContactInfo
	CustomData   map[string]interface{}
}

// ContactInfo contém informações de contato para o footer
type ContactInfo struct {
	Nome     string
	Email    string
	Telefone string
}

// TemplateEngine define a interface para o motor de templates
type TemplateEngine interface {
	RenderEventPage(templateID string, data *EventPageData) ([]byte, error)
	ValidateTemplate(templateID string) error
	GetAvailableTemplates() ([]*TemplateMetadata, error)
	LoadTemplate(templateID string) (*template.Template, error)
}

// TemplateConfig representa uma configuração de template (precisa estar no domain)
type TemplateConfig struct {
	IsBespoke          bool                    `json:"is_bespoke"`
	StandardTemplateID string                  `json:"standard_template_id,omitempty"`
	BespokeFileName    string                  `json:"bespoke_file_name,omitempty"`
	PaletaCores        eventDomain.PaletaCores `json:"paleta_cores,omitempty"`
}

// TemplateService define a interface para o serviço de templates
type TemplateService interface {
	RenderPublicPage(urlSlug string) ([]byte, error)
	GetEventPageData(evento *eventDomain.Evento) (*EventPageData, error)
	GetTemplateMetadata(templateID string) (*TemplateMetadata, error)
	ListAvailableTemplates() ([]*TemplateMetadata, error)
	UpdateEventTemplate(eventID, userID uuid.UUID, config TemplateConfig) error
}

// Métodos de domínio para TemplateMetadata

func NewTemplateMetadata(id, nome, descricao string, tipo TipoTemplate, caminhoArquivo string) *TemplateMetadata {
	paletaDefault := eventDomain.PaletaCores{
		"primary":    "#2563eb",
		"secondary":  "#f1f5f9", 
		"accent":     "#10b981",
		"background": "#ffffff",
		"text":       "#1f2937",
	}

	return &TemplateMetadata{
		ID:               id,
		Nome:             nome,
		Descricao:        descricao,
		Tipo:             tipo,
		CaminhoArquivo:   caminhoArquivo,
		PaletaDefault:    paletaDefault,
		SuportaGifts:     true,
		SuportaGallery:   true,
		SuportaMessages:  true,
		SuportaRSVP:      true,
		CriadoEm:         time.Now(),
	}
}

func (tm *TemplateMetadata) IsValid() error {
	if tm.ID == "" {
		return errors.New("ID do template é obrigatório")
	}
	if tm.Nome == "" {
		return errors.New("nome do template é obrigatório")
	}
	if tm.CaminhoArquivo == "" {
		return errors.New("caminho do arquivo é obrigatório")
	}
	if tm.Tipo != TipoStandard && tm.Tipo != TipoBespoke {
		return errors.New("tipo de template inválido")
	}
	return nil
}

func (tm *TemplateMetadata) GetFullPath() string {
	if tm.Tipo == TipoBespoke {
		return "templates/bespoke/" + tm.CaminhoArquivo
	}
	return "templates/standard/" + tm.CaminhoArquivo
}

// Métodos de domínio para EventPageData

func NewEventPageData(evento *eventDomain.Evento) *EventPageData {
	return &EventPageData{
		Event:        evento,
		PaletaCores:  evento.PaletaCores(),
		GuestGroups:  []*guestDomain.GrupoDeConvidados{},
		Gifts:        []*giftDomain.Presente{},
		Messages:     []*mbDomain.Recado{},
		Photos:       []*galleryDomain.Foto{},
		ShowGifts:    true,
		ShowGallery:  true,
		ShowMessages: true,
		ShowRSVP:     true,
		CustomData:   make(map[string]interface{}),
	}
}

func (epd *EventPageData) AddGuestGroup(group *guestDomain.GrupoDeConvidados) {
	if group != nil {
		epd.GuestGroups = append(epd.GuestGroups, group)
	}
}

func (epd *EventPageData) AddGift(gift *giftDomain.Presente) {
	if gift != nil {
		epd.Gifts = append(epd.Gifts, gift)
	}
}

func (epd *EventPageData) AddMessage(message *mbDomain.Recado) {
	if message != nil {
		epd.Messages = append(epd.Messages, message)
	}
}

func (epd *EventPageData) AddPhoto(photo *galleryDomain.Foto) {
	if photo != nil {
		epd.Photos = append(epd.Photos, photo)
	}
}

func (epd *EventPageData) SetContact(nome, email, telefone string) {
	epd.Contact = &ContactInfo{
		Nome:     strings.TrimSpace(nome),
		Email:    strings.TrimSpace(email),
		Telefone: strings.TrimSpace(telefone),
	}
}

func (epd *EventPageData) SetCustomData(key string, value interface{}) {
	if epd.CustomData == nil {
		epd.CustomData = make(map[string]interface{})
	}
	epd.CustomData[key] = value
}

func (epd *EventPageData) GetCustomData(key string) interface{} {
	if epd.CustomData == nil {
		return nil
	}
	return epd.CustomData[key]
}

// Validação dos dados da página
func (epd *EventPageData) Validate() error {
	if epd.Event == nil {
		return ErrDadosIncompletos
	}
	if epd.Event.Nome() == "" {
		return ErrDadosIncompletos
	}
	if epd.PaletaCores == nil || len(epd.PaletaCores) == 0 {
		return ErrDadosIncompletos
	}
	return nil
}

// Templates padrão disponíveis
func GetStandardTemplates() []*TemplateMetadata {
	return []*TemplateMetadata{
		{
			ID:             "template_moderno",
			Nome:           "Moderno",
			Descricao:      "Template moderno e minimalista com design clean",
			Tipo:           TipoStandard,
			CaminhoArquivo: "template_moderno.html",
			PaletaDefault: eventDomain.PaletaCores{
				"primary":    "#2563eb",
				"secondary":  "#f1f5f9",
				"accent":     "#10b981",
				"background": "#ffffff",
				"text":       "#1f2937",
			},
			SuportaGifts:    true,
			SuportaGallery:  true,
			SuportaMessages: true,
			SuportaRSVP:     true,
		},
		{
			ID:             "template_classico",
			Nome:           "Clássico",
			Descricao:      "Template tradicional elegante com tipografia clássica",
			Tipo:           TipoStandard,
			CaminhoArquivo: "template_classico.html",
			PaletaDefault: eventDomain.PaletaCores{
				"primary":    "#8b5a3c",
				"secondary":  "#f5f5dc",
				"accent":     "#d4af37",
				"background": "#fdfdf8",
				"text":       "#2c1810",
			},
			SuportaGifts:    true,
			SuportaGallery:  true,
			SuportaMessages: true,
			SuportaRSVP:     true,
		},
		{
			ID:             "template_elegante",
			Nome:           "Elegante",
			Descricao:      "Template sofisticado com elementos luxuosos",
			Tipo:           TipoStandard,
			CaminhoArquivo: "template_elegante.html",
			PaletaDefault: eventDomain.PaletaCores{
				"primary":    "#1a1a2e",
				"secondary":  "#16213e",
				"accent":     "#e94560",
				"background": "#0f0f23",
				"text":       "#ffffff",
			},
			SuportaGifts:    true,
			SuportaGallery:  true,
			SuportaMessages: true,
			SuportaRSVP:     true,
		},
	}
}