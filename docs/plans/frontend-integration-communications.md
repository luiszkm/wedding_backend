# 📢 Documentação de Integração Frontend - Módulo de Comunicações

## Visão Geral

O módulo de Comunicações permite criar e gerenciar comunicados/anúncios para os convidados. Inclui CRUD completo para administradores e visualização pública dos comunicados aprovados.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **Listar Comunicados Públicos**

**Endpoint:** `GET /eventos/{idEvento}/comunicados/publico`

**Resposta:**
```json
{
  "comunicados": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "titulo": "Atualização sobre o Local",
      "conteudo": "Informamos que o local da cerimônia foi alterado...",
      "data_publicacao": "2024-01-15T10:30:00Z",
      "prioritario": true
    }
  ]
}
```

### 2. ➕ **Criar Comunicado (Autenticado)**

**Endpoint:** `POST /eventos/{idEvento}/comunicados`

**Body:**
```json
{
  "titulo": "Atualização sobre o Local",
  "conteudo": "Informamos que o local da cerimônia foi alterado para...",
  "prioritario": false
}
```

### 3. ✏️ **Atualizar Comunicado (Autenticado)**

**Endpoint:** `PUT /comunicados/{idComunicado}`

### 4. 🗑️ **Deletar Comunicado (Autenticado)**

**Endpoint:** `DELETE /comunicados/{idComunicado}`

## 🎨 Componente React

### Hook para Comunicações

```javascript
// hooks/useCommunications.js
import { useState } from 'react';

export const useCommunications = (eventId, token = null) => {
  const [communications, setCommunications] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPublicCommunications = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/v1/eventos/${eventId}/comunicados/publico`);
      const data = await response.json();
      setCommunications(data.comunicados);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createCommunication = async (commData) => {
    if (!token) throw new Error('Token necessário');
    
    try {
      const response = await fetch(`/v1/eventos/${eventId}/comunicados`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(commData)
      });
      
      if (!response.ok) throw new Error('Erro ao criar comunicado');
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  return {
    communications,
    loading,
    error,
    fetchPublicCommunications,
    createCommunication
  };
};
```

### Componente Público de Comunicados

```javascript
// components/PublicCommunications.jsx
import React, { useEffect } from 'react';
import { useCommunications } from '../hooks/useCommunications';

const PublicCommunications = ({ eventId }) => {
  const { communications, loading, error, fetchPublicCommunications } = useCommunications(eventId);

  useEffect(() => {
    fetchPublicCommunications();
  }, [eventId]);

  if (loading) return <div className="loading">Carregando comunicados...</div>;
  if (error) return <div className="error">Erro: {error}</div>;

  const priorityCommunications = communications.filter(c => c.prioritario);
  const regularCommunications = communications.filter(c => !c.prioritario);

  return (
    <div className="communications">
      <h2>Comunicados</h2>
      
      {priorityCommunications.length > 0 && (
        <div className="priority-section">
          <h3>🚨 Importantes</h3>
          {priorityCommunications.map(comm => (
            <CommunicationCard key={comm.id} communication={comm} priority />
          ))}
        </div>
      )}
      
      {regularCommunications.length > 0 && (
        <div className="regular-section">
          {priorityCommunications.length > 0 && <h3>📋 Informativos</h3>}
          {regularCommunications.map(comm => (
            <CommunicationCard key={comm.id} communication={comm} />
          ))}
        </div>
      )}
      
      {communications.length === 0 && (
        <div className="empty-state">
          Nenhum comunicado disponível no momento.
        </div>
      )}
    </div>
  );
};

const CommunicationCard = ({ communication, priority = false }) => (
  <div className={`communication-card ${priority ? 'priority' : ''}`}>
    <div className="communication-header">
      <h4>{communication.titulo}</h4>
      <span className="date">
        {new Date(communication.data_publicacao).toLocaleDateString('pt-BR')}
      </span>
    </div>
    <div className="communication-content">
      {communication.conteudo}
    </div>
  </div>
);

export default PublicCommunications;
```

## 🎨 Estilos CSS

```css
.communications {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.priority-section,
.regular-section {
  margin-bottom: 30px;
}

.priority-section h3 {
  color: #f44336;
  border-left: 4px solid #f44336;
  padding-left: 10px;
}

.communication-card {
  background: white;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 15px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.communication-card.priority {
  border-left: 4px solid #f44336;
  background: #fff8f8;
}

.communication-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.communication-header h4 {
  margin: 0;
  color: #333;
}

.date {
  color: #666;
  font-size: 12px;
}

.communication-content {
  line-height: 1.6;
  color: #555;
  white-space: pre-wrap;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
  font-style: italic;
}

@media (max-width: 768px) {
  .communication-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
```

## 📱 Considerações para UX

### 1. **Priorização**
- Comunicados prioritários destacados
- Ordenação por data
- Notificações visuais

### 2. **Formatação**
- Suporte a quebras de linha
- Links automáticos
- Rich text editor para admin

### 3. **Responsividade**
- Layout adaptável
- Legibilidade mobile
- Carregamento eficiente

---

Para testes, consulte `tests/http/communications.http`.