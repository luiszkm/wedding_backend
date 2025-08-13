# 👥 Documentação de Integração Frontend - Módulo de Gestão de Convidados

## Visão Geral

O módulo de Gestão de Convidados permite criar grupos de convidados com chaves de acesso únicas para RSVP. Inclui funcionalidades públicas para confirmação de presença e interfaces administrativas para gerenciar grupos. O sistema utiliza chaves de acesso para permitir que convidados confirmem presença sem necessidade de registro individual.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📋 **RSVP - Confirmação de Presença (Público)**

**Endpoint:** `POST /rsvps`

**Descrição:** Permite que convidados confirmem presença usando chave de acesso. Endpoint público - não requer autenticação.

**Headers:**
```
Content-Type: application/json
```

**Body da Requisição:**
```json
{
  "chave_acesso": "padrinhos",
  "confirmacoes": [
    {
      "nome": "João Silva",
      "confirmado": true
    },
    {
      "nome": "Maria Silva", 
      "confirmado": false
    }
  ]
}
```

**Resposta de Sucesso (200):**
```json
{
  "message": "RSVP processado com sucesso",
  "confirmacoes_processadas": 2
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function submitRSVP(accessKey, confirmations) {
  try {
    const response = await fetch('/v1/rsvps', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        chave_acesso: accessKey,
        confirmacoes: confirmations
      })
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao processar RSVP:', error);
    throw error;
  }
}
```

---

### 2. 🔍 **Obter Grupo por Chave de Acesso (Público)**

**Endpoint:** `GET /acesso-convidado?chave={chaveAcesso}`

**Descrição:** Retorna dados do grupo de convidados usando a chave de acesso. Usado para exibir nomes no formulário de RSVP.

**Parâmetros de Query:**
- `chave`: Chave de acesso do grupo (ex: "padrinhos")

**Resposta de Sucesso (200):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "chave_acesso": "padrinhos",
  "nomes": ["João Silva", "Maria Silva", "Pedro Santos"],
  "id_casamento": "456e7890-e89b-12d3-a456-426614174001"
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function getGuestGroup(accessKey) {
  try {
    const response = await fetch(`/v1/acesso-convidado?chave=${encodeURIComponent(accessKey)}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao buscar grupo:', error);
    throw error;
  }
}
```

---

### 3. ➕ **Criar Grupo de Convidados (Autenticado)**

**Endpoint:** `POST /casamentos/{idCasamento}/grupos-de-convidados`

**Descrição:** Cria um novo grupo de convidados. Requer autenticação JWT.

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
  "chave_acesso": "padrinhos",
  "nomes": ["João Silva", "Maria Silva", "Pedro Santos"]
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
async function createGuestGroup(weddingId, groupData, token) {
  try {
    const response = await fetch(`/v1/casamentos/${weddingId}/grupos-de-convidados`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(groupData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao criar grupo:', error);
    throw error;
  }
}
```

---

### 4. ✏️ **Atualizar Grupo de Convidados (Autenticado)**

**Endpoint:** `PUT /grupos-de-convidados/{idGrupo}`

**Descrição:** Atualiza um grupo existente de convidados. Requer autenticação e propriedade do casamento.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Parâmetros:**
- `idGrupo` (path): UUID do grupo

**Body da Requisição:**
```json
{
  "nomes": ["João Silva", "Maria Silva", "Ana Santos"]
}
```

**Resposta de Sucesso (204):** Sem conteúdo

**Exemplo de uso (JavaScript):**
```javascript
async function updateGuestGroup(groupId, groupData, token) {
  try {
    const response = await fetch(`/v1/grupos-de-convidados/${groupId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(groupData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return true;
  } catch (error) {
    console.error('Erro ao atualizar grupo:', error);
    throw error;
  }
}
```

## 🎨 Componentes React

### Hook Customizado para Gestão de Convidados

```javascript
// hooks/useGuests.js
import { useState, useEffect } from 'react';

export const useGuests = (weddingId, token = null) => {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const getGuestGroup = async (accessKey) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/v1/acesso-convidado?chave=${encodeURIComponent(accessKey)}`);
      
      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Chave de acesso não encontrada');
        }
        throw new Error('Erro ao buscar grupo de convidados');
      }
      
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const submitRSVP = async (accessKey, confirmations) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/v1/rsvps', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          chave_acesso: accessKey,
          confirmacoes: confirmations
        })
      });
      
      if (!response.ok) {
        throw new Error('Erro ao processar RSVP');
      }
      
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const createGroup = async (groupData) => {
    if (!token) throw new Error('Token necessário para criar grupo');
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/v1/casamentos/${weddingId}/grupos-de-convidados`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(groupData)
      });
      
      if (!response.ok) {
        throw new Error('Erro ao criar grupo');
      }
      
      return await response.json();
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const updateGroup = async (groupId, groupData) => {
    if (!token) throw new Error('Token necessário para atualizar grupo');
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/v1/grupos-de-convidados/${groupId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(groupData)
      });
      
      if (!response.ok) {
        throw new Error('Erro ao atualizar grupo');
      }
      
      return true;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    groups,
    loading,
    error,
    getGuestGroup,
    submitRSVP,
    createGroup,
    updateGroup
  };
};
```

### Componente RSVP Público

```javascript
// components/RSVPForm.jsx
import React, { useState } from 'react';
import { useGuests } from '../hooks/useGuests';

const RSVPForm = () => {
  const { getGuestGroup, submitRSVP, loading, error } = useGuests();
  const [step, setStep] = useState('access-key'); // 'access-key' | 'confirmations' | 'success'
  const [accessKey, setAccessKey] = useState('');
  const [guestGroup, setGuestGroup] = useState(null);
  const [confirmations, setConfirmations] = useState([]);
  const [submitting, setSubmitting] = useState(false);

  const handleAccessKeySubmit = async (e) => {
    e.preventDefault();
    
    if (!accessKey.trim()) {
      return;
    }
    
    try {
      const group = await getGuestGroup(accessKey.trim());
      setGuestGroup(group);
      
      // Inicializar confirmações com todos como não confirmado
      const initialConfirmations = group.nomes.map(nome => ({
        nome,
        confirmado: false
      }));
      setConfirmations(initialConfirmations);
      setStep('confirmations');
    } catch (error) {
      console.error('Erro ao buscar grupo:', error);
    }
  };

  const handleConfirmationChange = (index, confirmed) => {
    setConfirmations(prev => prev.map((conf, i) => 
      i === index ? { ...conf, confirmado: confirmed } : conf
    ));
  };

  const handleRSVPSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    
    try {
      await submitRSVP(accessKey, confirmations);
      setStep('success');
    } catch (error) {
      console.error('Erro ao enviar RSVP:', error);
    } finally {
      setSubmitting(false);
    }
  };

  const resetForm = () => {
    setStep('access-key');
    setAccessKey('');
    setGuestGroup(null);
    setConfirmations([]);
  };

  if (step === 'success') {
    return (
      <div className="rsvp-success">
        <h2>✅ RSVP Confirmado!</h2>
        <p>Suas confirmações foram processadas com sucesso.</p>
        <button onClick={resetForm} className="secondary-button">
          Fazer Outro RSVP
        </button>
      </div>
    );
  }

  return (
    <div className="rsvp-form">
      <h2>Confirme sua Presença</h2>
      
      {error && (
        <div className="error-message">
          {error}
        </div>
      )}
      
      {step === 'access-key' && (
        <form onSubmit={handleAccessKeySubmit}>
          <div className="form-group">
            <label htmlFor="accessKey">
              Chave de Acesso:
            </label>
            <input
              type="text"
              id="accessKey"
              value={accessKey}
              onChange={(e) => setAccessKey(e.target.value)}
              placeholder="Digite sua chave de acesso"
              disabled={loading}
              required
            />
            <small className="help-text">
              A chave de acesso foi enviada junto com seu convite
            </small>
          </div>
          
          <button 
            type="submit" 
            disabled={loading || !accessKey.trim()}
            className="submit-button"
          >
            {loading ? 'Buscando...' : 'Continuar'}
          </button>
        </form>
      )}
      
      {step === 'confirmations' && guestGroup && (
        <form onSubmit={handleRSVPSubmit}>
          <div className="guest-group-info">
            <h3>Grupo: {guestGroup.chave_acesso}</h3>
            <p>Confirme a presença para cada pessoa:</p>
          </div>
          
          <div className="confirmations-list">
            {confirmations.map((conf, index) => (
              <div key={index} className="confirmation-item">
                <span className="guest-name">{conf.nome}</span>
                <div className="confirmation-buttons">
                  <label className="radio-label">
                    <input
                      type="radio"
                      name={`guest-${index}`}
                      checked={conf.confirmado === true}
                      onChange={() => handleConfirmationChange(index, true)}
                    />
                    <span className="confirm-yes">✅ Sim</span>
                  </label>
                  <label className="radio-label">
                    <input
                      type="radio"
                      name={`guest-${index}`}
                      checked={conf.confirmado === false}
                      onChange={() => handleConfirmationChange(index, false)}
                    />
                    <span className="confirm-no">❌ Não</span>
                  </label>
                </div>
              </div>
            ))}
          </div>
          
          <div className="form-actions">
            <button 
              type="button" 
              onClick={() => setStep('access-key')}
              className="secondary-button"
            >
              Voltar
            </button>
            <button 
              type="submit" 
              disabled={submitting}
              className="submit-button"
            >
              {submitting ? 'Enviando...' : 'Confirmar RSVP'}
            </button>
          </div>
        </form>
      )}
    </div>
  );
};

export default RSVPForm;
```

### Componente Administrativo de Grupos

```javascript
// components/GuestGroupAdmin.jsx
import React, { useState } from 'react';
import { useGuests } from '../hooks/useGuests';

const GuestGroupAdmin = ({ weddingId, token }) => {
  const { createGroup, updateGroup, loading, error } = useGuests(weddingId, token);
  const [showForm, setShowForm] = useState(false);
  const [editingGroup, setEditingGroup] = useState(null);
  const [formData, setFormData] = useState({
    chave_acesso: '',
    nomes: ['']
  });
  const [groups, setGroups] = useState([]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Filtrar nomes vazios
    const filteredNames = formData.nomes.filter(nome => nome.trim());
    
    if (filteredNames.length === 0) {
      alert('Adicione pelo menos um nome');
      return;
    }
    
    try {
      const groupData = {
        ...formData,
        nomes: filteredNames
      };
      
      if (editingGroup) {
        await updateGroup(editingGroup.id, { nomes: filteredNames });
        setGroups(prev => prev.map(g => 
          g.id === editingGroup.id 
            ? { ...g, nomes: filteredNames }
            : g
        ));
      } else {
        const newGroup = await createGroup(groupData);
        setGroups(prev => [...prev, { ...groupData, id: newGroup.id }]);
      }
      
      resetForm();
    } catch (error) {
      console.error('Erro ao salvar grupo:', error);
    }
  };

  const handleEdit = (group) => {
    setEditingGroup(group);
    setFormData({
      chave_acesso: group.chave_acesso,
      nomes: [...group.nomes]
    });
    setShowForm(true);
  };

  const resetForm = () => {
    setFormData({
      chave_acesso: '',
      nomes: ['']
    });
    setEditingGroup(null);
    setShowForm(false);
  };

  const addNameField = () => {
    setFormData(prev => ({
      ...prev,
      nomes: [...prev.nomes, '']
    }));
  };

  const removeNameField = (index) => {
    setFormData(prev => ({
      ...prev,
      nomes: prev.nomes.filter((_, i) => i !== index)
    }));
  };

  const updateName = (index, value) => {
    setFormData(prev => ({
      ...prev,
      nomes: prev.nomes.map((nome, i) => i === index ? value : nome)
    }));
  };

  return (
    <div className="guest-group-admin">
      <div className="header">
        <h2>Gerenciar Grupos de Convidados</h2>
        <button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancelar' : 'Novo Grupo'}
        </button>
      </div>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {showForm && (
        <form onSubmit={handleSubmit} className="group-form">
          <h3>{editingGroup ? 'Editar Grupo' : 'Novo Grupo'}</h3>
          
          {!editingGroup && (
            <div className="form-group">
              <label>Chave de Acesso:</label>
              <input
                type="text"
                value={formData.chave_acesso}
                onChange={(e) => setFormData({...formData, chave_acesso: e.target.value})}
                placeholder="Ex: padrinhos, familia-noivo"
                required
                disabled={loading}
              />
              <small className="help-text">
                Chave única que os convidados usarão para acessar o RSVP
              </small>
            </div>
          )}

          <div className="form-group">
            <label>Nomes dos Convidados:</label>
            {formData.nomes.map((nome, index) => (
              <div key={index} className="name-field">
                <input
                  type="text"
                  value={nome}
                  onChange={(e) => updateName(index, e.target.value)}
                  placeholder="Nome completo do convidado"
                  disabled={loading}
                />
                {formData.nomes.length > 1 && (
                  <button
                    type="button"
                    onClick={() => removeNameField(index)}
                    className="remove-button"
                  >
                    ✕
                  </button>
                )}
              </div>
            ))}
            
            <button
              type="button"
              onClick={addNameField}
              className="add-button"
            >
              + Adicionar Nome
            </button>
          </div>

          <div className="form-actions">
            <button type="submit" disabled={loading}>
              {loading ? 'Salvando...' : (editingGroup ? 'Atualizar' : 'Criar')}
            </button>
            <button type="button" onClick={resetForm}>
              Cancelar
            </button>
          </div>
        </form>
      )}

      <div className="groups-list">
        {groups.map((group) => (
          <div key={group.id} className="group-card">
            <div className="group-header">
              <h4>🔑 {group.chave_acesso}</h4>
              <button onClick={() => handleEdit(group)}>
                Editar
              </button>
            </div>
            <div className="group-content">
              <p><strong>Convidados ({group.nomes.length}):</strong></p>
              <ul className="guests-list">
                {group.nomes.map((nome, index) => (
                  <li key={index}>{nome}</li>
                ))}
              </ul>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default GuestGroupAdmin;
```

## 🎨 Estilos CSS

```css
/* RSVP Form */
.rsvp-form {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: #fff;
}

.rsvp-success {
  text-align: center;
  padding: 40px 20px;
  border: 2px solid #4CAF50;
  border-radius: 8px;
  background: #f8fff8;
  color: #2e7d32;
}

.guest-group-info {
  background: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.confirmations-list {
  margin: 20px 0;
}

.confirmation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 10px;
  background: #fafafa;
}

.guest-name {
  font-weight: bold;
  flex: 1;
}

.confirmation-buttons {
  display: flex;
  gap: 15px;
}

.radio-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 5px;
}

.radio-label input[type="radio"] {
  margin-right: 5px;
}

.confirm-yes {
  color: #4CAF50;
  font-weight: bold;
}

.confirm-no {
  color: #f44336;
  font-weight: bold;
}

/* Admin Components */
.guest-group-admin {
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

.group-form {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
  border: 1px solid #ddd;
}

.name-field {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.name-field input {
  flex: 1;
}

.remove-button {
  background: #f44336;
  color: white;
  border: none;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-button {
  background: #2196F3;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.groups-list {
  display: grid;
  gap: 20px;
}

.group-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  background: #fff;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.guests-list {
  list-style: none;
  padding: 0;
  margin: 10px 0;
}

.guests-list li {
  padding: 5px 0;
  border-bottom: 1px solid #f0f0f0;
}

.guests-list li:last-child {
  border-bottom: none;
}

/* Common Styles */
.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.help-text {
  color: #666;
  font-size: 12px;
  margin-top: 5px;
  display: block;
}

.submit-button {
  background: #2196F3;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.secondary-button {
  background: #666;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 4px;
  cursor: pointer;
}

.error-message {
  background: #ffebee;
  color: #c62828;
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 15px;
  border: 1px solid #ffcdd2;
}

/* Responsive */
@media (max-width: 768px) {
  .confirmation-item {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
  
  .confirmation-buttons {
    width: 100%;
    justify-content: center;
  }
  
  .header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | RSVP processado com sucesso |
| 201 | Criado | Grupo criado com sucesso |
| 204 | Sem Conteúdo | Grupo atualizado com sucesso |
| 400 | Bad Request | Dados inválidos (chave vazia, nomes duplicados) |
| 401 | Unauthorized | Token JWT inválido ou ausente |
| 404 | Not Found | Chave de acesso ou grupo não encontrado |
| 409 | Conflict | Chave de acesso já existe |
| 500 | Internal Server Error | Erro interno do servidor |

### Exemplos de Respostas de Erro

```json
// Erro 404 - Chave não encontrada
{
  "error": "Chave de acesso não encontrada",
  "details": "A chave informada não existe no sistema"
}

// Erro 409 - Chave já existe
{
  "error": "Chave de acesso já existe",
  "details": "Esta chave já está sendo usada por outro grupo"
}

// Erro 400 - Dados inválidos
{
  "error": "Lista de nomes é obrigatória",
  "details": "Pelo menos um nome deve ser fornecido"
}
```

## 📱 Considerações para UX

### 1. **Fluxo de RSVP**
- Interface em etapas para melhor usabilidade
- Validação da chave antes de mostrar confirmações
- Feedback visual claro para confirmações

### 2. **Gestão de Nomes**
- Interface dinâmica para adicionar/remover nomes
- Validação de nomes duplicados
- Suporte para nomes compostos

### 3. **Estados de Loading**
- Indicadores durante busca de grupo
- Desabilitar controles durante submissão
- Feedback de sucesso/erro claro

### 4. **Acessibilidade**
- Labels apropriados para screen readers
- Navegação por teclado funcional
- Contraste adequado de cores

## 🔐 Segurança

### 1. **Chaves de Acesso**
- Chaves únicas por grupo
- Validação no backend
- Não expor IDs internos

### 2. **Autorização**
- Verificação de propriedade para operações admin
- Validação de token JWT
- Rate limiting recomendado

### 3. **Validação de Dados**
- Sanitização de nomes de convidados
- Validação de formato de chave
- Prevenção de injection

---

Esta documentação fornece tudo necessário para integrar o sistema de gestão de convidados no frontend. Para testes, consulte os arquivos HTTP em `tests/http/guests.http`.