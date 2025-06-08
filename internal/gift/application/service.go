// file: internal/gift/application/service.go
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	"github.com/luiszkm/wedding_backend/internal/gift/domain"
)

type GiftService struct {
	repo        domain.PresenteRepository
	selecaoRepo domain.SelecaoRepository
	eventRepo   eventDomain.EventoRepository // <-- Nova dependência

}

func NewGiftService(presenteRepo domain.PresenteRepository, selecaoRepo domain.SelecaoRepository, eventRepo eventDomain.EventoRepository) *GiftService {
	return &GiftService{repo: presenteRepo, selecaoRepo: selecaoRepo, eventRepo: eventRepo}
}

func (s *GiftService) CriarNovoPresente(ctx context.Context, userID, idEvento uuid.UUID, nome, desc, fotoURL, categoria string, favorito bool, detalhes domain.DetalhesPresente) (*domain.Presente, error) {
	// 1. AUTORIZAÇÃO: Verifica se o usuário logado é o dono do evento.
	_, err := s.eventRepo.FindByID(ctx, userID, idEvento)
	if err != nil {
		// Retorna o erro do repositório (ex: não encontrado), que o handler traduzirá para 403 ou 404.
		return nil, fmt.Errorf("permissão negada ou evento não encontrado: %w", err)
	}

	// 2. LÓGICA DE NEGÓCIO: Se a autorização passou, cria o presente.
	novoPresente, err := domain.NewPresente(idEvento, nome, desc, fotoURL, favorito, categoria, detalhes)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, novoPresente); err != nil {
		return nil, fmt.Errorf("falha ao salvar novo presente: %w", err)
	}

	return novoPresente, nil
}
func (s *GiftService) ListarPresentesDisponiveis(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Presente, error) {
	presentes, err := s.repo.ListarDisponiveisPorCasamento(ctx, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista de presentes: %w", err)
	}
	return presentes, nil
}
func (s *GiftService) FinalizarSelecaoDePresentes(ctx context.Context, chaveDeAcesso string, idsDosPresentes []uuid.UUID) (*domain.Selecao, error) {
	if len(idsDosPresentes) == 0 {
		return nil, errors.New("a lista de presentes não pode estar vazia")
	}

	selecao, err := s.selecaoRepo.SalvarSelecao(ctx, chaveDeAcesso, idsDosPresentes)
	if err != nil {
		return nil, fmt.Errorf("falha no serviço ao finalizar seleção: %w", err)
	}
	return selecao, nil
}
