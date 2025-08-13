# 🎨 Documentação de Integração Frontend - Módulo de Templates

## Visão Geral

O módulo de Templates permite seleção e customização de templates para páginas de eventos. Inclui templates padrão, personalização de cores e preview em tempo real.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **Listar Templates Disponíveis**

**Endpoint:** `GET /templates/disponiveis`

**Resposta:**
```json
{
  "templates": [
    {
      "id": "template_moderno",
      "nome": "Moderno",
      "descricao": "Template moderno e minimalista com design clean",
      "tipo": "STANDARD",
      "paleta_default": {
        "primary": "#2563eb",
        "secondary": "#f1f5f9",
        "accent": "#10b981",
        "background": "#ffffff",
        "text": "#1f2937"
      },
      "suporta_gifts": true,
      "suporta_gallery": true,
      "suporta_messages": true,
      "suporta_rsvp": true,
      "preview_url": "https://storage.example.com/templates/moderno_preview.jpg"
    }
  ]
}
```

### 2. 🔍 **Obter Template Específico**

**Endpoint:** `GET /templates/{templateId}`

**Resposta:**
```json
{
  "id": "template_moderno",
  "nome": "Moderno",
  "descricao": "Template moderno e minimalista",
  "configuracoes_disponiveis": {
    "paletas_cores": [
      {
        "nome": "Azul Moderno",
        "cores": {
          "primary": "#2563eb",
          "secondary": "#f1f5f9"
        }
      }
    ],
    "fontes": ["Inter", "Roboto", "Open Sans"],
    "layouts": ["single-page", "multi-section"]
  }
}
```

## 🎨 Componentes React

### Hook para Templates

```javascript
// hooks/useTemplates.js
import { useState, useEffect } from 'react';

export const useTemplates = () => {
  const [templates, setTemplates] = useState([]);
  const [selectedTemplate, setSelectedTemplate] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchAvailableTemplates = async () => {
    setLoading(true);
    try {
      const response = await fetch('/v1/templates/disponiveis');
      const data = await response.json();
      setTemplates(data.templates);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchTemplateDetails = async (templateId) => {
    try {
      const response = await fetch(`/v1/templates/${templateId}`);
      const data = await response.json();
      setSelectedTemplate(data);
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  useEffect(() => {
    fetchAvailableTemplates();
  }, []);

  return {
    templates,
    selectedTemplate,
    loading,
    error,
    fetchTemplateDetails,
    setSelectedTemplate
  };
};
```

### Componente Seletor de Templates

```javascript
// components/TemplateSelector.jsx
import React, { useState } from 'react';
import { useTemplates } from '../hooks/useTemplates';

const TemplateSelector = ({ onTemplateSelect, selectedTemplateId }) => {
  const { templates, loading, error, fetchTemplateDetails } = useTemplates();
  const [selectedTemplate, setSelectedTemplate] = useState(null);
  const [customizations, setCustomizations] = useState({});

  const handleTemplateClick = async (template) => {
    try {
      const details = await fetchTemplateDetails(template.id);
      setSelectedTemplate(details);
      setCustomizations({
        paleta_cores: details.paleta_default,
        fonte_principal: details.configuracoes_disponiveis?.fontes?.[0] || 'Inter'
      });
    } catch (error) {
      console.error('Erro ao carregar template:', error);
    }
  };

  const handleCustomizationChange = (key, value) => {
    const newCustomizations = { ...customizations, [key]: value };
    setCustomizations(newCustomizations);
    
    if (selectedTemplate) {
      onTemplateSelect?.({
        template_id: selectedTemplate.id,
        configuracoes_template: newCustomizations
      });
    }
  };

  if (loading) return <div className="loading">Carregando templates...</div>;
  if (error) return <div className="error">Erro: {error}</div>;

  return (
    <div className="template-selector">
      <h3>Escolha um Template</h3>
      
      <div className="templates-grid">
        {templates.map((template) => (
          <div 
            key={template.id}
            className={`template-card ${
              selectedTemplateId === template.id ? 'selected' : ''
            }`}
            onClick={() => handleTemplateClick(template)}
          >
            {template.preview_url && (
              <img 
                src={template.preview_url} 
                alt={template.nome}
                className="template-preview"
              />
            )}
            <div className="template-info">
              <h4>{template.nome}</h4>
              <p>{template.descricao}</p>
              
              <div className="template-features">
                {template.suporta_gifts && <span className="feature">Presentes</span>}
                {template.suporta_gallery && <span className="feature">Galeria</span>}
                {template.suporta_messages && <span className="feature">Recados</span>}
                {template.suporta_rsvp && <span className="feature">RSVP</span>}
              </div>
            </div>
          </div>
        ))}
      </div>
      
      {selectedTemplate && (
        <div className="customization-panel">
          <h4>Personalizar {selectedTemplate.nome}</h4>
          
          {/* Color Palette Selector */}
          <div className="customization-group">
            <label>Paleta de Cores:</label>
            <div className="color-palette">
              {Object.entries(customizations.paleta_cores || {}).map(([key, color]) => (
                <div key={key} className="color-input-group">
                  <label>{key}:</label>
                  <input
                    type="color"
                    value={color}
                    onChange={(e) => handleCustomizationChange('paleta_cores', {
                      ...customizations.paleta_cores,
                      [key]: e.target.value
                    })}
                  />
                  <span className="color-value">{color}</span>
                </div>
              ))}
            </div>
          </div>
          
          {/* Font Selector */}
          {selectedTemplate.configuracoes_disponiveis?.fontes && (
            <div className="customization-group">
              <label>Fonte Principal:</label>
              <select
                value={customizations.fonte_principal || ''}
                onChange={(e) => handleCustomizationChange('fonte_principal', e.target.value)}
              >
                {selectedTemplate.configuracoes_disponiveis.fontes.map(font => (
                  <option key={font} value={font}>{font}</option>
                ))}
              </select>
            </div>
          )}
          
          <div className="preview-section">
            <h5>Preview</h5>
            <div 
              className="template-preview-box"
              style={{
                backgroundColor: customizations.paleta_cores?.background,
                color: customizations.paleta_cores?.text,
                fontFamily: customizations.fonte_principal
              }}
            >
              <div 
                className="preview-header"
                style={{ backgroundColor: customizations.paleta_cores?.primary, color: 'white' }}
              >
                João & Maria
              </div>
              <div className="preview-content">
                <p>Exemplo de conteúdo do template</p>
                <button 
                  style={{ 
                    backgroundColor: customizations.paleta_cores?.accent,
                    color: 'white',
                    border: 'none',
                    padding: '8px 16px',
                    borderRadius: '4px'
                  }}
                >
                  Botão de Exemplo
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TemplateSelector;
```

## 🎨 Estilos CSS

```css
.template-selector {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.template-card {
  border: 2px solid #eee;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
}

.template-card:hover {
  border-color: #2196F3;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.template-card.selected {
  border-color: #2196F3;
  background: #f0f8ff;
}

.template-preview {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.template-info {
  padding: 15px;
}

.template-info h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.template-info p {
  margin: 0 0 15px 0;
  color: #666;
  font-size: 14px;
}

.template-features {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.feature {
  background: #e3f2fd;
  color: #1976d2;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.customization-panel {
  background: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  margin-top: 20px;
}

.customization-group {
  margin-bottom: 20px;
}

.customization-group label {
  display: block;
  margin-bottom: 10px;
  font-weight: bold;
}

.color-palette {
  display: grid;
  gap: 15px;
}

.color-input-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.color-input-group label {
  min-width: 80px;
  margin: 0;
  font-weight: normal;
  text-transform: capitalize;
}

.color-input-group input[type="color"] {
  width: 40px;
  height: 30px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.color-value {
  font-family: monospace;
  font-size: 12px;
  color: #666;
}

.template-preview-box {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  margin-top: 10px;
}

.preview-header {
  padding: 15px;
  text-align: center;
  font-size: 18px;
  font-weight: bold;
}

.preview-content {
  padding: 20px;
  text-align: center;
}

.preview-content p {
  margin-bottom: 15px;
}

@media (max-width: 768px) {
  .templates-grid {
    grid-template-columns: 1fr;
  }
  
  .color-input-group {
    flex-direction: column;
    align-items: flex-start;
  }
}
```

## 📱 Considerações para UX

### 1. **Preview Visual**
- Imagens de preview dos templates
- Preview em tempo real das customizações
- Comparação lado a lado

### 2. **Customização**
- Seletor de cores intuitivo
- Paletas pré-definidas
- Preview instantâneo

### 3. **Usabilidade**
- Indicação clara de funcionalidades suportadas
- Seleção visual clara
- Reset para configurações padrão

---

Para testes, consulte a documentação da API de templates.