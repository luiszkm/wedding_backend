// file: internal/billing/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/billing/application"
	"github.com/luiszkm/wedding_backend/internal/billing/domain"
	"github.com/luiszkm/wedding_backend/internal/platform/auth"
	"github.com/luiszkm/wedding_backend/internal/platform/web"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type BillingHandler struct {
	service       *application.BillingService
	webhookSecret string
}

func NewBillingHandler(service *application.BillingService, webhookSecret string) *BillingHandler {
	return &BillingHandler{
		service:       service,
		webhookSecret: webhookSecret,
	}
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

	checkoutURL, novaAssinatura, err := h.service.IniciarNovaAssinatura(r.Context(), userID, planoID)
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
		CheckoutURL:  checkoutURL,
	}

	// Conforme a documentação, retornamos 202 Accepted.
	web.Respond(w, r, respDTO, http.StatusAccepted)
}

var webhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

func (h *BillingHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERRO ao ler o corpo do webhook: %v", err)
		// Não podemos processar, mas não é culpa da Stripe. Não retornamos erro aqui.
		// Apenas logamos e saímos. Um status 200 implícito será enviado.
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")

	// Usa o segredo injetado na struct do handler
	event, err := webhook.ConstructEvent(payload, signatureHeader, h.webhookSecret)
	if err != nil {
		log.Printf("ERRO na verificação da assinatura do webhook: %v", err)
		// A assinatura é inválida, requisição maliciosa ou mal configurada. Rejeitamos.
		web.RespondError(w, r, "ASSINATURA_INVALIDA", "Assinatura do webhook inválida.", http.StatusBadRequest)
		return
	}

	// Processa apenas os eventos que nos interessam.
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("ERRO ao decodificar o objeto da sessão de checkout: %v", err)
			// O evento está malformado, não adianta a Stripe reenviar. Respondemos OK.
			w.WriteHeader(http.StatusOK)
			return
		}

		// Usamos o ClientReferenceID, que é o nosso ID de assinatura.
		assinaturaID, err := uuid.Parse(session.ClientReferenceID)
		if err != nil {
			log.Printf("ERRO: ClientReferenceID inválido no evento da Stripe: %v", err)
			w.WriteHeader(http.StatusOK) // Mesmo caso do anterior.
			return
		}

		// Chama nosso serviço de aplicação para finalizar o processo.
		if err := h.service.AtivarAssinatura(r.Context(), assinaturaID); err != nil {
			log.Printf("ERRO ao ativar assinatura %s: %v", assinaturaID, err)
			// ESTE é um erro do nosso lado (ex: banco de dados).
			// Retornamos 500 para que a Stripe tente reenviar o webhook mais tarde.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		log.Printf("Evento de webhook não tratado recebido: %s", event.Type)
	}

	// Para todos os eventos recebidos com sucesso (mesmo os não tratados),
	// respondemos 200 OK para a Stripe.
	w.WriteHeader(http.StatusOK)
}
