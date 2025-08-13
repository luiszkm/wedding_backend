# 🔐 Documentação de Integração Frontend - Módulo IAM (Autenticação)

## Visão Geral

O módulo IAM (Identity and Access Management) gerencia a autenticação de usuários no sistema de casamentos. Fornece endpoints para registro, login e gerenciamento de sessões usando tokens JWT. Este é o módulo fundamental que protege todas as operações administrativas do sistema.

## Endpoints da API

### Base URL
```
http://localhost:3000/v1
```

### 1. 📝 **Registrar Usuário**

**Endpoint:** `POST /usuarios/registrar`

**Descrição:** Registra um novo usuário no sistema. Endpoint público - não requer autenticação.

**Headers:**
```
Content-Type: application/json
```

**Body da Requisição:**
```json
{
  "nome": "João Silva",
  "email": "joao@exemplo.com", 
  "senha": "minhasenha123"
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "nome": "João Silva",
  "email": "joao@exemplo.com"
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function registerUser(userData) {
  try {
    const response = await fetch('/v1/usuarios/registrar', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(userData)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Erro ao registrar usuário:', error);
    throw error;
  }
}
```

---

### 2. 🔑 **Login de Usuário**

**Endpoint:** `POST /usuarios/login`

**Descrição:** Autentica um usuário e retorna token JWT para acesso às rotas protegidas.

**Headers:**
```
Content-Type: application/json
```

**Body da Requisição:**
```json
{
  "email": "joao@exemplo.com",
  "senha": "minhasenha123"
}
```

**Resposta de Sucesso (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usuario": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "nome": "João Silva",
    "email": "joao@exemplo.com"
  }
}
```

**Exemplo de uso (JavaScript):**
```javascript
async function loginUser(credentials) {
  try {
    const response = await fetch('/v1/usuarios/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(credentials)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Armazenar token para uso posterior
    localStorage.setItem('authToken', data.token);
    localStorage.setItem('user', JSON.stringify(data.usuario));
    
    return data;
  } catch (error) {
    console.error('Erro ao fazer login:', error);
    throw error;
  }
}
```

## 🎨 Componentes React

### Hook Customizado para Autenticação

```javascript
// hooks/useAuth.js
import { useState, useEffect, createContext, useContext } from 'react';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth deve ser usado dentro de AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Verificar se há token armazenado ao inicializar
    const storedToken = localStorage.getItem('authToken');
    const storedUser = localStorage.getItem('user');
    
    if (storedToken && storedUser) {
      setToken(storedToken);
      setUser(JSON.parse(storedUser));
    }
    
    setLoading(false);
  }, []);

  const register = async (userData) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/v1/usuarios/registrar', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Erro ao registrar usuário');
      }
      
      const data = await response.json();
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const login = async (credentials) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/v1/usuarios/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(credentials)
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Credenciais inválidas');
      }
      
      const data = await response.json();
      
      // Armazenar dados de autenticação
      localStorage.setItem('authToken', data.token);
      localStorage.setItem('user', JSON.stringify(data.usuario));
      
      setToken(data.token);
      setUser(data.usuario);
      
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
    setToken(null);
    setUser(null);
    setError(null);
  };

  const isAuthenticated = () => {
    return !!token && !!user;
  };

  const value = {
    user,
    token,
    loading,
    error,
    register,
    login,
    logout,
    isAuthenticated
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
```

### Componente de Login

```javascript
// components/LoginForm.jsx
import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

const LoginForm = ({ onLoginSuccess }) => {
  const { login, loading, error } = useAuth();
  const [formData, setFormData] = useState({
    email: '',
    senha: ''
  });
  const [validationErrors, setValidationErrors] = useState({});

  const validateForm = () => {
    const errors = {};
    
    if (!formData.email) {
      errors.email = 'Email é obrigatório';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      errors.email = 'Email inválido';
    }
    
    if (!formData.senha) {
      errors.senha = 'Senha é obrigatória';
    } else if (formData.senha.length < 6) {
      errors.senha = 'Senha deve ter pelo menos 6 caracteres';
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
    
    setValidationErrors({});
    
    try {
      await login(formData);
      onLoginSuccess?.();
    } catch (error) {
      console.error('Erro no login:', error);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // Limpar erro do campo quando usuário começar a digitar
    if (validationErrors[name]) {
      setValidationErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="login-form">
      <h2>Entrar</h2>
      
      {error && (
        <div className="error-message">
          {error}
        </div>
      )}
      
      <div className="form-group">
        <label htmlFor="email">Email:</label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          className={validationErrors.email ? 'error' : ''}
          disabled={loading}
          required
        />
        {validationErrors.email && (
          <span className="field-error">{validationErrors.email}</span>
        )}
      </div>

      <div className="form-group">
        <label htmlFor="senha">Senha:</label>
        <input
          type="password"
          id="senha"
          name="senha"
          value={formData.senha}
          onChange={handleChange}
          className={validationErrors.senha ? 'error' : ''}
          disabled={loading}
          required
        />
        {validationErrors.senha && (
          <span className="field-error">{validationErrors.senha}</span>
        )}
      </div>

      <button 
        type="submit" 
        disabled={loading}
        className="submit-button"
      >
        {loading ? 'Entrando...' : 'Entrar'}
      </button>
    </form>
  );
};

export default LoginForm;
```

### Componente de Registro

```javascript
// components/RegisterForm.jsx
import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

const RegisterForm = ({ onRegisterSuccess }) => {
  const { register, loading, error } = useAuth();
  const [formData, setFormData] = useState({
    nome: '',
    email: '',
    senha: '',
    confirmarSenha: ''
  });
  const [validationErrors, setValidationErrors] = useState({});

  const validateForm = () => {
    const errors = {};
    
    if (!formData.nome.trim()) {
      errors.nome = 'Nome é obrigatório';
    } else if (formData.nome.length < 2) {
      errors.nome = 'Nome deve ter pelo menos 2 caracteres';
    }
    
    if (!formData.email) {
      errors.email = 'Email é obrigatório';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      errors.email = 'Email inválido';
    }
    
    if (!formData.senha) {
      errors.senha = 'Senha é obrigatória';
    } else if (formData.senha.length < 6) {
      errors.senha = 'Senha deve ter pelo menos 6 caracteres';
    }
    
    if (formData.senha !== formData.confirmarSenha) {
      errors.confirmarSenha = 'Senhas não coincidem';
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
    
    setValidationErrors({});
    
    try {
      const { confirmarSenha, ...userData } = formData;
      await register(userData);
      onRegisterSuccess?.();
    } catch (error) {
      console.error('Erro no registro:', error);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // Limpar erro do campo quando usuário começar a digitar
    if (validationErrors[name]) {
      setValidationErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="register-form">
      <h2>Criar Conta</h2>
      
      {error && (
        <div className="error-message">
          {error}
        </div>
      )}
      
      <div className="form-group">
        <label htmlFor="nome">Nome Completo:</label>
        <input
          type="text"
          id="nome"
          name="nome"
          value={formData.nome}
          onChange={handleChange}
          className={validationErrors.nome ? 'error' : ''}
          disabled={loading}
          required
        />
        {validationErrors.nome && (
          <span className="field-error">{validationErrors.nome}</span>
        )}
      </div>

      <div className="form-group">
        <label htmlFor="email">Email:</label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          className={validationErrors.email ? 'error' : ''}
          disabled={loading}
          required
        />
        {validationErrors.email && (
          <span className="field-error">{validationErrors.email}</span>
        )}
      </div>

      <div className="form-group">
        <label htmlFor="senha">Senha:</label>
        <input
          type="password"
          id="senha"
          name="senha"
          value={formData.senha}
          onChange={handleChange}
          className={validationErrors.senha ? 'error' : ''}
          disabled={loading}
          required
        />
        {validationErrors.senha && (
          <span className="field-error">{validationErrors.senha}</span>
        )}
      </div>

      <div className="form-group">
        <label htmlFor="confirmarSenha">Confirmar Senha:</label>
        <input
          type="password"
          id="confirmarSenha"
          name="confirmarSenha"
          value={formData.confirmarSenha}
          onChange={handleChange}
          className={validationErrors.confirmarSenha ? 'error' : ''}
          disabled={loading}
          required
        />
        {validationErrors.confirmarSenha && (
          <span className="field-error">{validationErrors.confirmarSenha}</span>
        )}
      </div>

      <button 
        type="submit" 
        disabled={loading}
        className="submit-button"
      >
        {loading ? 'Criando...' : 'Criar Conta'}
      </button>
    </form>
  );
};

export default RegisterForm;
```

### Componente de Rota Protegida

```javascript
// components/ProtectedRoute.jsx
import React from 'react';
import { useAuth } from '../hooks/useAuth';

const ProtectedRoute = ({ children, fallback = null }) => {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div className="loading">Verificando autenticação...</div>;
  }

  if (!isAuthenticated()) {
    return fallback || <div className="unauthorized">Acesso negado. Faça login.</div>;
  }

  return children;
};

export default ProtectedRoute;
```

## 🎨 Estilos CSS

```css
/* Auth Forms */
.login-form,
.register-form {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.login-form h2,
.register-form h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #2196F3;
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
}

.form-group input.error {
  border-color: #f44336;
}

.field-error {
  display: block;
  color: #f44336;
  font-size: 12px;
  margin-top: 5px;
}

.submit-button {
  width: 100%;
  padding: 12px;
  background: #2196F3;
  color: white;
  border: none;
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
  margin-bottom: 15px;
  border: 1px solid #ffcdd2;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}

.unauthorized {
  text-align: center;
  padding: 40px;
  color: #f44336;
  background: #ffebee;
  border-radius: 4px;
  margin: 20px;
}

/* Responsive */
@media (max-width: 480px) {
  .login-form,
  .register-form {
    margin: 10px;
    padding: 15px;
  }
}
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| 200 | Sucesso | Login realizado com sucesso |
| 201 | Criado | Usuário registrado com sucesso |
| 400 | Bad Request | Dados inválidos (email malformado, senha muito curta) |
| 401 | Unauthorized | Credenciais inválidas no login |
| 409 | Conflict | Email já existe no registro |
| 422 | Unprocessable Entity | Dados válidos mas não processáveis |
| 500 | Internal Server Error | Erro interno do servidor |

### Exemplos de Respostas de Erro

```json
// Erro 400 - Dados inválidos
{
  "error": "Email deve ter formato válido",
  "details": "O campo email deve conter um endereço válido"
}

// Erro 409 - Email já existe  
{
  "error": "Email já cadastrado",
  "details": "Este email já está sendo usado por outro usuário"
}

// Erro 401 - Credenciais inválidas
{
  "error": "Credenciais inválidas",
  "details": "Email ou senha incorretos"
}
```

## 🔐 Segurança

### 1. **Armazenamento de Token**
- Token JWT armazenado no localStorage para persistência
- Considere usar httpOnly cookies para maior segurança em produção
- Implemente renovação automática de tokens

### 2. **Validação de Entrada**
- Validação no frontend e backend
- Sanitização de dados antes do envio
- Proteção contra XSS e injection

### 3. **Gerenciamento de Sessão**
- Logout automático em caso de token expirado
- Limpeza completa dos dados de autenticação
- Redirecionamento para login quando necessário

### 4. **Headers de Segurança**
```javascript
// Exemplo de interceptador para requests autenticados
const createAuthenticatedRequest = (url, options = {}) => {
  const token = localStorage.getItem('authToken');
  
  return fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers
    }
  });
};
```

## 📱 Considerações para UX

### 1. **Estados de Loading**
- Mostrar indicadores de carregamento durante requisições
- Desabilitar formulários durante submissão
- Feedback visual para usuário

### 2. **Validação em Tempo Real**
- Validar campos conforme usuário digita
- Mostrar senhas strength indicator
- Confirmação visual de campos válidos

### 3. **Mensagens de Erro**
- Mensagens claras e específicas
- Não expor informações sensíveis
- Sugestões de correção quando possível

### 4. **Persistência de Estado**
- Manter usuário logado entre sessões
- Lembrar último email usado
- Redirecionamento inteligente após login

---

Esta documentação fornece tudo necessário para integrar o sistema de autenticação no frontend. Para testes, consulte os arquivos HTTP em `tests/http/auth.http`.