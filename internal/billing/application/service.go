// file: internal/billing/application/service.go
package application

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/billing/domain"
)

type BillingService struct {
	planoRepo      domain.PlanoRepository
	assinaturaRepo domain.AssinaturaRepository
}

func NewBillingService(planoRepo domain.PlanoRepository, assinaturaRepo domain.AssinaturaRepository) *BillingService {
	return &BillingService{planoRepo: planoRepo, assinaturaRepo: assinaturaRepo}
}

// ListarPlanos é o caso de uso para buscar todos os planos disponíveis.
func (s *BillingService) ListarPlanos(ctx context.Context) ([]*domain.Plano, error) {
	planos, err := s.planoRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista de planos no serviço: %w", err)
	}
	return planos, nil
}

func (s *BillingService) IniciarNovaAssinatura(ctx context.Context, userID, planoID uuid.UUID) (*domain.Assinatura, error) {
	// 1. Valida se o plano escolhido existe.
	_, err := s.planoRepo.FindByID(ctx, planoID)
	if err != nil {
		return nil, fmt.Errorf("plano selecionado é inválido: %w", err)
	}

	// (Aqui entraria a lógica de verificar se o usuário já tem uma assinatura ativa)

	// 2. Cria a nova assinatura no estado PENDENTE.
	novaAssinatura := domain.NewAssinatura(userID, planoID)

	// 3. Salva no banco de dados.
	if err := s.assinaturaRepo.Save(ctx, novaAssinatura); err != nil {
		return nil, fmt.Errorf("falha ao criar registro de assinatura: %w", err)
	}

	// 4. (Simulação) Aqui dispararíamos o processo de pagamento.
	log.Printf("Processo de pagamento iniciado para a assinatura %s do usuário %s", novaAssinatura.ID(), userID)

	return novaAssinatura, nil
}
