package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	"github.com/luiszkm/wedding_backend/internal/gift/domain"
)

type ItemSelecao struct {
	IDPresente uuid.UUID
	Quantidade int
}

type GiftService struct {
	repo        domain.PresenteRepository
	selecaoRepo domain.SelecaoRepository
	eventRepo   eventDomain.EventoRepository
}

func NewGiftService(presenteRepo domain.PresenteRepository, selecaoRepo domain.SelecaoRepository, eventRepo eventDomain.EventoRepository) *GiftService {
	return &GiftService{repo: presenteRepo, selecaoRepo: selecaoRepo, eventRepo: eventRepo}
}

func (s *GiftService) CriarPresenteIntegral(ctx context.Context, userID, idEvento uuid.UUID, nome, desc, fotoURL, categoria string, favorito bool, detalhes domain.DetalhesPresente) (*domain.Presente, error) {
	// Verificar permissão
	_, err := s.eventRepo.FindByID(ctx, userID, idEvento)
	if err != nil {
		return nil, fmt.Errorf("permissão negada ou evento não encontrado: %w", err)
	}

	// Criar presente integral
	novoPresente, err := domain.NewPresenteIntegral(idEvento, nome, desc, fotoURL, favorito, categoria, detalhes)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, novoPresente); err != nil {
		return nil, fmt.Errorf("falha ao salvar novo presente: %w", err)
	}

	return novoPresente, nil
}

func (s *GiftService) CriarPresenteFracionado(ctx context.Context, userID, idEvento uuid.UUID, nome, desc, fotoURL, categoria string, favorito bool, detalhes domain.DetalhesPresente, valorTotal float64, numeroCotas int) (*domain.Presente, error) {
	// Verificar permissão
	_, err := s.eventRepo.FindByID(ctx, userID, idEvento)
	if err != nil {
		return nil, fmt.Errorf("permissão negada ou evento não encontrado: %w", err)
	}

	// Criar presente fracionado
	novoPresente, err := domain.NewPresenteFracionado(idEvento, nome, desc, fotoURL, favorito, categoria, detalhes, valorTotal, numeroCotas)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveWithCotas(ctx, novoPresente); err != nil {
		return nil, fmt.Errorf("falha ao salvar novo presente fracionado: %w", err)
	}

	return novoPresente, nil
}

// Método legacy mantido para compatibilidade (será removido depois de atualizar handlers)
func (s *GiftService) CriarNovoPresente(ctx context.Context, userID, idEvento uuid.UUID, nome, desc, fotoURL, categoria string, favorito bool, detalhes domain.DetalhesPresente) (*domain.Presente, error) {
	return s.CriarPresenteIntegral(ctx, userID, idEvento, nome, desc, fotoURL, categoria, favorito, detalhes)
}

func (s *GiftService) ListarPresentesDisponiveis(ctx context.Context, casamentoID uuid.UUID) ([]*domain.Presente, error) {
	presentes, err := s.repo.ListarDisponiveisPorCasamento(ctx, casamentoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar lista de presentes: %w", err)
	}
	return presentes, nil
}

func (s *GiftService) ListarTodosPresentesPorEvento(ctx context.Context, userID, eventoID uuid.UUID) ([]*domain.PresenteComSelecao, error) {
	// Verificar permissão: evento deve pertencer ao usuário
	_, err := s.eventRepo.FindByID(ctx, userID, eventoID)
	if err != nil {
		return nil, fmt.Errorf("permissão negada ou evento não encontrado: %w", err)
	}

	presentesComSelecao, err := s.repo.ListarTodosPorEvento(ctx, eventoID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar todos os presentes do evento: %w", err)
	}

	return presentesComSelecao, nil
}

func (s *GiftService) FinalizarSelecaoDePresentes(ctx context.Context, chaveDeAcesso string, itens []ItemSelecao) (*domain.Selecao, error) {
	if len(itens) == 0 {
		return nil, errors.New("a lista de presentes não pode estar vazia")
	}

	// Extrair IDs únicos dos presentes
	presenteIDs := make([]uuid.UUID, 0, len(itens))
	itensMap := make(map[uuid.UUID]int)

	for _, item := range itens {
		if item.Quantidade <= 0 {
			return nil, fmt.Errorf("quantidade deve ser positiva para presente %s", item.IDPresente.String())
		}
		presenteIDs = append(presenteIDs, item.IDPresente)
		itensMap[item.IDPresente] = item.Quantidade
	}

	// Buscar presentes
	presentes, err := s.repo.FindByIDs(ctx, presenteIDs)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar presentes: %w", err)
	}

	if len(presentes) != len(presenteIDs) {
		return nil, errors.New("um ou mais presentes não foram encontrados")
	}

	// Validar disponibilidade e preparar seleção
	presentesConfirmados := make([]domain.PresenteConfirmado, 0, len(presentes))
	conflitantes := make([]uuid.UUID, 0)

	for _, presente := range presentes {
		quantidade := itensMap[presente.ID()]

		if presente.EhIntegral() {
			// Presente integral: quantidade deve ser 1 e deve estar disponível
			if quantidade != 1 {
				return nil, fmt.Errorf("presente integral %s deve ter quantidade 1", presente.Nome())
			}
			if presente.Status() != domain.StatusDisponivel {
				conflitantes = append(conflitantes, presente.ID())
				continue
			}

			presentesConfirmados = append(presentesConfirmados, domain.PresenteConfirmado{
				ID:         presente.ID(),
				Nome:       presente.Nome(),
				Quantidade: 1,
				ValorCota:  nil,
			})
		} else {
			// Presente fracionado: verificar cotas disponíveis
			cotasDisponiveis := presente.ContarCotasDisponiveis()
			if quantidade > cotasDisponiveis {
				return nil, fmt.Errorf("presente %s tem apenas %d cotas disponíveis, solicitado %d", presente.Nome(), cotasDisponiveis, quantidade)
			}

			valorCota := presente.ObterValorCota()
			presentesConfirmados = append(presentesConfirmados, domain.PresenteConfirmado{
				ID:         presente.ID(),
				Nome:       presente.Nome(),
				Quantidade: quantidade,
				ValorCota:  &valorCota,
			})
		}
	}

	if len(conflitantes) > 0 {
		return nil, &domain.ErrPresentesConflitantes{PresentesIDs: conflitantes}
	}

	// Finalizar seleção passando o mapa de quantidades
	selecao, err := s.selecaoRepo.SalvarSelecao(ctx, chaveDeAcesso, itensMap)
	if err != nil {
		return nil, fmt.Errorf("falha no serviço ao finalizar seleção: %w", err)
	}

	return selecao, nil
}

// Método legacy mantido para compatibilidade
func (s *GiftService) FinalizarSelecaoDepresentes(ctx context.Context, chaveDeAcesso string, idsDosPresentes []uuid.UUID) (*domain.Selecao, error) {
	// Converter para novo formato (todos com quantidade 1)
	itens := make([]ItemSelecao, len(idsDosPresentes))
	for i, id := range idsDosPresentes {
		itens[i] = ItemSelecao{
			IDPresente: id,
			Quantidade: 1,
		}
	}

	return s.FinalizarSelecaoDePresentes(ctx, chaveDeAcesso, itens)
}

func (s *GiftService) DeletarPresente(ctx context.Context, userID, presenteID uuid.UUID) error {
	// Buscar o presente para verificar permissões e validar regras de negócio
	presente, err := s.repo.FindByID(ctx, userID, presenteID)
	if err != nil {
		return err
	}

	// Validar se o presente pode ser deletado
	// Não permitir deletar presentes que já foram selecionados
	if presente.Status() == domain.StatusSelecionado {
		return errors.New("não é possível deletar um presente que já foi selecionado")
	}

	// Permitir deletar presentes parcialmente selecionados com aviso no log
	// (pode ser mudado para bloquear se necessário)
	if presente.Status() == domain.StatusParcialmenteSelecionado {
		// Log de aviso - presente parcialmente selecionado será deletado
		// As cotas serão removidas mas as seleções ficarão órfãs
	}

	return s.repo.Delete(ctx, userID, presenteID)
}

func (s *GiftService) AtualizarPresente(ctx context.Context, userID, presenteID uuid.UUID, nome, descricao, categoria, fotoURL string, ehFavorito bool, detalhes domain.DetalhesPresente) error {
	// Buscar o presente para verificar permissões
	presente, err := s.repo.FindByID(ctx, userID, presenteID)
	if err != nil {
		return err
	}

	// Atualizar dados no domínio
	if err := presente.AtualizarDados(nome, descricao, categoria, ehFavorito, detalhes); err != nil {
		return err
	}

	// Atualizar foto se fornecida
	if fotoURL != "" {
		presente.AtualizarFoto(fotoURL)
	}

	// Persistir mudanças
	return s.repo.Update(ctx, presente)
}
