// file: internal/pagetemplate/domain/template_test.go
package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	eventDomain "github.com/luiszkm/wedding_backend/internal/event/domain"
)

func TestNewTemplateMetadata(t *testing.T) {
	t.Run("deve criar template metadata válido", func(t *testing.T) {
		// Arrange
		id := "template_test"
		nome := "Template Teste"
		descricao := "Template para testes"
		tipo := TipoStandard
		caminho := "template_test.html"

		// Act
		metadata := NewTemplateMetadata(id, nome, descricao, tipo, caminho)

		// Assert
		assert.Equal(t, id, metadata.ID)
		assert.Equal(t, nome, metadata.Nome)
		assert.Equal(t, descricao, metadata.Descricao)
		assert.Equal(t, tipo, metadata.Tipo)
		assert.Equal(t, caminho, metadata.CaminhoArquivo)
		assert.True(t, metadata.SuportaGifts)
		assert.True(t, metadata.SuportaGallery)
		assert.True(t, metadata.SuportaMessages)
		assert.True(t, metadata.SuportaRSVP)
		assert.NotNil(t, metadata.PaletaDefault)
		assert.NotZero(t, metadata.CriadoEm)
	})
}

func TestTemplateMetadata_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		metadata *TemplateMetadata
		wantErr  bool
	}{
		{
			name: "metadata válido",
			metadata: &TemplateMetadata{
				ID:             "template_test",
				Nome:           "Template Teste",
				CaminhoArquivo: "template_test.html",
				Tipo:           TipoStandard,
			},
			wantErr: false,
		},
		{
			name: "ID vazio",
			metadata: &TemplateMetadata{
				Nome:           "Template Teste",
				CaminhoArquivo: "template_test.html",
				Tipo:           TipoStandard,
			},
			wantErr: true,
		},
		{
			name: "Nome vazio",
			metadata: &TemplateMetadata{
				ID:             "template_test",
				CaminhoArquivo: "template_test.html",
				Tipo:           TipoStandard,
			},
			wantErr: true,
		},
		{
			name: "Caminho arquivo vazio",
			metadata: &TemplateMetadata{
				ID:   "template_test",
				Nome: "Template Teste",
				Tipo: TipoStandard,
			},
			wantErr: true,
		},
		{
			name: "Tipo inválido",
			metadata: &TemplateMetadata{
				ID:             "template_test",
				Nome:           "Template Teste",
				CaminhoArquivo: "template_test.html",
				Tipo:           TipoTemplate("INVALIDO"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metadata.IsValid()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplateMetadata_GetFullPath(t *testing.T) {
	tests := []struct {
		name     string
		metadata *TemplateMetadata
		want     string
	}{
		{
			name: "template padrão",
			metadata: &TemplateMetadata{
				Tipo:           TipoStandard,
				CaminhoArquivo: "template_moderno.html",
			},
			want: "templates/standard/template_moderno.html",
		},
		{
			name: "template bespoke",
			metadata: &TemplateMetadata{
				Tipo:           TipoBespoke,
				CaminhoArquivo: "cliente_xyz.html",
			},
			want: "templates/bespoke/cliente_xyz.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.metadata.GetFullPath()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewEventPageData(t *testing.T) {
	t.Run("deve criar EventPageData válido", func(t *testing.T) {
		// Arrange
		evento, _ := eventDomain.NewEvento(
			uuid.New(),
			"Casamento Teste",
			time.Now(),
			eventDomain.TipoCasamento,
			"casamento-teste",
		)

		// Act
		pageData := NewEventPageData(evento)

		// Assert
		assert.Equal(t, evento, pageData.Event)
		assert.Equal(t, evento.PaletaCores(), pageData.PaletaCores)
		assert.Empty(t, pageData.GuestGroups)
		assert.Empty(t, pageData.Gifts)
		assert.Empty(t, pageData.Messages)
		assert.Empty(t, pageData.Photos)
		assert.True(t, pageData.ShowGifts)
		assert.True(t, pageData.ShowGallery)
		assert.True(t, pageData.ShowMessages)
		assert.True(t, pageData.ShowRSVP)
		assert.NotNil(t, pageData.CustomData)
	})
}

func TestEventPageData_SetContact(t *testing.T) {
	t.Run("deve definir informações de contato", func(t *testing.T) {
		// Arrange
		evento, _ := eventDomain.NewEvento(
			uuid.New(),
			"Casamento Teste",
			time.Now(),
			eventDomain.TipoCasamento,
			"casamento-teste",
		)
		pageData := NewEventPageData(evento)

		// Act
		pageData.SetContact("João Silva", "joao@teste.com", "11999999999")

		// Assert
		assert.NotNil(t, pageData.Contact)
		assert.Equal(t, "João Silva", pageData.Contact.Nome)
		assert.Equal(t, "joao@teste.com", pageData.Contact.Email)
		assert.Equal(t, "11999999999", pageData.Contact.Telefone)
	})
}

func TestEventPageData_SetCustomData(t *testing.T) {
	t.Run("deve definir dados customizados", func(t *testing.T) {
		// Arrange
		evento, _ := eventDomain.NewEvento(
			uuid.New(),
			"Casamento Teste",
			time.Now(),
			eventDomain.TipoCasamento,
			"casamento-teste",
		)
		pageData := NewEventPageData(evento)

		// Act
		pageData.SetCustomData("chave_teste", "valor_teste")

		// Assert
		assert.Equal(t, "valor_teste", pageData.GetCustomData("chave_teste"))
		assert.Nil(t, pageData.GetCustomData("chave_inexistente"))
	})
}

func TestEventPageData_Validate(t *testing.T) {
	tests := []struct {
		name     string
		pageData *EventPageData
		wantErr  bool
	}{
		{
			name: "dados válidos",
			pageData: func() *EventPageData {
				evento, _ := eventDomain.NewEvento(
					uuid.New(),
					"Evento Teste",
					time.Now(),
					eventDomain.TipoCasamento,
					"evento-teste",
				)
				return &EventPageData{
					Event: evento,
					PaletaCores: eventDomain.PaletaCores{
						"primary": "#000000",
					},
				}
			}(),
			wantErr: false,
		},
		{
			name: "evento nil",
			pageData: &EventPageData{
				PaletaCores: eventDomain.PaletaCores{
					"primary": "#000000",
				},
			},
			wantErr: true,
		},
		{
			name: "paleta cores nil",
			pageData: &EventPageData{
				Event: &eventDomain.Evento{},
			},
			wantErr: true,
		},
		{
			name: "paleta cores vazia",
			pageData: &EventPageData{
				Event:       &eventDomain.Evento{},
				PaletaCores: eventDomain.PaletaCores{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pageData.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetStandardTemplates(t *testing.T) {
	t.Run("deve retornar templates padrão", func(t *testing.T) {
		// Act
		templates := GetStandardTemplates()

		// Assert
		assert.NotEmpty(t, templates)
		assert.Len(t, templates, 3) // moderno, clássico, elegante

		// Verificar que todos têm os dados necessários
		for _, tmpl := range templates {
			assert.NotEmpty(t, tmpl.ID)
			assert.NotEmpty(t, tmpl.Nome)
			assert.NotEmpty(t, tmpl.Descricao)
			assert.Equal(t, TipoStandard, tmpl.Tipo)
			assert.NotEmpty(t, tmpl.CaminhoArquivo)
			assert.NotNil(t, tmpl.PaletaDefault)
			assert.NoError(t, tmpl.IsValid())
		}

		// Verificar IDs específicos
		templateIDs := make([]string, len(templates))
		for i, tmpl := range templates {
			templateIDs[i] = tmpl.ID
		}
		assert.Contains(t, templateIDs, "template_moderno")
		assert.Contains(t, templateIDs, "template_classico")
		assert.Contains(t, templateIDs, "template_elegante")
	})
}