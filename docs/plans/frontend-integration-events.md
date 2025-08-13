# 💒 Documentação de Integração Frontend - Módulo de Eventos

## Visão Geral

O módulo de Eventos gerencia a criação e configuração de eventos de casamento, incluindo páginas públicas personalizadas com templates. Permite definir informações básicas, URL slug personalizada e integração com outros módulos.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 🌐 **Página Pública do Evento**

**Endpoint:** `GET /eventos/{urlSlug}/pagina`

**Descrição:** Renderiza a página pública de um evento usando seu template configurado.

**Parâmetros:**
- `urlSlug`: URL amigável do evento (ex: "casamento-joao-maria-2024")

**Headers de Resposta:**
- Content-Type: `text/html; charset=utf-8`
- Cache-Control: `public, max-age=300`

**Exemplo de uso:**
```javascript
// Redirecionar para página do evento
window.location.href = `/v1/eventos/${urlSlug}/pagina`;

// Ou carregar em iframe
const iframe = document.createElement('iframe');
iframe.src = `/v1/eventos/${urlSlug}/pagina`;
```

### 2. ➕ **Criar Evento (Autenticado)**

**Endpoint:** `POST /eventos`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Body da Requisição:**
```json
{
  "nome": "Casamento João e Maria",
  "descricao": "Celebração do nosso casamento",
  "data_evento": "2024-06-15T15:00:00Z",
  "local": "Igreja São Pedro, São Paulo",
  "url_slug": "casamento-joao-maria-2024",
  "template_id": "template_moderno",
  "configuracoes_template": {
    "paleta_cores": {
      "primary": "#2563eb",
      "secondary": "#f1f5f9",
      "accent": "#10b981"
    },
    "fonte_principal": "Inter",
    "mostrar_contador": true
  }
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "url_publica": "https://api.exemplo.com/v1/eventos/casamento-joao-maria-2024/pagina"
}
```

### 3. ✏️ **Atualizar Evento (Autenticado)**

**Endpoint:** `PUT /eventos/{idEvento}`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

### 4. 📋 **Listar Eventos do Usuário (Autenticado)**

**Endpoint:** `GET /usuarios/{idUsuario}/eventos`

**Resposta:**
```json
{
  "eventos": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "nome": "Casamento João e Maria",
      "data_evento": "2024-06-15T15:00:00Z",
      "url_slug": "casamento-joao-maria-2024",
      "status": "ativo",
      "template_id": "template_moderno",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## 🎨 Componentes React

### Hook para Gestão de Eventos

```javascript
// hooks/useEvents.js
import { useState } from 'react';

export const useEvents = (token) => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const createEvent = async (eventData) => {
    if (!token) throw new Error('Token necessário');
    
    setLoading(true);
    try {
      const response = await fetch('/v1/eventos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(eventData)
      });
      
      if (!response.ok) throw new Error('Erro ao criar evento');
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const fetchUserEvents = async (userId) => {
    if (!token) return;
    
    setLoading(true);
    try {
      const response = await fetch(`/v1/usuarios/${userId}/eventos`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      setEvents(data.eventos);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const updateEvent = async (eventId, eventData) => {
    if (!token) throw new Error('Token necessário');
    
    try {
      const response = await fetch(`/v1/eventos/${eventId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(eventData)
      });
      
      if (!response.ok) throw new Error('Erro ao atualizar evento');
      
      // Atualizar lista local
      setEvents(prev => prev.map(event => 
        event.id === eventId ? { ...event, ...eventData } : event
      ));
      
      return true;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  return {
    events,
    loading,
    error,
    createEvent,
    fetchUserEvents,
    updateEvent
  };
};
```

### Componente de Criação de Evento

```javascript
// components/EventCreationForm.jsx
import React, { useState } from 'react';
import { useEvents } from '../hooks/useEvents';
import { useAuth } from '../hooks/useAuth';

const EventCreationForm = ({ onEventCreated }) => {
  const { token } = useAuth();
  const { createEvent, loading, error } = useEvents(token);
  const [formData, setFormData] = useState({
    nome: '',
    descricao: '',
    data_evento: '',
    local: '',
    url_slug: '',
    template_id: 'template_moderno'
  });
  const [validationErrors, setValidationErrors] = useState({});

  const generateSlug = (nome) => {
    return nome
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  };

  const validateForm = () => {
    const errors = {};
    
    if (!formData.nome.trim()) {
      errors.nome = 'Nome é obrigatório';
    }
    
    if (!formData.data_evento) {
      errors.data_evento = 'Data do evento é obrigatória';
    } else {
      const eventDate = new Date(formData.data_evento);
      if (eventDate <= new Date()) {
        errors.data_evento = 'Data deve ser no futuro';
      }
    }
    
    if (!formData.local.trim()) {
      errors.local = 'Local é obrigatório';
    }
    
    if (!formData.url_slug.trim()) {
      errors.url_slug = 'URL é obrigatória';
    } else if (!/^[a-z0-9-]+$/.test(formData.url_slug)) {
      errors.url_slug = 'URL deve conter apenas letras minúsculas, números e hífens';
    }
    
    return errors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const errors = validateForm();
    if (Object.keys(errors).length > 0) {
      setValidationErrors(errors);
      return;
    }
    
    try {
      const result = await createEvent(formData);
      onEventCreated?.(result);
      // Reset form
      setFormData({
        nome: '',
        descricao: '',
        data_evento: '',
        local: '',
        url_slug: '',
        template_id: 'template_moderno'
      });
      setValidationErrors({});
    } catch (error) {
      console.error('Erro ao criar evento:', error);
    }
  };

  const handleNameChange = (e) => {
    const nome = e.target.value;
    const slug = generateSlug(nome);
    setFormData({
      ...formData,
      nome,
      url_slug: slug
    });
  };

  return (
    <form onSubmit={handleSubmit} className="event-form">
      <h2>Criar Novo Evento</h2>
      
      {error && (
        <div className="error-message">{error}</div>
      )}
      
      <div className="form-group">
        <label>Nome do Evento:</label>
        <input
          type="text"
          value={formData.nome}
          onChange={handleNameChange}
          className={validationErrors.nome ? 'error' : ''}
          placeholder="Ex: Casamento João e Maria"
          disabled={loading}
          required
        />
        {validationErrors.nome && (
          <span className="field-error">{validationErrors.nome}</span>
        )}
      </div>

      <div className="form-group">
        <label>Descrição:</label>
        <textarea
          value={formData.descricao}
          onChange={(e) => setFormData({...formData, descricao: e.target.value})}
          placeholder="Descrição do evento (opcional)"
          rows={3}
          disabled={loading}
        />
      </div>

      <div className="form-row">
        <div className="form-group">
          <label>Data e Hora:</label>
          <input
            type="datetime-local"
            value={formData.data_evento}
            onChange={(e) => setFormData({...formData, data_evento: e.target.value})}
            className={validationErrors.data_evento ? 'error' : ''}
            disabled={loading}
            required
          />
          {validationErrors.data_evento && (
            <span className="field-error">{validationErrors.data_evento}</span>
          )}
        </div>

        <div className="form-group">
          <label>Template:</label>
          <select
            value={formData.template_id}
            onChange={(e) => setFormData({...formData, template_id: e.target.value})}
            disabled={loading}
          >
            <option value="template_moderno">Moderno</option>
            <option value="template_classico">Clássico</option>
            <option value="template_romantico">Romântico</option>
          </select>
        </div>
      </div>

      <div className="form-group">
        <label>Local:</label>
        <input
          type="text"
          value={formData.local}
          onChange={(e) => setFormData({...formData, local: e.target.value})}
          className={validationErrors.local ? 'error' : ''}
          placeholder="Ex: Igreja São Pedro, São Paulo"
          disabled={loading}
          required
        />
        {validationErrors.local && (
          <span className="field-error">{validationErrors.local}</span>
        )}
      </div>

      <div className="form-group">
        <label>URL Personalizada:</label>
        <div className="url-preview">
          <span className="url-base">meusite.com/</span>
          <input
            type="text"
            value={formData.url_slug}
            onChange={(e) => setFormData({...formData, url_slug: e.target.value.toLowerCase()})}
            className={validationErrors.url_slug ? 'error' : ''}
            placeholder="casamento-joao-maria-2024"
            disabled={loading}
            required
          />
        </div>
        {validationErrors.url_slug && (
          <span className="field-error">{validationErrors.url_slug}</span>
        )}
        <small className="help-text">
          Esta será a URL pública do seu evento
        </small>
      </div>

      <button type="submit" disabled={loading} className="submit-button">
        {loading ? 'Criando...' : 'Criar Evento'}
      </button>
    </form>
  );
};

export default EventCreationForm;
```

## 🎨 Estilos CSS

```css
.event-form {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
}

.form-group input.error,
.form-group textarea.error {
  border-color: #f44336;
}

.url-preview {
  display: flex;
  align-items: center;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
}

.url-base {
  background: #f5f5f5;
  padding: 12px;
  color: #666;
  border-right: 1px solid #ddd;
  white-space: nowrap;
}

.url-preview input {
  border: none;
  border-radius: 0;
  flex: 1;
}

.field-error {
  display: block;
  color: #f44336;
  font-size: 12px;
  margin-top: 5px;
}

.help-text {
  display: block;
  color: #666;
  font-size: 12px;
  margin-top: 5px;
  font-style: italic;
}

.submit-button {
  width: 100%;
  background: #2196F3;
  color: white;
  border: none;
  padding: 15px;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submit-button:hover:not(:disabled) {
  background: #1976D2;
}

.submit-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.error-message {
  background: #ffebee;
  color: #c62828;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 20px;
  border: 1px solid #ffcdd2;
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .url-preview {
    flex-direction: column;
  }
  
  .url-base {
    border-right: none;
    border-bottom: 1px solid #ddd;
  }
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | Página renderizada com sucesso |
| 201 | Criado | Evento criado com sucesso |
| 400 | Bad Request | Dados inválidos (URL slug inválida) |
| 401 | Unauthorized | Token JWT inválido |
| 404 | Not Found | Evento não encontrado |
| 409 | Conflict | URL slug já existe |

## 📱 Considerações para UX

### 1. **URL Amigável**
- Geração automática de slug a partir do nome
- Validação em tempo real
- Preview da URL final

### 2. **Templates**
- Seleção visual de templates
- Preview em tempo real
- Customização de cores

### 3. **Validação**
- Verificação de disponibilidade de URL
- Validação de data futura
- Campos obrigatórios claros

## 🔐 Segurança

### 1. **Validação de URL**
- Sanitização de slug
- Verificação de caracteres permitidos
- Prevenção de URLs maliciosas

### 2. **Autorização**
- Verificação de propriedade
- Token JWT válido
- Rate limiting

---

Para testes, consulte `tests/http/events.http`.