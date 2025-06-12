// file: internal/billing/application/service.go
package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/billing/domain"
)

type BillingService struct {
	planoRepo      domain.PlanoRepository
	assinaturaRepo domain.AssinaturaRepository
	gateway        domain.PaymentGateway
}

func NewBillingService(planoRepo domain.PlanoRepository, assinaturaRepo domain.AssinaturaRepository, gateway domain.PaymentGateway) *BillingService {
	return &BillingService{
		planoRepo:      planoRepo,
		assinaturaRepo: assinaturaRepo,
		gateway:        gateway,
	}
}

// ListarPlanos é o caso de uso para buscar todos os planos disponíveis.
func (s *BillingService) ListarPlanos(ctx context.Context) ([]*domain.Plano, error) {
	planos, err := s.planoRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista de planos no serviço: %w", err)
	}
	return planos, nil
}

func (s *BillingService) IniciarNovaAssinatura(ctx context.Context, userID, planoID uuid.UUID) (string, *domain.Assinatura, error) {
	// 1. Valida se o plano escolhido existe.
	plano, err := s.planoRepo.FindByID(ctx, planoID)
	if err != nil {
		return "", nil, fmt.Errorf("plano selecionado é inválido: %w", err)
	}

	// 2. Cria a nova assinatura no estado PENDENTE.
	novaAssinatura := domain.NewAssinatura(userID, planoID)
	if err := s.assinaturaRepo.Save(ctx, novaAssinatura); err != nil {
		return "", nil, fmt.Errorf("falha ao criar registro de assinatura: %w", err)
	}

	// 3. Usa a interface do gateway para criar a sessão de checkout.
	// O serviço não sabe se é Stripe, Pagar.me ou outro.
	checkoutURL, err := s.gateway.CriarSessaoCheckout(ctx, novaAssinatura, plano)
	if err != nil {
		// Aqui poderíamos ter uma lógica para reverter a criação da assinatura
		return "", nil, fmt.Errorf("falha ao iniciar processo de pagamento: %w", err)
	}

	return checkoutURL, novaAssinatura, nil
}

func (s *BillingService) AtivarAssinatura(ctx context.Context, assinaturaID uuid.UUID, stripeSubscriptionID string) error {
	// 1. Busca os detalhes da assinatura no gateway de pagamento para obter as datas corretas.
	details, err := s.gateway.GetSubscriptionDetails(ctx, stripeSubscriptionID)
	if err != nil {
		return fmt.Errorf("falha ao obter detalhes do gateway de pagamento: %w", err)
	}

	// 2. Busca nossa assinatura interna no nosso banco.
	assinatura, err := s.assinaturaRepo.FindByID(ctx, assinaturaID)
	if err != nil {
		return fmt.Errorf("falha ao buscar assinatura para ativação: %w", err)
	}

	// 3. Validação de negócio: só ativamos se estiver pendente.
	if assinatura.Status() != domain.StatusPendente {
		log.Printf("Aviso: tentativa de ativar assinatura %s que já não está pendente.", assinaturaID)
		return nil
	}

	// 4. Ativa a assinatura no nosso domínio com os dados obtidos.
	assinatura.Ativar(details.ID, details.CurrentPeriodStart, details.CurrentPeriodEnd)

	// 5. Salva o novo estado.
	if err := s.assinaturaRepo.Update(ctx, assinatura); err != nil {
		return fmt.Errorf("falha ao atualizar assinatura para ativa: %w", err)
	}

	log.Printf("Assinatura %s ativada com sucesso. Válida até %s", assinatura.ID(), details.CurrentPeriodEnd.Format("02/01/2006"))
	return nil
}

// RenovarAssinatura atualiza a data de fim de uma assinatura após um pagamento de renovação.
func (s *BillingService) RenovarAssinatura(ctx context.Context, stripeSubID string, novoFimPeriodo time.Time) error {
	assinatura, err := s.assinaturaRepo.FindByStripeSubscriptionID(ctx, stripeSubID)
	if err != nil {
		return fmt.Errorf("assinatura com id stripe %s não encontrada para renovação: %w", stripeSubID, err)
	}

	assinatura.Renovar(novoFimPeriodo)

	if err := s.assinaturaRepo.Update(ctx, assinatura); err != nil {
		return fmt.Errorf("falha ao renovar assinatura: %w", err)
	}

	log.Printf("Assinatura com ID Stripe %s renovada. Nova validade: %s", stripeSubID, novoFimPeriodo.Format("02/01/2006"))
	return nil
}

// CancelarAssinatura atualiza o status de uma assinatura para CANCELADA.
func (s *BillingService) CancelarAssinatura(ctx context.Context, stripeSubID string) error {
	assinatura, err := s.assinaturaRepo.FindByStripeSubscriptionID(ctx, stripeSubID)
	if err != nil {
		return fmt.Errorf("assinatura com id stripe %s não encontrada para cancelamento: %w", stripeSubID, err)
	}

	assinatura.Cancelar()

	if err := s.assinaturaRepo.Update(ctx, assinatura); err != nil {
		return fmt.Errorf("falha ao cancelar assinatura: %w", err)
	}

	log.Printf("Assinatura com ID Stripe %s foi cancelada.", stripeSubID)
	return nil
}
