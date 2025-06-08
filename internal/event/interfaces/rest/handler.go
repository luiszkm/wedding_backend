// file: internal/event/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/event/application"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type EventHandler struct {
	service *application.EventService
}

func NewEventHandler(service *application.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) HandleCriarEvento(w http.ResponseWriter, r *http.Request) {
	// 1. Extrai o ID do usuário do contexto. O middleware garante que ele exista.
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}

	// 2. Decodifica o corpo da requisição para o DTO.
	var reqDTO CriarEventoRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "O corpo da requisição está malformado.", http.StatusBadRequest)
		return
	}

	// 3. Chama o serviço, passando o ID do usuário logado.
	novoEvento, err := h.service.CriarNovoEvento(
		r.Context(),
		userID,
		reqDTO.Nome,
		reqDTO.Data,
		reqDTO.Tipo,
		reqDTO.UrlSlug,
	)
	if err != nil {
		log.Printf("ERRO ao criar evento: %v", err)
		web.RespondError(w, r, "ERRO_CRIACAO_EVENTO", err.Error(), http.StatusBadRequest)
		return
	}

	// Conforme a documentação, retornamos 201 Created com o ID do evento.
	respDTO := CriarEventoResponseDTO{IDEvento: novoEvento.ID().String()}
	web.Respond(w, r, respDTO, http.StatusCreated)
}
