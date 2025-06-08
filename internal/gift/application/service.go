// file: internal/gift/application/service.go
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/gift/domain"
)

type GiftService struct {
	repo        domain.PresenteRepository
	selecaoRepo domain.SelecaoRepository
}

func NewGiftService(repo domain.PresenteRepository, selecaoRepo domain.SelecaoRepository) *GiftService {
	return &GiftService{repo: repo, selecaoRepo: selecaoRepo}
}

func (s *GiftService) CriarNovoPresente(ctx context.Context, idCasamento uuid.UUID, nome, desc,
	fotoURL string, favorito bool, categoria string, detalhes domain.DetalhesPresente) (*domain.Presente, error) {
	novoPresente, err := domain.NewPresente(idCasamento, nome, desc, fotoURL, favorito, categoria, detalhes)
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
