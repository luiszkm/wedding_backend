// file: internal/billing/interfaces/rest/handler.go
package rest

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, h.webhookSecret)
	if err != nil {
		log.Printf("ERRO na verificação da assinatura do webhook: %v", err)
		web.RespondError(w, r, "ASSINATURA_INVALIDA", "Assinatura do webhook inválida.", http.StatusBadRequest)
		return
	}

	// Processa os eventos que nos interessam.
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("ERRO ao decodificar o objeto da sessão: %v", err)
			w.WriteHeader(http.StatusOK) // Responde OK para a Stripe não reenviar
			return
		}

		if session.Subscription == nil {
			log.Printf("ERRO: evento checkout.session.completed sem dados de assinatura")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Usamos o ClientReferenceID, que é o nosso ID de assinatura (UUID).
		assinaturaID, err := uuid.Parse(session.ClientReferenceID)
		if err != nil {
			log.Printf("ERRO: ClientReferenceID inválido no evento da Stripe: %v", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		// E pegamos o ID da assinatura da Stripe (string 'sub_...').
		stripeSubscriptionID := session.Subscription.ID

		// Chama nosso serviço de aplicação para finalizar o processo de ativação.
		if err := h.service.AtivarAssinatura(r.Context(), assinaturaID, stripeSubscriptionID); err != nil {
			log.Printf("ERRO ao ativar assinatura %s: %v", assinaturaID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			log.Printf("ERRO ao decodificar o objeto de fatura: %v", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Verificamos se o pagamento é para um ciclo de assinatura (renovação)
		if invoice.BillingReason == stripe.InvoiceBillingReasonSubscriptionCycle {
			stripeSubID := invoice.ID
			// A fatura nos dá o novo final do período.
			novoFimPeriodo := time.Unix(invoice.PeriodEnd, 0)
			if err := h.service.RenovarAssinatura(r.Context(), stripeSubID, novoFimPeriodo); err != nil {
				log.Printf("ERRO ao renovar assinatura %s: %v", stripeSubID, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			log.Printf("ERRO ao decodificar o objeto de assinatura (delete): %v", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := h.service.CancelarAssinatura(r.Context(), subscription.ID); err != nil {
			log.Printf("ERRO ao cancelar assinatura %s: %v", subscription.ID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		log.Printf("Evento de webhook não tratado recebido: %s", event.Type)
	}

	// Para todos os eventos processados com sucesso (ou não tratados),
	// respondemos 200 OK para a Stripe.
	w.WriteHeader(http.StatusOK)
}
