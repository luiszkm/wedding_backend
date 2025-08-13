# 💬 Documentação de Integração Frontend - Módulo de Mural de Recados

## Visão Geral

O módulo de Mural de Recados permite que convidados deixem mensagens carinhosas para os noivos, com sistema de moderação administrativa. Inclui interface pública para visualização e envio, além de painel administrativo para aprovação.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **Listar Recados Públicos**

**Endpoint:** `GET /casamentos/{idCasamento}/recados/publico`

**Resposta de Sucesso (200):**
```json
{
  "recados": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "autor": "João Silva",
      "mensagem": "Parabéns pelo casamento! Que sejam muito felizes!",
      "aprovado": true,
      "data_criacao": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 2. ✍️ **Deixar Recado (Público)**

**Endpoint:** `POST /recados`

**Body da Requisição:**
```json
{
  "id_casamento": "456e7890-e89b-12d3-a456-426614174001",
  "autor": "Maria Santos",
  "mensagem": "Felicidades para o casal!"
}
```

### 3. 🔍 **Listar Recados Admin (Autenticado)**

**Endpoint:** `GET /casamentos/{idCasamento}/recados/admin`

### 4. ✅ **Moderar Recado (Autenticado)**

**Endpoint:** `PATCH /recados/{idRecado}`

**Body:** `{"aprovado": true}`

## 🎨 Componentes React

### Hook para Mural de Recados

```javascript
// hooks/useMessages.js
import { useState } from 'react';

export const useMessages = (weddingId, token = null) => {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPublicMessages = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/recados/publico`);
      const data = await response.json();
      setMessages(data.recados);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const sendMessage = async (author, message) => {
    try {
      const response = await fetch('/v1/recados', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id_casamento: weddingId,
          autor: author,
          mensagem: message
        })
      });
      
      if (!response.ok) throw new Error('Erro ao enviar recado');
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  const moderateMessage = async (messageId, approved) => {
    if (!token) return;
    
    try {
      await fetch(`/v1/recados/${messageId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ aprovado: approved })
      });
      
      setMessages(prev => prev.map(msg => 
        msg.id === messageId ? { ...msg, aprovado: approved } : msg
      ));
    } catch (err) {
      setError(err.message);
    }
  };

  return {
    messages,
    loading,
    error,
    fetchPublicMessages,
    sendMessage,
    moderateMessage
  };
};
```

### Componente Público do Mural

```javascript
// components/PublicMessageBoard.jsx
import React, { useEffect, useState } from 'react';
import { useMessages } from '../hooks/useMessages';

const PublicMessageBoard = ({ weddingId }) => {
  const { messages, loading, error, fetchPublicMessages, sendMessage } = useMessages(weddingId);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ autor: '', mensagem: '' });
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchPublicMessages();
  }, [weddingId]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!formData.autor.trim() || !formData.mensagem.trim()) return;
    
    setSubmitting(true);
    try {
      await sendMessage(formData.autor, formData.mensagem);
      setFormData({ autor: '', mensagem: '' });
      setShowForm(false);
      // Mensagem será visível após moderação
    } catch (error) {
      console.error('Erro ao enviar:', error);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return <div className="loading">Carregando recados...</div>;
  if (error) return <div className="error">Erro: {error}</div>;

  return (
    <div className="message-board">
      <div className="header">
        <h2>Mural de Recados</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancelar' : 'Deixar Recado'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="message-form">
          <div className="form-group">
            <input
              type="text"
              placeholder="Seu nome"
              value={formData.autor}
              onChange={(e) => setFormData({...formData, autor: e.target.value})}
              disabled={submitting}
              required
            />
          </div>
          <div className="form-group">
            <textarea
              placeholder="Sua mensagem para os noivos..."
              value={formData.mensagem}
              onChange={(e) => setFormData({...formData, mensagem: e.target.value})}
              rows={4}
              disabled={submitting}
              required
            />
          </div>
          <div className="form-actions">
            <button type="submit" disabled={submitting}>
              {submitting ? 'Enviando...' : 'Enviar Recado'}
            </button>
          </div>
          <small className="help-text">
            Seu recado será revisado antes de aparecer no mural.
          </small>
        </form>
      )}

      <div className="messages-list">
        {messages.length === 0 ? (
          <div className="empty-state">
            Seja o primeiro a deixar um recado!
          </div>
        ) : (
          messages.map((message) => (
            <div key={message.id} className="message-card">
              <div className="message-header">
                <h4>{message.autor}</h4>
                <span className="date">
                  {new Date(message.data_criacao).toLocaleDateString('pt-BR')}
                </span>
              </div>
              <p className="message-content">{message.mensagem}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default PublicMessageBoard;
```

## 🎨 Estilos CSS

```css
.message-board {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.message-form {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
  border: 1px solid #ddd;
}

.form-group {
  margin-bottom: 15px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: inherit;
  font-size: 14px;
}

.help-text {
  color: #666;
  font-style: italic;
  margin-top: 10px;
  display: block;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message-card {
  background: white;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.message-header h4 {
  margin: 0;
  color: #333;
  font-size: 16px;
}

.date {
  color: #666;
  font-size: 12px;
}

.message-content {
  margin: 0;
  line-height: 1.6;
  color: #555;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #999;
  font-style: italic;
  background: #f9f9f9;
  border-radius: 8px;
}

.form-actions button {
  background: #2196F3;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.form-actions button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
  
  .message-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | Recados carregados com sucesso |
| 201 | Criado | Recado enviado com sucesso |
| 400 | Bad Request | Dados inválidos (mensagem vazia) |
| 401 | Unauthorized | Token JWT inválido (moderação) |
| 404 | Not Found | Recado ou casamento não encontrado |

## 📱 Considerações para UX

### 1. **Moderação Transparente**
- Informar que recados passam por moderação
- Feedback de envio bem-sucedido
- Interface admin para aprovar/rejeitar

### 2. **Validação**
- Campos obrigatórios
- Limite de caracteres
- Sanitização contra spam

### 3. **Performance**
- Paginação para muitos recados
- Loading states
- Otimização mobile

## 🔐 Segurança

### 1. **Moderação**
- Todos os recados passam por aprovação
- Filtros automáticos de conteúdo
- Rate limiting

### 2. **Validação**
- Sanitização de entrada
- Proteção XSS
- Verificação de spam

---

Para testes, consulte `tests/http/messages.http`.