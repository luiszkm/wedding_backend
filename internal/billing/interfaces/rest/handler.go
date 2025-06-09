// file: internal/billing/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/billing/application"
	"github.com/luiszkm/wedding_backend/internal/billing/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
)

type BillingHandler struct {
	service *application.BillingService
}

func NewBillingHandler(service *application.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) HandleListarPlanos(w http.ResponseWriter, r *http.Request) {
	// 1. Chama o serviço de aplicação.
	planos, err := h.service.ListarPlanos(r.Context())
	if err != nil {
		log.Printf("ERRO ao listar planos: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar a lista de planos.", http.StatusInternalServerError)
		return
	}

	// 2. Mapeia os objetos de domínio para DTOs de resposta.
	respDTO := make([]PlanoDTO, len(planos))
	for i, p := range planos {
		respDTO[i] = PlanoDTO{
			ID:                  p.ID().String(),
			Nome:                p.Nome(),
			PrecoEmCentavos:     p.PrecoEmCentavos(),
			NumeroMaximoEventos: p.NumeroMaximoEventos(),
			DuracaoEmDias:       p.DuracaoEmDias(),
		}
	}

	// 3. Responde com sucesso.
	web.Respond(w, r, respDTO, http.StatusOK)
}

func (h *BillingHandler) HandleCriarAssinatura(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserContextKey).(uuid.UUID)
	if !ok {
		web.RespondError(w, r, "TOKEN_INVALIDO", "ID de usuário ausente no token.", http.StatusUnauthorized)
		return
	}

	var reqDTO CriarAssinaturaRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		web.RespondError(w, r, "CORPO_INVALIDO", "Corpo da requisição malformado.", http.StatusBadRequest)
		return
	}

	planoID, err := uuid.Parse(reqDTO.IDPlano)
	if err != nil {
		web.RespondError(w, r, "DADOS_INVALIDOS", "O ID do plano é inválido.", http.StatusBadRequest)
		return
	}

	novaAssinatura, err := h.service.IniciarNovaAssinatura(r.Context(), userID, planoID)
	if err != nil {
		if errors.Is(err, domain.ErrPlanoNaoEncontrado) {
			web.RespondError(w, r, "PLANO_NAO_ENCONTRADO", err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("ERRO ao criar assinatura: %v", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao processar assinatura.", http.StatusInternalServerError)
		return
	}

	respDTO := CriarAssinaturaResponseDTO{
		IDAssinatura: novaAssinatura.ID().String(),
		Status:       string(novaAssinatura.Status()),
	}

	// Conforme a documentação, retornamos 202 Accepted.
	web.Respond(w, r, respDTO, http.StatusAccepted)
}
