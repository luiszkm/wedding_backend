// file: internal/pagetemplate/application/service.go
package application

import (
	"context"
	"fmt"
	"time"

	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
	galleryDomain "github.com/luiszkm/wedding_backend/internal/gallery/domain"
	giftDomain "github.com/luiszkm/wedding_backend/internal/gift/domain"
	guestDomain "github.com/luiszkm/wedding_backend/internal/guest/domain"
	mbDomain "github.com/luiszkm/wedding_backend/internal/messageboard/domain"
	"github.com/luiszkm/wedding_backend/internal/pagetemplate/domain"
	"github.com/google/uuid"
)

// PageTemplateService implementa a lógica de aplicação para templates de página
type PageTemplateService struct {
	templateEngine   domain.TemplateEngine
	eventRepo        eventDomain.EventoRepository
	guestRepo        guestDomain.GroupRepository
	giftRepo         giftDomain.PresenteRepository
	messageRepo      mbDomain.RecadoRepository
	galleryRepo      galleryDomain.FotoRepository
}

// NewPageTemplateService cria uma nova instância do serviço
func NewPageTemplateService(
	templateEngine domain.TemplateEngine,
	eventRepo eventDomain.EventoRepository,
	guestRepo guestDomain.GroupRepository,
	giftRepo giftDomain.PresenteRepository,
	messageRepo mbDomain.RecadoRepository,
	galleryRepo galleryDomain.FotoRepository,
) domain.TemplateService {
	return &PageTemplateService{
		templateEngine: templateEngine,
		eventRepo:      eventRepo,
		guestRepo:      guestRepo,
		giftRepo:       giftRepo,
		messageRepo:    messageRepo,
		galleryRepo:    galleryRepo,
	}
}

// RenderPublicPage renderiza a página pública de um evento
func (s *PageTemplateService) RenderPublicPage(urlSlug string) ([]byte, error) {
	ctx := context.Background()

	// Buscar evento pelo slug
	evento, err := s.eventRepo.FindBySlug(ctx, urlSlug)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar evento: %w", err)
	}

	// Obter dados completos da página
	pageData, err := s.GetEventPageData(evento)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter dados da página: %w", err)
	}

	// Determinar qual template usar
	templateID := evento.GetTemplateAtivo()

	// Renderizar página
	html, err := s.templateEngine.RenderEventPage(templateID, pageData)
	if err != nil {
		return nil, fmt.Errorf("erro ao renderizar página: %w", err)
	}

	return html, nil
}

// GetEventPageData coleta todos os dados necessários para renderizar uma página
func (s *PageTemplateService) GetEventPageData(evento *eventDomain.Evento) (*domain.EventPageData, error) {
	// ctx := context.Background() // TODO: usar quando implementar repositórios
	
	// Criar estrutura base de dados
	pageData := domain.NewEventPageData(evento)

	// Buscar grupos de convidados - por enquanto vazio, implementar depois
	// TODO: Implementar método FindByEventID no guestRepo
	// guestGroups, err := s.guestRepo.FindByEventID(ctx, evento.ID())

	// TODO: Implementar métodos nos repositórios
	// Por enquanto, dados ficam vazios - funcionalidade será implementada gradualmente
	
	// gifts, err := s.giftRepo.FindByEventID(ctx, evento.ID())
	// messages, err := s.messageRepo.FindApprovedByEventID(ctx, evento.ID())  
	// photos, err := s.galleryRepo.FindByEventID(ctx, evento.ID())

	// Configurar flags de exibição - por enquanto, todas false
	// Quando implementarmos os repositórios, isso será baseado nos dados reais
	pageData.ShowGifts = false
	pageData.ShowGallery = false
	pageData.ShowMessages = false
	pageData.ShowRSVP = false

	// Adicionar informações de contato (pode ser configurado via admin)
	// Por enquanto, usar dados básicos do evento
	pageData.SetContact(evento.Nome(), "", "")

	return pageData, nil
}

// GetTemplateMetadata retorna metadados de um template específico
func (s *PageTemplateService) GetTemplateMetadata(templateID string) (*domain.TemplateMetadata, error) {
	templates, err := s.templateEngine.GetAvailableTemplates()
	if err != nil {
		return nil, err
	}

	for _, tmpl := range templates {
		if tmpl.ID == templateID {
			return tmpl, nil
		}
	}

	return nil, domain.ErrTemplateNaoEncontrado
}

// ListAvailableTemplates lista todos os templates disponíveis
func (s *PageTemplateService) ListAvailableTemplates() ([]*domain.TemplateMetadata, error) {
	return s.templateEngine.GetAvailableTemplates()
}

// UpdateEventTemplate atualiza o template de um evento
func (s *PageTemplateService) UpdateEventTemplate(eventID uuid.UUID, userID uuid.UUID, templateConfig domain.TemplateConfig) error {
	ctx := context.Background()

	// Buscar evento
	evento, err := s.eventRepo.FindByID(ctx, userID, eventID)
	if err != nil {
		return fmt.Errorf("erro ao buscar evento: %w", err)
	}

	// Aplicar configuração do template
	if templateConfig.IsBespoke && templateConfig.BespokeFileName != "" {
		// Validar template bespoke
		if err := s.templateEngine.ValidateTemplate(templateConfig.BespokeFileName); err != nil {
			return fmt.Errorf("template bespoke inválido: %w", err)
		}
		
		if err := evento.DefinirTemplateBespoke(templateConfig.BespokeFileName); err != nil {
			return fmt.Errorf("erro ao definir template bespoke: %w", err)
		}
	} else if templateConfig.StandardTemplateID != "" {
		// Usar template padrão
		if err := evento.DefinirTemplate(templateConfig.StandardTemplateID); err != nil {
			return fmt.Errorf("erro ao definir template padrão: %w", err)
		}
	}

	// Atualizar paleta de cores se fornecida
	if templateConfig.PaletaCores != nil {
		if err := evento.DefinirPaletaCores(templateConfig.PaletaCores); err != nil {
			return fmt.Errorf("erro ao definir paleta de cores: %w", err)
		}
	}

	// Salvar alterações
	if err := s.eventRepo.Update(ctx, evento); err != nil {
		return fmt.Errorf("erro ao salvar alterações do evento: %w", err)
	}

	return nil
}

// ValidateTemplateConfig valida uma configuração de template
func (s *PageTemplateService) ValidateTemplateConfig(config domain.TemplateConfig) error {
	if config.IsBespoke {
		if config.BespokeFileName == "" {
			return fmt.Errorf("nome do arquivo bespoke é obrigatório")
		}
		
		// Validar se template bespoke existe
		if err := s.templateEngine.ValidateTemplate(config.BespokeFileName); err != nil {
			return fmt.Errorf("template bespoke não encontrado ou inválido: %w", err)
		}
	} else {
		if config.StandardTemplateID == "" {
			return fmt.Errorf("ID do template padrão é obrigatório")
		}
		
		// Validar se template padrão existe
		templates := domain.GetStandardTemplates()
		found := false
		for _, tmpl := range templates {
			if tmpl.ID == config.StandardTemplateID {
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("template padrão não encontrado: %s", config.StandardTemplateID)
		}
	}

	// Validar paleta de cores se fornecida
	if config.PaletaCores != nil {
		coresObrigatorias := []string{"primary", "secondary", "background", "text"}
		for _, cor := range coresObrigatorias {
			if _, existe := config.PaletaCores[cor]; !existe {
				return fmt.Errorf("cor obrigatória não encontrada: %s", cor)
			}
		}
	}

	return nil
}

// PreviewTemplate gera uma prévia de como ficará o template (usando dados de exemplo)
func (s *PageTemplateService) PreviewTemplate(templateID string, config domain.TemplateConfig) ([]byte, error) {
	// Criar dados de exemplo para prévia
	exampleEvent := createExampleEvent(config.PaletaCores)
	pageData := createExamplePageData(exampleEvent)

	// Determinar template a usar
	finalTemplateID := templateID
	if config.IsBespoke && config.BespokeFileName != "" {
		finalTemplateID = config.BespokeFileName
	}

	// Renderizar prévia
	return s.templateEngine.RenderEventPage(finalTemplateID, pageData)
}

// Funções auxiliares para criar dados de exemplo

func createExampleEvent(paleta eventDomain.PaletaCores) *eventDomain.Evento {
	// Criar evento de exemplo
	exemplo, _ := eventDomain.NewEvento(
		uuid.New(),
		"João & Maria - Exemplo",
		time.Date(2024, 6, 15, 15, 0, 0, 0, time.Local),
		eventDomain.TipoCasamento,
		"joao-maria-exemplo",
	)
	
	if paleta != nil {
		exemplo.DefinirPaletaCores(paleta)
	}
	
	return exemplo
}

func createExamplePageData(evento *eventDomain.Evento) *domain.EventPageData {
	pageData := domain.NewEventPageData(evento)
	
	// Adicionar dados de exemplo
	pageData.ShowGifts = true
	pageData.ShowGallery = true
	pageData.ShowMessages = true
	pageData.ShowRSVP = true
	
	pageData.SetContact("João & Maria", "contato@casamento.com", "(11) 99999-9999")
	pageData.SetCustomData("exemplo", true)
	
	return pageData
}