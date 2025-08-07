package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/itinerary/application"
	"github.com/luiszkm/wedding_backend/internal/itinerary/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type ItineraryHandler struct {
	service *application.ItineraryService
}

func NewItineraryHandler(service *application.ItineraryService) *ItineraryHandler {
	return &ItineraryHandler{service: service}
}

// HandleCreateItineraryItem cria um novo item do roteiro (autenticado)
// POST /eventos/{idEvento}/roteiro
func (h *ItineraryHandler) HandleCreateItineraryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente ou inválido no token.", http.StatusUnauthorized)
		return
	}

	eventIDStr := chi.URLParam(r, "idEvento")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do evento é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO CreateItineraryItemRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	itemID, err := h.service.CreateItineraryItem(
		r.Context(),
		userID,
		eventID,
		reqDTO.Horario,
		reqDTO.TituloAtividade,
		reqDTO.DescricaoAtividade,
	)
	if err != nil {
		// Mapeia erros do domínio para respostas HTTP apropriadas
		if errors.Is(err, domain.ErrTituloAtividadeObrigatorio) ||
			errors.Is(err, domain.ErrTituloAtividadeMuitoLongo) ||
			errors.Is(err, domain.ErrHorarioObrigatorio) ||
			errors.Is(err, domain.ErrEventoObrigatorio) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao criar item do roteiro: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	respDTO := CreateItineraryItemResponseDTO{ID: itemID.String()}
	web.Respond(w, r, respDTO, http.StatusCreated)
}

// HandleGetItinerary retorna todos os itens do roteiro de um evento (público)
// GET /eventos/{idEvento}/roteiro
func (h *ItineraryHandler) HandleGetItinerary(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "idEvento")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do evento é inválido.", http.StatusBadRequest)
		return
	}

	items, err := h.service.GetItineraryByEventID(r.Context(), eventID)
	if err != nil {
		log.Printf("ERRO ao obter roteiro: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	itemsDTO := make([]ItineraryItemDTO, len(items))
	for i, item := range items {
		itemsDTO[i] = ItineraryItemDTO{
			ID:                item.ID().String(),
			IDEvento:          item.IDEvento().String(),
			Horario:           item.Horario(),
			TituloAtividade:   item.TituloAtividade(),
			DescricaoAtividade: item.DescricaoAtividade(),
			CreatedAt:         item.CreatedAt(),
			UpdatedAt:         item.UpdatedAt(),
		}
	}

	respDTO := GetItineraryResponseDTO{Items: itemsDTO}
	web.Respond(w, r, respDTO, http.StatusOK)
}

// HandleUpdateItineraryItem atualiza um item do roteiro (autenticado)
// PUT /roteiro/{idItemRoteiro}
func (h *ItineraryHandler) HandleUpdateItineraryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente ou inválido no token.", http.StatusUnauthorized)
		return
	}

	itemIDStr := chi.URLParam(r, "idItemRoteiro")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do item do roteiro é inválido.", http.StatusBadRequest)
		return
	}

	var reqDTO UpdateItineraryItemRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateItineraryItem(
		r.Context(),
		userID,
		itemID,
		reqDTO.Horario,
		reqDTO.TituloAtividade,
		reqDTO.DescricaoAtividade,
	)
	if err != nil {
		if errors.Is(err, domain.ErrItemRoteiroNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Item do roteiro não encontrado.", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTituloAtividadeObrigatorio) ||
			errors.Is(err, domain.ErrTituloAtividadeMuitoLongo) ||
			errors.Is(err, domain.ErrHorarioObrigatorio) {
			web.RespondError(w, r, "DADOS_INVALIDOS", err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERRO ao atualizar item do roteiro: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleDeleteItineraryItem remove um item do roteiro (autenticado)
// DELETE /roteiro/{idItemRoteiro}
func (h *ItineraryHandler) HandleDeleteItineraryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente ou inválido no token.", http.StatusUnauthorized)
		return
	}

	itemIDStr := chi.URLParam(r, "idItemRoteiro")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		web.RespondError(w, r, "PARAMETRO_INVALIDO", "O ID do item do roteiro é inválido.", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteItineraryItem(r.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrItemRoteiroNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Item do roteiro não encontrado.", http.StatusNotFound)
			return
		}
		log.Printf("ERRO ao deletar item do roteiro: %v\n", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar sua requisição.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}