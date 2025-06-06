// file: internal/guest/application/service.go
package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/luiszkm/wedding_backend/internal/guest/domain"
)

type GuestService struct {
	repo domain.GroupRepository
}

func NewGuestService(repo domain.GroupRepository) *GuestService {
	return &GuestService{repo: repo}
}

// CriarNovoGrupo é um caso de uso da aplicação.
func (s *GuestService) CriarNovoGrupo(ctx context.Context, idCasamento uuid.UUID, chaveDeAcesso string, nomesDosConvidados []string) (uuid.UUID, error) {
	// 1. Usa a fábrica do domínio para criar o agregado. A lógica de negócio está protegida.
	novoGrupo, err := domain.NewGrupoDeConvidados(idCasamento, chaveDeAcesso, nomesDosConvidados)
	if err != nil {
		return uuid.Nil, fmt.Errorf("falha ao criar novo grupo de convidados: %w", err)
	}

	// 2. Usa o repositório para persistir o novo agregado.
	if err := s.repo.Save(ctx, novoGrupo); err != nil {
		return uuid.Nil, fmt.Errorf("falha ao salvar novo grupo de convidados: %w", err)
	}

	// 3. Retorna o resultado.
	return novoGrupo.ID(), nil
}

// ObterGrupoPorChaveDeAcesso é o caso de uso para a busca.
func (s *GuestService) ObterGrupoPorChaveDeAcesso(ctx context.Context, accessKey string) (*domain.GrupoDeConvidados, error) {
	grupo, err := s.repo.FindByAccessKey(ctx, accessKey)
	if err != nil {
		// Apenas repassa o erro (seja ele "não encontrado" ou um erro técnico).
		return nil, fmt.Errorf("falha ao obter grupo: %w", err)
	}
	return grupo, nil
}
