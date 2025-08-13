# 📋 Documentação de Integração Frontend - Módulo de Roteiro/Itinerário

## Visão Geral

O módulo de Roteiro/Itinerário permite que os anfitriões criem e gerenciem um cronograma detalhado do dia do casamento, que pode ser visualizado publicamente pelos convidados. Esta funcionalidade foi implementada seguindo a **ADR-005** e oferece endpoints para operações CRUD completas.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📖 **Listar Roteiro (Público)**

**Endpoint:** `GET /eventos/{idEvento}/roteiro`

**Descrição:** Retorna todos os itens do roteiro de um evento, ordenados por horário. Acesso público - não requer autenticação.

**Parâmetros:**
- `idEvento` (path): UUID do evento

**Resposta de Sucesso (200):**
```json
{
  "items": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "idEvento": "123e4567-e89b-12d3-a456-426614174001",
      "horario": "2024-12-25T15:00:00Z",
      "tituloAtividade": "Cerimônia Religiosa",
      "descricaoAtividade": "Cerimônia na Igreja São José",
      "createdAt": "2024-08-06T10:30:00Z",
      "updatedAt": "2024-08-06T10:30:00Z"
    }
  ]
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function fetchItinerary(eventId) {
  try {
    const response = await fetch(`/v1/eventos/${eventId}/roteiro`);
    const data = await response.json();
    return data.items;
  } catch (error) {
    console.error('Erro ao carregar roteiro:', error);
    return [];
  }
}
```

---

### 2. ➕ **Criar Item do Roteiro (Autenticado)**

**Endpoint:** `POST /eventos/{idEvento}/roteiro`

**Descrição:** Cria um novo item no roteiro do evento. Requer autenticação JWT.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Parâmetros:**
- `idEvento` (path): UUID do evento

**Body da Requisição:**
```json
{
  "horario": "2024-12-25T15:00:00Z",
  "tituloAtividade": "Cerimônia Religiosa",
  "descricaoAtividade": "Cerimônia na Igreja São José" // Opcional
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000"
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function createItineraryItem(eventId, itemData, token) {
  try {
    const response = await fetch(`/v1/eventos/${eventId}/roteiro`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(itemData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao criar item do roteiro:', error);
    throw error;
  }
}
```

---

### 3. ✏️ **Atualizar Item do Roteiro (Autenticado)**

**Endpoint:** `PUT /roteiro/{idItemRoteiro}`

**Descrição:** Atualiza um item existente do roteiro. Requer autenticação e propriedade do evento.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Parâmetros:**
- `idItemRoteiro` (path): UUID do item do roteiro

**Body da Requisição:**
```json
{
  "horario": "2024-12-25T15:30:00Z",
  "tituloAtividade": "Cerimônia Religiosa - Atualizada",
  "descricaoAtividade": "Nova descrição" // Opcional
}
```

**Resposta de Sucesso (204):** Sem conteúdo

**Exemplo de uso (JavaScript):**
```javascript
async function updateItineraryItem(itemId, itemData, token) {
  try {
    const response = await fetch(`/v1/roteiro/${itemId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(itemData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return true;
  } catch (error) {
    console.error('Erro ao atualizar item do roteiro:', error);
    throw error;
  }
}
```

---

### 4. 🗑️ **Deletar Item do Roteiro (Autenticado)**

**Endpoint:** `DELETE /roteiro/{idItemRoteiro}`

**Descrição:** Remove um item do roteiro. Requer autenticação e propriedade do evento.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Parâmetros:**
- `idItemRoteiro` (path): UUID do item do roteiro

**Resposta de Sucesso (204):** Sem conteúdo

**Exemplo de uso (JavaScript):**
```javascript
async function deleteItineraryItem(itemId, token) {
  try {
    const response = await fetch(`/v1/roteiro/${itemId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return true;
  } catch (error) {
    console.error('Erro ao deletar item do roteiro:', error);
    throw error;
  }
}
```

## 🎨 Componente React de Exemplo

### Hook Customizado para Gerenciar Roteiro

```javascript
// hooks/useItinerary.js
import { useState, useEffect } from 'react';

export const useItinerary = (eventId, token = null) => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchItinerary = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/v1/eventos/${eventId}/roteiro`);
      if (!response.ok) throw new Error('Falha ao carregar roteiro');
      const data = await response.json();
      setItems(data.items);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createItem = async (itemData) => {
    if (!token) throw new Error('Token necessário para criar item');
    
    try {
      const response = await fetch(`/v1/eventos/${eventId}/roteiro`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(itemData)
      });
      
      if (!response.ok) throw new Error('Falha ao criar item');
      await fetchItinerary(); // Recarrega a lista
    } catch (err) {
      setError(err.message);
    }
  };

  const updateItem = async (itemId, itemData) => {
    if (!token) throw new Error('Token necessário para atualizar item');
    
    try {
      const response = await fetch(`/v1/roteiro/${itemId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(itemData)
      });
      
      if (!response.ok) throw new Error('Falha ao atualizar item');
      await fetchItinerary(); // Recarrega a lista
    } catch (err) {
      setError(err.message);
    }
  };

  const deleteItem = async (itemId) => {
    if (!token) throw new Error('Token necessário para deletar item');
    
    try {
      const response = await fetch(`/v1/roteiro/${itemId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (!response.ok) throw new Error('Falha ao deletar item');
      await fetchItinerary(); // Recarrega a lista
    } catch (err) {
      setError(err.message);
    }
  };

  useEffect(() => {
    if (eventId) fetchItinerary();
  }, [eventId]);

  return {
    items,
    loading,
    error,
    createItem,
    updateItem,
    deleteItem,
    refetch: fetchItinerary
  };
};
```

### Componente de Exibição Pública

```javascript
// components/PublicItinerary.jsx
import React from 'react';
import { useItinerary } from '../hooks/useItinerary';

const PublicItinerary = ({ eventId }) => {
  const { items, loading, error } = useItinerary(eventId);

  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString('pt-BR', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  if (loading) return <div className="loading">Carregando roteiro...</div>;
  if (error) return <div className="error">Erro: {error}</div>;
  if (!items.length) return <div className="empty">Nenhum item no roteiro ainda.</div>;

  return (
    <div className="itinerary">
      <h2>Roteiro do Evento</h2>
      <div className="timeline">
        {items.map((item, index) => (
          <div key={item.id} className="timeline-item">
            <div className="timeline-time">
              <strong>{formatTime(item.horario)}</strong>
              <small>{formatDate(item.horario)}</small>
            </div>
            <div className="timeline-content">
              <h3>{item.tituloAtividade}</h3>
              {item.descricaoAtividade && (
                <p>{item.descricaoAtividade}</p>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default PublicItinerary;
```

### Componente de Administração

```javascript
// components/ItineraryAdmin.jsx
import React, { useState } from 'react';
import { useItinerary } from '../hooks/useItinerary';

const ItineraryAdmin = ({ eventId, token }) => {
  const { items, loading, error, createItem, updateItem, deleteItem } = useItinerary(eventId, token);
  const [showForm, setShowForm] = useState(false);
  const [editingItem, setEditingItem] = useState(null);
  const [formData, setFormData] = useState({
    horario: '',
    tituloAtividade: '',
    descricaoAtividade: ''
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingItem) {
        await updateItem(editingItem.id, formData);
      } else {
        await createItem(formData);
      }
      resetForm();
    } catch (error) {
      console.error('Erro ao salvar item:', error);
    }
  };

  const handleEdit = (item) => {
    setEditingItem(item);
    setFormData({
      horario: item.horario.slice(0, 16), // Para input datetime-local
      tituloAtividade: item.tituloAtividade,
      descricaoAtividade: item.descricaoAtividade || ''
    });
    setShowForm(true);
  };

  const handleDelete = async (itemId) => {
    if (window.confirm('Tem certeza que deseja deletar este item?')) {
      await deleteItem(itemId);
    }
  };

  const resetForm = () => {
    setFormData({
      horario: '',
      tituloAtividade: '',
      descricaoAtividade: ''
    });
    setEditingItem(null);
    setShowForm(false);
  };

  return (
    <div className="itinerary-admin">
      <div className="header">
        <h2>Gerenciar Roteiro</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancelar' : 'Adicionar Item'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="itinerary-form">
          <h3>{editingItem ? 'Editar Item' : 'Novo Item'}</h3>
          
          <div className="form-group">
            <label>Horário:</label>
            <input
              type="datetime-local"
              value={formData.horario}
              onChange={(e) => setFormData({...formData, horario: e.target.value})}
              required
            />
          </div>

          <div className="form-group">
            <label>Título da Atividade:</label>
            <input
              type="text"
              value={formData.tituloAtividade}
              onChange={(e) => setFormData({...formData, tituloAtividade: e.target.value})}
              maxLength={255}
              required
            />
          </div>

          <div className="form-group">
            <label>Descrição (opcional):</label>
            <textarea
              value={formData.descricaoAtividade}
              onChange={(e) => setFormData({...formData, descricaoAtividade: e.target.value})}
              rows={3}
            />
          </div>

          <div className="form-actions">
            <button type="submit">
              {editingItem ? 'Atualizar' : 'Criar'}
            </button>
            <button type="button" onClick={resetForm}>
              Cancelar
            </button>
          </div>
        </form>
      )}

      {loading && <div>Carregando...</div>}
      {error && <div className="error">{error}</div>}

      <div className="items-list">
        {items.map((item) => (
          <div key={item.id} className="item-card">
            <div className="item-header">
              <h4>{item.tituloAtividade}</h4>
              <div className="item-actions">
                <button onClick={() => handleEdit(item)}>Editar</button>
                <button onClick={() => handleDelete(item.id)}>Deletar</button>
              </div>
            </div>
            <div className="item-details">
              <p><strong>Horário:</strong> {new Date(item.horario).toLocaleString('pt-BR')}</p>
              {item.descricaoAtividade && (
                <p><strong>Descrição:</strong> {item.descricaoAtividade}</p>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ItineraryAdmin;
```

## 🎨 Estilos CSS de Exemplo

```css
/* Estilos para timeline pública */
.itinerary {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.timeline {
  position: relative;
  padding-left: 30px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 0;
  height: 100%;
  width: 2px;
  background: #e0e0e0;
}

.timeline-item {
  position: relative;
  margin-bottom: 30px;
  padding-left: 40px;
}

.timeline-item::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 5px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #2196F3;
  border: 3px solid #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.timeline-time {
  font-size: 14px;
  color: #666;
  margin-bottom: 5px;
}

.timeline-content h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.timeline-content p {
  margin: 0;
  color: #666;
  line-height: 1.5;
}

/* Estilos para admin */
.itinerary-admin {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.itinerary-form {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-actions {
  display: flex;
  gap: 10px;
}

.item-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 15px;
  background: #fff;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.item-actions {
  display: flex;
  gap: 10px;
}

.item-actions button {
  padding: 5px 10px;
  font-size: 12px;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

button:hover {
  opacity: 0.8;
}

.error {
  color: #f44336;
  padding: 10px;
  background: #ffebee;
  border-radius: 4px;
  margin-bottom: 15px;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #999;
  font-style: italic;
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | Listagem de roteiro |
| 201 | Criado | Item criado com sucesso |
| 204 | Sem Conteúdo | Item atualizado/deletado |
| 400 | Bad Request | Dados inválidos (título vazio, muito longo, etc.) |
| 401 | Unauthorized | Token JWT inválido ou ausente |
| 404 | Not Found | Item ou evento não encontrado |
| 500 | Internal Server Error | Erro interno do servidor |

### Exemplos de Respostas de Erro

```json
// Erro 400 - Título obrigatório
{
  "codigo": "DADOS_INVALIDOS",
  "mensagem": "título da atividade é obrigatório",
  "status": 400
}

// Erro 401 - Não autenticado
{
  "codigo": "TOKEN_INVALIDO",
  "mensagem": "ID de usuário ausente ou inválido no token",
  "status": 401
}

// Erro 404 - Item não encontrado
{
  "codigo": "NAO_ENCONTRADO",
  "mensagem": "Item do roteiro não encontrado",
  "status": 404
}
```

## 📱 Considerações para UX

### 1. **Ordenação Automática**
- Os itens são sempre retornados ordenados por horário
- O frontend deve manter essa ordenação na interface

### 2. **Formatação de Datas**
- Use `toLocaleString()` para exibir datas no formato local
- Considere usar bibliotecas como `date-fns` ou `dayjs` para formatação avançada

### 3. **Estados de Carregamento**
- Implemente loading states durante as operações
- Use skeleton screens para melhor UX

### 4. **Validação no Frontend**
- Valide dados antes de enviar para a API
- Limite de 255 caracteres para título
- Campos obrigatórios: horário e título

### 5. **Confirmações**
- Sempre confirme antes de deletar itens
- Considere implementar undo para ações destrutivas

## 🔐 Segurança

### 1. **Autenticação**
- Rotas de escrita exigem token JWT válido
- Rota de leitura é pública para facilitar acesso dos convidados

### 2. **Autorização**
- Usuários só podem modificar itens de seus próprios eventos
- Verificação de propriedade feita no backend

### 3. **Validação**
- Validação duplicada (frontend + backend)
- Sanitização de dados no backend

---

Esta documentação fornece tudo que você precisa para integrar o módulo de Roteiro/Itinerário no seu frontend. Para dúvidas ou problemas, consulte os arquivos de teste em `tests/http/itinerary.http`.