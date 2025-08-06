// file: internal/pagetemplate/interfaces/rest/dto.go
package rest

import (
	"fmt"
	"time"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	"github.com/luiszkm/wedding_backend/internal/pagetemplate/domain"
)

// TemplateMetadataDTO representa metadados de um template para a API
type TemplateMetadataDTO struct {
	ID               string                      `json:"id"`
	Nome             string                      `json:"nome"`
	Descricao        string                      `json:"descricao"`
	Tipo             string                      `json:"tipo"`
	PaletaDefault    eventDomain.PaletaCores     `json:"paleta_default"`
	SuportaGifts     bool                        `json:"suporta_gifts"`
	SuportaGallery   bool                        `json:"suporta_gallery"`
	SuportaMessages  bool                        `json:"suporta_messages"`
	SuportaRSVP      bool                        `json:"suporta_rsvp"`
	CriadoEm         time.Time                   `json:"criado_em"`
}

// TemplateConfigRequest representa uma requisição para atualizar configuração de template
type TemplateConfigRequest struct {
	IsBespoke          bool                        `json:"is_bespoke"`
	StandardTemplateID string                      `json:"standard_template_id,omitempty"`
	BespokeFileName    string                      `json:"bespoke_file_name,omitempty"`
	PaletaCores        eventDomain.PaletaCores     `json:"paleta_cores,omitempty"`
}

// TemplatePreviewRequest representa uma requisição de prévia de template
type TemplatePreviewRequest struct {
	TemplateID string                      `json:"template_id"`
	Config     TemplateConfigRequest       `json:"config"`
}

// TemplatesListResponse representa a resposta da listagem de templates
type TemplatesListResponse struct {
	Templates []TemplateMetadataDTO `json:"templates"`
	Total     int                   `json:"total"`
}

// EventTemplateInfoResponse representa informações do template atual de um evento
type EventTemplateInfoResponse struct {
	EventID            string                      `json:"event_id"`
	EventName          string                      `json:"event_name"`
	CurrentTemplate    TemplateMetadataDTO         `json:"current_template"`
	IsUsingBespoke     bool                        `json:"is_using_bespoke"`
	BespokeFileName    *string                     `json:"bespoke_file_name,omitempty"`
	PaletaCores        eventDomain.PaletaCores     `json:"paleta_cores"`
	PublicURL          string                      `json:"public_url"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse representa uma resposta de sucesso
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Funções de conversão

// FromDomainTemplateMetadata converte de domain para DTO
func FromDomainTemplateMetadata(tmpl *domain.TemplateMetadata) TemplateMetadataDTO {
	return TemplateMetadataDTO{
		ID:               tmpl.ID,
		Nome:             tmpl.Nome,
		Descricao:        tmpl.Descricao,
		Tipo:             string(tmpl.Tipo),
		PaletaDefault:    tmpl.PaletaDefault,
		SuportaGifts:     tmpl.SuportaGifts,
		SuportaGallery:   tmpl.SuportaGallery,
		SuportaMessages:  tmpl.SuportaMessages,
		SuportaRSVP:      tmpl.SuportaRSVP,
		CriadoEm:         tmpl.CriadoEm,
	}
}

// ToDomainTemplateConfig converte DTO para domain config
func (req *TemplateConfigRequest) ToDomainTemplateConfig() domain.TemplateConfig {
	return domain.TemplateConfig{
		IsBespoke:          req.IsBespoke,
		StandardTemplateID: req.StandardTemplateID,
		BespokeFileName:    req.BespokeFileName,
		PaletaCores:        req.PaletaCores,
	}
}

// Validate valida uma requisição de configuração de template
func (req *TemplateConfigRequest) Validate() error {
	if req.IsBespoke {
		if req.BespokeFileName == "" {
			return fmt.Errorf("bespoke_file_name é obrigatório quando is_bespoke é true")
		}
	} else {
		if req.StandardTemplateID == "" {
			return fmt.Errorf("standard_template_id é obrigatório quando is_bespoke é false")
		}
	}

	// Validar formato das cores se fornecidas
	if req.PaletaCores != nil {
		coresObrigatorias := []string{"primary", "secondary", "background", "text"}
		for _, cor := range coresObrigatorias {
			if valor, existe := req.PaletaCores[cor]; !existe || valor == "" {
				return fmt.Errorf("cor obrigatória não encontrada ou vazia: %s", cor)
			}
		}
	}

	return nil
}

// CreateTemplatesListResponse cria uma resposta de listagem de templates
func CreateTemplatesListResponse(templates []*domain.TemplateMetadata) TemplatesListResponse {
	templateDTOs := make([]TemplateMetadataDTO, len(templates))
	for i, tmpl := range templates {
		templateDTOs[i] = FromDomainTemplateMetadata(tmpl)
	}

	return TemplatesListResponse{
		Templates: templateDTOs,
		Total:     len(templateDTOs),
	}
}

// CreateEventTemplateInfoResponse cria resposta com informações do template do evento
func CreateEventTemplateInfoResponse(
	evento *eventDomain.Evento,
	templateMetadata *domain.TemplateMetadata,
	baseURL string,
) EventTemplateInfoResponse {
	return EventTemplateInfoResponse{
		EventID:            evento.ID().String(),
		EventName:          evento.Nome(),
		CurrentTemplate:    FromDomainTemplateMetadata(templateMetadata),
		IsUsingBespoke:     evento.UsaTemplateBespoke(),
		BespokeFileName:    evento.IDTemplateArquivo(),
		PaletaCores:        evento.PaletaCores(),
		PublicURL:          fmt.Sprintf("%s/v1/eventos/%s/pagina", baseURL, evento.UrlSlug()),
	}
}