// file: internal/pagetemplate/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	"github.com/luiszkm/wedding_backend/internal/pagetemplate/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

// PageTemplateHandler gerencia requests relacionados a templates de página
type PageTemplateHandler struct {
	templateService domain.TemplateService
}

// NewPageTemplateHandler cria um novo handler para templates de página
func NewPageTemplateHandler(templateService domain.TemplateService) *PageTemplateHandler {
	return &PageTemplateHandler{
		templateService: templateService,
	}
}

// HandleRenderPublicPage renderiza a página pública de um evento
func (h *PageTemplateHandler) HandleRenderPublicPage(w http.ResponseWriter, r *http.Request) {
	urlSlug := chi.URLParam(r, "urlSlug")
	if urlSlug == "" {
		web.RespondError(w, r, "SLUG_OBRIGATORIO", "URL slug é obrigatório", http.StatusBadRequest)
		return
	}

	// Renderizar página
	htmlContent, err := h.templateService.RenderPublicPage(urlSlug)
	if err != nil {
		if err == eventDomain.ErrEventoNaoEncontrado {
			web.RespondError(w, r, "EVENTO_NAO_ENCONTRADO", "Evento não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_RENDERIZACAO", "Erro ao renderizar página: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Definir headers para HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300") // Cache de 5 minutos
	w.WriteHeader(http.StatusOK)
	w.Write(htmlContent)
}

// HandleListAvailableTemplates lista todos os templates disponíveis
func (h *PageTemplateHandler) HandleListAvailableTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.templateService.ListAvailableTemplates()
	if err != nil {
		web.RespondError(w, r, "ERRO_LISTAGEM_TEMPLATES", "Erro ao listar templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateTemplatesListResponse(templates)
	web.Respond(w, r, response, http.StatusOK)
}

// HandleUpdateEventTemplate atualiza o template de um evento
func (h *PageTemplateHandler) HandleUpdateEventTemplate(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do evento da URL
	eventIDStr := chi.URLParam(r, "eventId")
	if eventIDStr == "" {
		web.RespondError(w, r, "EVENT_ID_OBRIGATORIO", "ID do evento é obrigatório", http.StatusBadRequest)
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		web.RespondError(w, r, "EVENT_ID_INVALIDO", "ID do evento é inválido", http.StatusBadRequest)
		return
	}

	// Extrair ID do usuário do contexto (adicionado pelo middleware de auth)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "USER_ID_OBRIGATORIO", "ID do usuário é obrigatório", http.StatusUnauthorized)
		return
	}

	// Parse do body da requisição
	var req TemplateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "JSON_INVALIDO", "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar requisição
	if err := req.Validate(); err != nil {
		web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
		return
	}

	// Converter para config do domain
	config := req.ToDomainTemplateConfig()

	// Atualizar template
	if err := h.templateService.UpdateEventTemplate(eventID, userID, config); err != nil {
		if err == eventDomain.ErrEventoNaoEncontrado {
			web.RespondError(w, r, "EVENTO_NAO_ENCONTRADO", "Evento não encontrado", http.StatusNotFound)
			return
		}
		if err == domain.ErrTemplateNaoEncontrado {
			web.RespondError(w, r, "TEMPLATE_NAO_ENCONTRADO", "Template não encontrado", http.StatusBadRequest)
			return
		}
		web.RespondError(w, r, "ERRO_ATUALIZACAO", "Erro ao atualizar template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := SuccessResponse{
		Message: "Template atualizado com sucesso",
	}
	web.Respond(w, r, response, http.StatusOK)
}

// HandleGetEventTemplateInfo retorna informações sobre o template atual de um evento
func (h *PageTemplateHandler) HandleGetEventTemplateInfo(w http.ResponseWriter, r *http.Request) {
	// Por enquanto, retorna não implementado
	// TODO: Implementar quando necessário
	web.RespondError(w, r, "NAO_IMPLEMENTADO", "Endpoint ainda não implementado", http.StatusNotImplemented)
}

// HandlePreviewTemplate gera uma prévia de template
func (h *PageTemplateHandler) HandlePreviewTemplate(w http.ResponseWriter, r *http.Request) {
	// Por enquanto, retorna não implementado
	// TODO: Implementar preview quando necessário
	web.RespondError(w, r, "NAO_IMPLEMENTADO", "Funcionalidade de prévia ainda não implementada", http.StatusNotImplemented)
}

// HandleGetTemplateMetadata retorna metadados de um template específico
func (h *PageTemplateHandler) HandleGetTemplateMetadata(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "templateId")
	if templateID == "" {
		web.RespondError(w, r, "TEMPLATE_ID_OBRIGATORIO", "ID do template é obrigatório", http.StatusBadRequest)
		return
	}

	// Decodificar URL se necessário
	templateID = strings.ReplaceAll(templateID, "%2F", "/")

	metadata, err := h.templateService.GetTemplateMetadata(templateID)
	if err != nil {
		if err == domain.ErrTemplateNaoEncontrado {
			web.RespondError(w, r, "TEMPLATE_NAO_ENCONTRADO", "Template não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_BUSCA_METADATA", "Erro ao buscar metadados: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := FromDomainTemplateMetadata(metadata)
	web.Respond(w, r, response, http.StatusOK)
}

// Middleware para validar Content-Type em requests que precisam de JSON
func RequireJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				web.RespondError(w, r, "CONTENT_TYPE_INVALIDO", "Content-Type deve ser application/json", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware para adicionar headers de cache em respostas de templates
func TemplateCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Para listagens de templates, cache por 1 hora
		if strings.Contains(r.URL.Path, "/templates") && r.Method == "GET" {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}