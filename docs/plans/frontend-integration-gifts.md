# 🎁 Documentação de Integração Frontend - Módulo de Lista de Presentes

## Visão Geral

O módulo de Lista de Presentes permite criar e gerenciar uma lista de presentes para o casamento, incluindo suporte a presentes fracionados com sistema de cotas. Oferece interfaces públicas para visualização e seleção de presentes, além de funcionalidades administrativas para gestão completa da lista.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **Listar Presentes Públicos**

**Endpoint:** `GET /casamentos/{idCasamento}/presentes-publico`

**Descrição:** Retorna lista pública de presentes disponíveis para seleção. Endpoint público - não requer autenticação.

**Parâmetros:**
- `idCasamento` (path): UUID do casamento

**Resposta de Sucesso (200):**
```json
{
  "presentes": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "nome": "Jogo de Panelas Premium",
      "descricao": "Conjunto completo de panelas antiaderentes",
      "preco": 299.99,
      "imagem_url": "https://storage.example.com/gifts/panelas.jpg",
      "selecionado": false,
      "tipo": "COMPLETO",
      "cotas_disponiveis": null,
      "cotas_selecionadas": null,
      "valor_por_cota": null
    },
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "nome": "Lua de Mel - Passagens",
      "descricao": "Contribuição para passagens da lua de mel",
      "preco": 2000.00,
      "imagem_url": "https://storage.example.com/gifts/viagem.jpg",
      "selecionado": false,
      "tipo": "FRACIONADO",
      "cotas_disponiveis": 10,
      "cotas_selecionadas": 3,
      "valor_por_cota": 200.00
    }
  ]
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function fetchPublicGifts(weddingId) {
  try {
    const response = await fetch(`/v1/casamentos/${weddingId}/presentes-publico`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return data.presentes;
  } catch (error) {
    console.error('Erro ao carregar presentes:', error);
    throw error;
  }
}
```

---

### 2. 🛒 **Selecionar Presente**

**Endpoint:** `POST /selecoes-de-presente`

**Descrição:** Permite selecionar um presente completo ou cotas de um presente fracionado. Endpoint público - não requer autenticação.

**Headers:**
```
Content-Type: application/json
```

**Body da Requisição:**
```json
{
  "id_presente": "123e4567-e89b-12d3-a456-426614174000",
  "nome_selecionador": "João Silva",
  "email_selecionador": "joao@exemplo.com",
  "quantidade_cotas": 2,
  "mensagem": "Parabéns pelo casamento!"
}
```

**Campos:**
- `id_presente`: UUID do presente (obrigatório)
- `nome_selecionador`: Nome de quem está dando o presente (obrigatório)
- `email_selecionador`: Email do selecionador (obrigatório)
- `quantidade_cotas`: Número de cotas (apenas para presentes fracionados)
- `mensagem`: Mensagem opcional

**Resposta de Sucesso (201):**
```json
{
  "id": "789a0123-e89b-12d3-a456-426614174002",
  "valor_total": 400.00,
  "cotas_selecionadas": 2
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function selectGift(selectionData) {
  try {
    const response = await fetch('/v1/selecoes-de-presente', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(selectionData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao selecionar presente:', error);
    throw error;
  }
}
```

---

### 3. ➕ **Criar Presente (Autenticado)**

**Endpoint:** `POST /casamentos/{idCasamento}/presentes`

**Descrição:** Cria um novo presente na lista. Requer autenticação JWT.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Parâmetros:**
- `idCasamento` (path): UUID do casamento

**Body da Requisição:**
```json
{
  "nome": "Jogo de Panelas Premium",
  "descricao": "Conjunto completo de panelas antiaderentes com revestimento cerâmico",
  "preco": 299.99,
  "imagem": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQ...",
  "tipo": "COMPLETO",
  "numero_cotas": null,
  "valor_por_cota": null
}
```

**Campos para Presente Fracionado:**
```json
{
  "nome": "Lua de Mel - Hotel",
  "descricao": "Contribuição para hospedagem da lua de mel",
  "preco": 1500.00,
  "imagem": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQ...",
  "tipo": "FRACIONADO",
  "numero_cotas": 10,
  "valor_por_cota": 150.00
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
async function createGift(weddingId, giftData, token) {
  try {
    const response = await fetch(`/v1/casamentos/${weddingId}/presentes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(giftData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao criar presente:', error);
    throw error;
  }
}
```

---

### 4. 📊 **Listar Presentes Administrativo (Autenticado)**

**Endpoint:** `GET /casamentos/{idCasamento}/presentes`

**Descrição:** Lista todos os presentes com informações administrativas incluindo seleções. Requer autenticação JWT.

**Headers:**
```
Authorization: Bearer {jwt_token}
```

**Resposta de Sucesso (200):**
```json
{
  "presentes": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "nome": "Jogo de Panelas Premium",
      "descricao": "Conjunto completo de panelas antiaderentes",
      "preco": 299.99,
      "imagem_url": "https://storage.example.com/gifts/panelas.jpg",
      "tipo": "COMPLETO",
      "selecionado": true,
      "selecoes": [
        {
          "id": "sel-001",
          "nome_selecionador": "Maria Santos",
          "email_selecionador": "maria@exemplo.com",
          "data_selecao": "2024-01-15T10:30:00Z",
          "mensagem": "Felicidades!"
        }
      ]
    },
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "nome": "Lua de Mel - Passagens",
      "tipo": "FRACIONADO",
      "numero_cotas": 10,
      "valor_por_cota": 200.00,
      "cotas_selecionadas": 7,
      "cotas_disponiveis": 3,
      "selecoes": [
        {
          "id": "sel-002",
          "nome_selecionador": "João Silva",
          "quantidade_cotas": 3,
          "valor_total": 600.00,
          "data_selecao": "2024-01-16T14:20:00Z"
        }
      ]
    }
  ]
}
```

## 🎨 Componentes React

### Hook Customizado para Gestão de Presentes

```javascript
// hooks/useGifts.js
import { useState, useEffect } from 'react';

export const useGifts = (weddingId, token = null) => {
  const [gifts, setGifts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPublicGifts = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/presentes-publico`);
      
      if (!response.ok) {
        throw new Error('Erro ao carregar presentes');
      }
      
      const data = await response.json();
      setGifts(data.presentes);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchAdminGifts = async () => {
    if (!token) throw new Error('Token necessário para buscar dados administrativos');
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/presentes`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (!response.ok) {
        throw new Error('Erro ao carregar presentes');
      }
      
      const data = await response.json();
      setGifts(data.presentes);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const selectGift = async (selectionData) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/v1/selecoes-de-presente', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(selectionData)
      });
      
      if (!response.ok) {
        throw new Error('Erro ao selecionar presente');
      }
      
      const result = await response.json();
      
      // Atualizar lista após seleção
      await fetchPublicGifts();
      
      return result;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const createGift = async (giftData) => {
    if (!token) throw new Error('Token necessário para criar presente');
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/presentes`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(giftData)
      });
      
      if (!response.ok) {
        throw new Error('Erro ao criar presente');
      }
      
      const result = await response.json();
      
      // Recarregar lista
      await fetchAdminGifts();
      
      return result;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    gifts,
    loading,
    error,
    fetchPublicGifts,
    fetchAdminGifts,
    selectGift,
    createGift
  };
};
```

### Componente de Lista Pública de Presentes

```javascript
// components/PublicGiftList.jsx
import React, { useEffect, useState } from 'react';
import { useGifts } from '../hooks/useGifts';

const PublicGiftList = ({ weddingId }) => {
  const { gifts, loading, error, fetchPublicGifts, selectGift } = useGifts(weddingId);
  const [selectedGift, setSelectedGift] = useState(null);
  const [showSelectionForm, setShowSelectionForm] = useState(false);

  useEffect(() => {
    fetchPublicGifts();
  }, [weddingId]);

  const handleGiftClick = (gift) => {
    if (gift.tipo === 'COMPLETO' && gift.selecionado) {
      return; // Presente já selecionado
    }
    
    if (gift.tipo === 'FRACIONADO' && gift.cotas_disponiveis === 0) {
      return; // Sem cotas disponíveis
    }
    
    setSelectedGift(gift);
    setShowSelectionForm(true);
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL'
    }).format(price);
  };

  const getGiftStatus = (gift) => {
    if (gift.tipo === 'COMPLETO') {
      return gift.selecionado ? 'Selecionado' : 'Disponível';
    } else {
      const remaining = gift.cotas_disponiveis;
      if (remaining === 0) return 'Completo';
      return `${remaining} cota${remaining > 1 ? 's' : ''} disponível${remaining > 1 ? 'eis' : ''}`;
    }
  };

  const getGiftStatusClass = (gift) => {
    if (gift.tipo === 'COMPLETO') {
      return gift.selecionado ? 'selected' : 'available';
    } else {
      return gift.cotas_disponiveis === 0 ? 'selected' : 'available';
    }
  };

  if (loading) {
    return <div className="loading">Carregando lista de presentes...</div>;
  }

  if (error) {
    return <div className="error">Erro: {error}</div>;
  }

  return (
    <div className="public-gift-list">
      <h2>Lista de Presentes</h2>
      
      {gifts.length === 0 ? (
        <div className="empty-state">
          Nenhum presente disponível no momento.
        </div>
      ) : (
        <div className="gifts-grid">
          {gifts.map((gift) => (
            <div 
              key={gift.id} 
              className={`gift-card ${getGiftStatusClass(gift)}`}
              onClick={() => handleGiftClick(gift)}
            >
              <div className="gift-image">
                {gift.imagem_url ? (
                  <img src={gift.imagem_url} alt={gift.nome} />
                ) : (
                  <div className="placeholder-image">🎁</div>
                )}
              </div>
              
              <div className="gift-content">
                <h3 className="gift-name">{gift.nome}</h3>
                <p className="gift-description">{gift.descricao}</p>
                
                <div className="gift-pricing">
                  {gift.tipo === 'COMPLETO' ? (
                    <div className="complete-price">
                      <span className="price">{formatPrice(gift.preco)}</span>
                    </div>
                  ) : (
                    <div className="fractional-price">
                      <div className="per-quota">
                        {formatPrice(gift.valor_por_cota)} por cota
                      </div>
                      <div className="total-price">
                        Total: {formatPrice(gift.preco)}
                      </div>
                      <div className="quota-info">
                        {gift.cotas_selecionadas}/{gift.cotas_disponiveis + gift.cotas_selecionadas} cotas selecionadas
                      </div>
                    </div>
                  )}
                </div>
                
                <div className={`gift-status ${getGiftStatusClass(gift)}`}>
                  {getGiftStatus(gift)}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
      
      {showSelectionForm && selectedGift && (
        <GiftSelectionModal
          gift={selectedGift}
          onClose={() => {
            setShowSelectionForm(false);
            setSelectedGift(null);
          }}
          onSelect={selectGift}
          onSuccess={() => {
            setShowSelectionForm(false);
            setSelectedGift(null);
            fetchPublicGifts();
          }}
        />
      )}
    </div>
  );
};

// Modal de Seleção de Presente
const GiftSelectionModal = ({ gift, onClose, onSelect, onSuccess }) => {
  const [formData, setFormData] = useState({
    nome_selecionador: '',
    email_selecionador: '',
    quantidade_cotas: gift.tipo === 'FRACIONADO' ? 1 : undefined,
    mensagem: ''
  });
  const [submitting, setSubmitting] = useState(false);
  const [validationErrors, setValidationErrors] = useState({});

  const validateForm = () => {
    const errors = {};
    
    if (!formData.nome_selecionador.trim()) {
      errors.nome_selecionador = 'Nome é obrigatório';
    }
    
    if (!formData.email_selecionador.trim()) {
      errors.email_selecionador = 'Email é obrigatório';
    } else if (!/\S+@\S+\.\S+/.test(formData.email_selecionador)) {
      errors.email_selecionador = 'Email inválido';
    }
    
    if (gift.tipo === 'FRACIONADO') {
      const cotas = parseInt(formData.quantidade_cotas);
      if (!cotas || cotas < 1) {
        errors.quantidade_cotas = 'Quantidade de cotas deve ser pelo menos 1';
      } else if (cotas > gift.cotas_disponiveis) {
        errors.quantidade_cotas = `Máximo ${gift.cotas_disponiveis} cotas disponíveis`;
      }
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
    
    setSubmitting(true);
    
    try {
      const selectionData = {
        id_presente: gift.id,
        nome_selecionador: formData.nome_selecionador,
        email_selecionador: formData.email_selecionador,
        mensagem: formData.mensagem || undefined
      };
      
      if (gift.tipo === 'FRACIONADO') {
        selectionData.quantidade_cotas = parseInt(formData.quantidade_cotas);
      }
      
      await onSelect(selectionData);
      onSuccess();
    } catch (error) {
      console.error('Erro ao selecionar presente:', error);
    } finally {
      setSubmitting(false);
    }
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL'
    }).format(price);
  };

  const calculateTotal = () => {
    if (gift.tipo === 'COMPLETO') {
      return gift.preco;
    } else {
      const cotas = parseInt(formData.quantidade_cotas) || 1;
      return gift.valor_por_cota * cotas;
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h3>Selecionar Presente</h3>
          <button className="close-button" onClick={onClose}>✕</button>
        </div>
        
        <div className="gift-summary">
          <h4>{gift.nome}</h4>
          <p>{gift.descricao}</p>
          
          {gift.tipo === 'FRACIONADO' && (
            <div className="quota-summary">
              <p>Valor por cota: {formatPrice(gift.valor_por_cota)}</p>
              <p>Cotas disponíveis: {gift.cotas_disponiveis}</p>
            </div>
          )}
        </div>
        
        <form onSubmit={handleSubmit} className="selection-form">
          <div className="form-group">
            <label>Seu Nome:</label>
            <input
              type="text"
              value={formData.nome_selecionador}
              onChange={(e) => setFormData({...formData, nome_selecionador: e.target.value})}
              className={validationErrors.nome_selecionador ? 'error' : ''}
              disabled={submitting}
              required
            />
            {validationErrors.nome_selecionador && (
              <span className="field-error">{validationErrors.nome_selecionador}</span>
            )}
          </div>

          <div className="form-group">
            <label>Seu Email:</label>
            <input
              type="email"
              value={formData.email_selecionador}
              onChange={(e) => setFormData({...formData, email_selecionador: e.target.value})}
              className={validationErrors.email_selecionador ? 'error' : ''}
              disabled={submitting}
              required
            />
            {validationErrors.email_selecionador && (
              <span className="field-error">{validationErrors.email_selecionador}</span>
            )}
          </div>

          {gift.tipo === 'FRACIONADO' && (
            <div className="form-group">
              <label>Quantidade de Cotas:</label>
              <input
                type="number"
                min="1"
                max={gift.cotas_disponiveis}
                value={formData.quantidade_cotas}
                onChange={(e) => setFormData({...formData, quantidade_cotas: e.target.value})}
                className={validationErrors.quantidade_cotas ? 'error' : ''}
                disabled={submitting}
                required
              />
              {validationErrors.quantidade_cotas && (
                <span className="field-error">{validationErrors.quantidade_cotas}</span>
              )}
            </div>
          )}

          <div className="form-group">
            <label>Mensagem (opcional):</label>
            <textarea
              value={formData.mensagem}
              onChange={(e) => setFormData({...formData, mensagem: e.target.value})}
              rows={3}
              disabled={submitting}
              placeholder="Deixe uma mensagem carinhosa para os noivos..."
            />
          </div>

          <div className="total-summary">
            <strong>Total: {formatPrice(calculateTotal())}</strong>
          </div>

          <div className="form-actions">
            <button type="button" onClick={onClose} disabled={submitting}>
              Cancelar
            </button>
            <button type="submit" disabled={submitting} className="primary">
              {submitting ? 'Selecionando...' : 'Confirmar Seleção'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default PublicGiftList;
```

## 🎨 Estilos CSS

```css
/* Public Gift List */
.public-gift-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.public-gift-list h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.gifts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.gift-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  transition: all 0.3s ease;
  cursor: pointer;
}

.gift-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  transform: translateY(-2px);
}

.gift-card.selected {
  opacity: 0.7;
  cursor: not-allowed;
}

.gift-card.selected:hover {
  transform: none;
  box-shadow: none;
}

.gift-image {
  height: 200px;
  overflow: hidden;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gift-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder-image {
  font-size: 48px;
  color: #999;
}

.gift-content {
  padding: 20px;
}

.gift-name {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 18px;
  font-weight: bold;
}

.gift-description {
  color: #666;
  font-size: 14px;
  line-height: 1.4;
  margin-bottom: 15px;
}

.gift-pricing {
  margin-bottom: 15px;
}

.complete-price .price {
  font-size: 20px;
  font-weight: bold;
  color: #2196F3;
}

.fractional-price .per-quota {
  font-size: 16px;
  font-weight: bold;
  color: #2196F3;
}

.fractional-price .total-price {
  font-size: 14px;
  color: #666;
  margin: 5px 0;
}

.quota-info {
  font-size: 12px;
  color: #888;
}

.gift-status {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: bold;
  text-align: center;
}

.gift-status.available {
  background: #e8f5e8;
  color: #2e7d32;
}

.gift-status.selected {
  background: #ffebee;
  color: #c62828;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.close-button {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #666;
}

.gift-summary {
  padding: 20px;
  background: #f9f9f9;
  border-bottom: 1px solid #eee;
}

.quota-summary {
  margin-top: 10px;
  font-size: 14px;
  color: #666;
}

.selection-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input.error,
.form-group textarea.error {
  border-color: #f44336;
}

.field-error {
  display: block;
  color: #f44336;
  font-size: 12px;
  margin-top: 5px;
}

.total-summary {
  background: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  margin: 20px 0;
  text-align: center;
  font-size: 18px;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.form-actions button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.form-actions button:not(.primary) {
  background: #666;
  color: white;
}

.form-actions button.primary {
  background: #2196F3;
  color: white;
}

.form-actions button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Loading and Error States */
.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  text-align: center;
  padding: 40px;
  color: #f44336;
  background: #ffebee;
  border-radius: 4px;
  margin: 20px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #999;
  font-style: italic;
}

/* Responsive */
@media (max-width: 768px) {
  .gifts-grid {
    grid-template-columns: 1fr;
  }
  
  .modal-content {
    width: 95%;
    margin: 10px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .form-actions button {
    width: 100%;
  }
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | Lista carregada com sucesso |
| 201 | Criado | Presente criado ou seleção realizada |
| 400 | Bad Request | Dados inválidos (preço negativo, cotas inválidas) |
| 401 | Unauthorized | Token JWT inválido ou ausente |
| 404 | Not Found | Presente ou casamento não encontrado |
| 409 | Conflict | Presente já selecionado ou cotas insuficientes |
| 422 | Unprocessable Entity | Dados válidos mas não processáveis |
| 500 | Internal Server Error | Erro interno do servidor |

### Exemplos de Respostas de Erro

```json
// Erro 409 - Presente já selecionado
{
  "error": "Presente já foi selecionado",
  "details": "Este presente não está mais disponível"
}

// Erro 409 - Cotas insuficientes  
{
  "error": "Cotas insuficientes",
  "details": "Restam apenas 2 cotas disponíveis"
}

// Erro 400 - Dados inválidos
{
  "error": "Quantidade de cotas inválida",
  "details": "A quantidade deve ser entre 1 e o máximo disponível"
}
```

## 📱 Considerações para UX

### 1. **Interface Intuitiva**
- Cards visuais para cada presente
- Status claro de disponibilidade
- Diferenciação visual entre presentes completos e fracionados

### 2. **Sistema de Cotas**
- Indicador visual de cotas disponíveis/selecionadas
- Seletor fácil de quantidade de cotas
- Cálculo automático do valor total

### 3. **Feedback de Seleção**
- Confirmação visual após seleção
- Estados de loading durante processamento
- Mensagens de erro claras

### 4. **Responsividade**
- Grid adaptável para diferentes tamanhos de tela
- Modal responsivo
- Interface touch-friendly

## 🔐 Segurança

### 1. **Validação de Dados**
- Validação de email e campos obrigatórios
- Verificação de disponibilidade em tempo real
- Sanitização de dados de entrada

### 2. **Controle de Acesso**
- Endpoints públicos apenas para visualização e seleção
- Operações administrativas requerem autenticação
- Verificação de propriedade de casamento

### 3. **Prevenção de Fraudes**
- Rate limiting recomendado
- Validação de cotas disponíveis no backend
- Log de seleções para auditoria

---

Esta documentação fornece tudo necessário para integrar o sistema de lista de presentes no frontend. Para testes, consulte os arquivos HTTP em `tests/http/gifts.http` e `tests/http/gift-selections.http`.