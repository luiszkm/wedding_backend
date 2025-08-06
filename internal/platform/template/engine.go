// file: internal/platform/template/engine.go
package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/luiszkm/wedding_backend/internal/pagetemplate/domain"
)

// GoTemplateEngine implementa o TemplateEngine usando html/template do Go
type GoTemplateEngine struct {
	templatesDir string
	templates    map[string]*template.Template
	mutex        sync.RWMutex
	funcMap      template.FuncMap
}

// NewGoTemplateEngine cria uma nova instância do template engine
func NewGoTemplateEngine(templatesDir string) *GoTemplateEngine {
	engine := &GoTemplateEngine{
		templatesDir: templatesDir,
		templates:    make(map[string]*template.Template),
		funcMap:      createFuncMap(),
	}
	
	// Carregar templates na inicialização
	engine.loadAllTemplates()
	
	return engine
}

// createFuncMap cria funções customizadas para os templates
func createFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatDate": func(format string, t interface{}) string {
			// Implementar formatação de data personalizada se necessário
			return fmt.Sprintf("%v", t)
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil
				}
				dict[key] = values[i+1]
			}
			return dict
		},
	}
}

// RenderEventPage renderiza uma página de evento usando o template especificado
func (gte *GoTemplateEngine) RenderEventPage(templateID string, data *domain.EventPageData) ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("dados inválidos: %w", err)
	}

	tmpl, err := gte.LoadTemplate(templateID)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar template %s: %w", templateID, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("erro ao executar template %s: %w", templateID, err)
	}

	return buf.Bytes(), nil
}

// LoadTemplate carrega um template específico
func (gte *GoTemplateEngine) LoadTemplate(templateID string) (*template.Template, error) {
	gte.mutex.RLock()
	if tmpl, exists := gte.templates[templateID]; exists {
		gte.mutex.RUnlock()
		return tmpl, nil
	}
	gte.mutex.RUnlock()

	// Template não está em cache, carregar do disco
	return gte.loadTemplate(templateID)
}

// loadTemplate carrega um template do sistema de arquivos
func (gte *GoTemplateEngine) loadTemplate(templateID string) (*template.Template, error) {
	gte.mutex.Lock()
	defer gte.mutex.Unlock()

	// Verificar se já foi carregado enquanto esperávamos o lock
	if tmpl, exists := gte.templates[templateID]; exists {
		return tmpl, nil
	}

	templatePath := gte.getTemplatePath(templateID)
	
	// Verificar se o arquivo existe
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template %s não encontrado em %s", templateID, templatePath)
	}

	// Criar novo template com funções personalizadas
	tmpl := template.New(templateID).Funcs(gte.funcMap)

	// Carregar partials primeiro
	partialsPath := filepath.Join(gte.templatesDir, "partials")
	if err := gte.loadPartials(tmpl, partialsPath); err != nil {
		return nil, fmt.Errorf("erro ao carregar partials: %w", err)
	}

	// Carregar template principal
	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do template %s: %w", templateID, err)
	}

	// Armazenar no cache
	gte.templates[templateID] = tmpl

	return tmpl, nil
}

// getTemplatePath retorna o caminho completo para um template
func (gte *GoTemplateEngine) getTemplatePath(templateID string) string {
	// Se terminar com .html, assumir que é um template bespoke
	if strings.HasSuffix(templateID, ".html") {
		return filepath.Join(gte.templatesDir, "bespoke", templateID)
	}
	
	// Caso contrário, é um template padrão
	return filepath.Join(gte.templatesDir, "standard", templateID+".html")
}

// loadPartials carrega todos os partials do diretório
func (gte *GoTemplateEngine) loadPartials(tmpl *template.Template, partialsPath string) error {
	if _, err := os.Stat(partialsPath); os.IsNotExist(err) {
		// Diretório de partials não existe, continuar sem erro
		return nil
	}

	return filepath.WalkDir(partialsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Obter nome relativo para o partial
		relPath, err := filepath.Rel(gte.templatesDir, path)
		if err != nil {
			return err
		}

		// Normalizar path para usar no template
		partialName := strings.ReplaceAll(relPath, "\\", "/")

		_, err = tmpl.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("erro ao fazer parse do partial %s: %w", partialName, err)
		}

		return nil
	})
}

// ValidateTemplate valida se um template é válido
func (gte *GoTemplateEngine) ValidateTemplate(templateID string) error {
	templatePath := gte.getTemplatePath(templateID)
	
	// Verificar se o arquivo existe
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return domain.ErrTemplateNaoEncontrado
	}

	// Tentar fazer parse do template
	tmpl := template.New("validation").Funcs(gte.funcMap)
	
	// Carregar partials
	partialsPath := filepath.Join(gte.templatesDir, "partials")
	if err := gte.loadPartials(tmpl, partialsPath); err != nil {
		return fmt.Errorf("erro nos partials: %w", err)
	}

	// Parse do template principal
	_, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		return domain.ErrTemplateMalformado
	}

	return nil
}

// GetAvailableTemplates retorna todos os templates disponíveis
func (gte *GoTemplateEngine) GetAvailableTemplates() ([]*domain.TemplateMetadata, error) {
	var templates []*domain.TemplateMetadata

	// Adicionar templates padrão
	standardTemplates := domain.GetStandardTemplates()
	for _, tmpl := range standardTemplates {
		// Verificar se o arquivo realmente existe
		if err := gte.ValidateTemplate(tmpl.ID); err == nil {
			templates = append(templates, tmpl)
		}
	}

	// Adicionar templates bespoke
	bespokePath := filepath.Join(gte.templatesDir, "bespoke")
	if _, err := os.Stat(bespokePath); err == nil {
		bespokeTemplates, err := gte.scanBespokeTemplates(bespokePath)
		if err == nil {
			templates = append(templates, bespokeTemplates...)
		}
	}

	return templates, nil
}

// scanBespokeTemplates escaneia o diretório bespoke por templates personalizados
func (gte *GoTemplateEngine) scanBespokeTemplates(bespokePath string) ([]*domain.TemplateMetadata, error) {
	var templates []*domain.TemplateMetadata

	err := filepath.WalkDir(bespokePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		// Criar metadata básico para template bespoke
		fileName := d.Name()
		templateID := strings.TrimSuffix(fileName, ".html")
		
		tmpl := &domain.TemplateMetadata{
			ID:             fileName, // Para bespoke, usamos o nome do arquivo completo
			Nome:           templateID,
			Descricao:      "Template personalizado",
			Tipo:           domain.TipoBespoke,
			CaminhoArquivo: fileName,
			SuportaGifts:   true, // Assumir suporte completo para templates bespoke
			SuportaGallery: true,
			SuportaMessages: true,
			SuportaRSVP:    true,
		}

		templates = append(templates, tmpl)
		return nil
	})

	return templates, err
}

// loadAllTemplates carrega todos os templates disponíveis no cache
func (gte *GoTemplateEngine) loadAllTemplates() {
	// Esta função pode ser executada em background para pre-carregar templates
	// Por enquanto, deixamos o carregamento lazy (sob demanda)
}

// ClearCache limpa o cache de templates
func (gte *GoTemplateEngine) ClearCache() {
	gte.mutex.Lock()
	defer gte.mutex.Unlock()
	gte.templates = make(map[string]*template.Template)
}

// ReloadTemplate força o reload de um template específico
func (gte *GoTemplateEngine) ReloadTemplate(templateID string) error {
	gte.mutex.Lock()
	delete(gte.templates, templateID)
	gte.mutex.Unlock()
	
	_, err := gte.loadTemplate(templateID)
	return err
}