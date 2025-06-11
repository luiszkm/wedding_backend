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

func (s *BillingService) AtivarAssinatura(ctx context.Context, assinaturaID uuid.UUID) error {
	// 1. Busca a assinatura em nosso banco de dados.
	assinatura, err := s.assinaturaRepo.FindByID(ctx, assinaturaID)
	if err != nil {
		return fmt.Errorf("falha ao buscar assinatura para ativação: %w", err)
	}

	// 2. Validação de negócio: só ativamos se estiver pendente.
	if assinatura.Status() != domain.StatusPendente {
		log.Printf("Aviso: tentativa de ativar assinatura %s que já não está pendente.", assinaturaID)
		return nil // Não é um erro, apenas ignoramos.
	}

	// 3. Busca o plano para saber a duração
	plano, err := s.planoRepo.FindByID(ctx, assinatura.IDPlano())
	if err != nil {
		return fmt.Errorf("falha ao buscar plano da assinatura: %w", err)
	}

	// 4. Define as datas e ativa a assinatura no nosso domínio.
	dataInicio := time.Now()
	dataFim := dataInicio.AddDate(0, 0, plano.DuracaoEmDias())
	assinatura.Ativar(dataInicio, dataFim)

	// 5. Salva o novo estado no banco de dados.
	if err := s.assinaturaRepo.Update(ctx, assinatura); err != nil {
		return fmt.Errorf("falha ao atualizar assinatura para ativa: %w", err)
	}

	log.Printf("Assinatura %s ativada com sucesso. Válida até %s", assinatura.ID(), dataFim.Format("02/01/2006"))
	return nil
}
