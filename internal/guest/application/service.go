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
func (s *GuestService) ObterGrupoPorChaveDeAcesso(ctx context.Context, eventID uuid.UUID, accessKey string) (*domain.GrupoDeConvidados, error) {
	grupo, err := s.repo.FindByAccessKey(ctx, eventID, accessKey)
	if err != nil {
		// Apenas repassa o erro (seja ele "não encontrado" ou um erro técnico).
		return nil, fmt.Errorf("falha ao obter grupo: %w", err)
	}
	return grupo, nil
}

func (s *GuestService) ConfirmarPresencaGrupo(ctx context.Context, eventID uuid.UUID, chaveDeAcesso string, respostas []domain.RespostaRSVP) error {
	// 1. Carregar o agregado pela chave de acesso.
	grupo, err := s.repo.FindByAccessKey(ctx, eventID, chaveDeAcesso)
	if err != nil {
		return fmt.Errorf("falha ao buscar grupo por chave: %w", err)
	}

	// 2. Executar a lógica de negócio no domínio.
	if err := grupo.ConfirmarPresenca(respostas); err != nil {
		return err // Retorna erros de negócio (status inválido, convidado não pertence, etc.)
	}

	// 3. Persistir o agregado inteiro com seu novo estado.
	if err := s.repo.UpdateRSVP(ctx, grupo); err != nil {
		return fmt.Errorf("falha ao salvar confirmação de presença: %w", err)
	}

	return nil
}
func (s *GuestService) RevisarGrupo(ctx context.Context, userID, groupID uuid.UUID, chaveDeAcesso string, convidadosParaRevisao []domain.ConvidadoParaRevisao) error {
	// 1. Carrega o agregado, já com a verificação de propriedade no repositório.
	grupo, err := s.repo.FindByID(ctx, userID, groupID)
	if err != nil {
		return fmt.Errorf("falha ao buscar grupo para revisão: %w", err)
	}

	// 2. Executa a lógica de negócio no domínio
	if err := grupo.Revisar(chaveDeAcesso, convidadosParaRevisao); err != nil {
		return err
	}

	// 3. Persiste as alterações, também com verificação de propriedade.
	if err := s.repo.Update(ctx, userID, grupo); err != nil {
		return fmt.Errorf("falha ao salvar revisão do grupo: %w", err)
	}
	return nil
}

// ListarGruposPorEvento retorna todos os grupos de um evento
func (s *GuestService) ListarGruposPorEvento(ctx context.Context, userID, eventID uuid.UUID, statusFilter string) ([]*domain.GrupoDeConvidados, error) {
	grupos, err := s.repo.FindAllByEventID(ctx, userID, eventID, statusFilter)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar grupos por evento: %w", err)
	}
	return grupos, nil
}

// ObterGrupoPorID retorna um grupo específico (admin)
func (s *GuestService) ObterGrupoPorID(ctx context.Context, userID, groupID uuid.UUID) (*domain.GrupoDeConvidados, error) {
	grupo, err := s.repo.FindByID(ctx, userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter grupo por ID: %w", err)
	}
	return grupo, nil
}

// RemoverGrupo remove um grupo completamente
func (s *GuestService) RemoverGrupo(ctx context.Context, userID, groupID uuid.UUID) error {
	// 1. Carregar o grupo para validação
	grupo, err := s.repo.FindByID(ctx, userID, groupID)
	if err != nil {
		return fmt.Errorf("falha ao buscar grupo para remoção: %w", err)
	}

	// 2. Validar se pode ser removido (regra de negócio)
	if err := grupo.PodeSerRemovido(); err != nil {
		return err
	}

	// 3. Remover o grupo
	if err := s.repo.Delete(ctx, userID, groupID); err != nil {
		return fmt.Errorf("falha ao remover grupo: %w", err)
	}

	return nil
}

// ObterEstatisticasRSVP retorna estatísticas de RSVP para um evento
func (s *GuestService) ObterEstatisticasRSVP(ctx context.Context, userID, eventID uuid.UUID) (*domain.RSVPStats, error) {
	stats, err := s.repo.GetRSVPStats(ctx, userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter estatísticas RSVP: %w", err)
	}
	return stats, nil
}
