# API Endpoints

Esta documentação descreve todos os endpoints disponíveis na Wedding Management API.

## Base URL

```
http://localhost:3000/v1
```

## Autenticação

A API utiliza JWT (JSON Web Tokens) para autenticação. Para endpoints protegidos, inclua o token no header:

```
Authorization: Bearer <seu-jwt-token>
```

---

## 🔓 Endpoints Públicos

### Autenticação (IAM)

#### Registrar Usuário
```http
POST /v1/usuarios/registrar
```

**Request Body:**
```json
{
  "nome": "João Silva",
  "email": "joao@exemplo.com",
  "senha": "minhasenha123"
}
```

**Response:**
```json
{
  "id": "uuid-do-usuario",
  "nome": "João Silva",
  "email": "joao@exemplo.com"
}
```

#### Login
```http
POST /v1/usuarios/login
```

**Request Body:**
```json
{
  "email": "joao@exemplo.com",
  "senha": "minhasenha123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usuario": {
    "id": "uuid-do-usuario",
    "nome": "João Silva",
    "email": "joao@exemplo.com"
  }
}
```

### RSVP (Confirmação de Presença)

#### Confirmar Presença
```http
POST /v1/rsvps
```

**Request Body:**
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

### Lista de Presentes Pública

#### Listar Presentes Públicos
```http
GET /v1/casamentos/{idCasamento}/presentes-publico
```

**Response:**
```json
{
  "presentes": [
    {
      "id": "uuid-do-presente",
      "nome": "Jogo de Panelas",
      "descricao": "Conjunto completo de panelas",
      "preco": 299.99,
      "imagem_url": "https://...",
      "selecionado": false
    }
  ]
}
```

### Mural de Recados Público

#### Listar Recados Públicos
```http
GET /v1/casamentos/{idCasamento}/recados/publico
```

**Response:**
```json
{
  "recados": [
    {
      "id": "uuid-do-recado",
      "autor": "João Silva",
      "mensagem": "Parabéns pelo casamento!",
      "aprovado": true,
      "data_criacao": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Páginas Públicas de Eventos

#### Renderizar Página Pública do Evento
```http
GET /v1/eventos/{urlSlug}/pagina
```

Renderiza a página pública de um evento usando seu template configurado.

**Parâmetros:**
- `urlSlug`: URL amigável do evento (ex: "casamento-joao-maria-2024")

**Response:**
- Content-Type: `text/html; charset=utf-8`
- Cache-Control: `public, max-age=300`
- HTML completo renderizado com template do evento

**Exemplo:**
```bash
curl https://api.exemplo.com/v1/eventos/casamento-joao-maria/pagina
```

### Templates

#### Listar Templates Disponíveis
```http
GET /v1/templates/disponiveis
```

Retorna lista de todos os templates padrão disponíveis para seleção.

**Response:**
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
      "criado_em": "2024-01-15T10:00:00Z"
    },
    {
      "id": "template_classico",
      "nome": "Clássico", 
      "descricao": "Template tradicional elegante com tipografia clássica",
      "tipo": "STANDARD",
      "paleta_default": {
        "primary": "#8b5a3c",
        "secondary": "#f5f5dc",
        "accent": "#d4af37",
        "background": "#fdfdf8",
        "text": "#2c1810"
      },
      "suporta_gifts": true,
      "suporta_gallery": true,
      "suporta_messages": true,
      "suporta_rsvp": true,
      "criado_em": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 3
}
```

#### Obter Metadados de Template Específico
```http
GET /v1/templates/{templateId}
```

**Parâmetros:**
- `templateId`: ID do template (ex: "template_moderno" ou "cliente_premium.html")

**Response:**
```json
{
  "id": "template_moderno",
  "nome": "Moderno",
  "descricao": "Template moderno e minimalista",
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
  "criado_em": "2024-01-15T10:00:00Z"
}
```

### Billing

#### Listar Planos
```http
GET /v1/planos
```

**Response:**
```json
{
  "planos": [
    {
      "id": "uuid-do-plano",
      "nome": "Básico",
      "descricao": "Plano básico para casamentos pequenos",
      "preco": 99.99,
      "features": ["100 convidados", "Lista de presentes", "Galeria básica"]
    }
  ]
}
```

#### Webhook Stripe
```http
POST /v1/webhooks/stripe
```
*Endpoint para receber webhooks do Stripe (uso interno)*

---

## 🔒 Endpoints Protegidos

*Todos os endpoints abaixo requerem autenticação JWT*

### Gestão de Convidados

#### Criar Grupo de Convidados
```http
POST /v1/casamentos/{idCasamento}/grupos-de-convidados
```

**Request Body:**
```json
{
  "chave_acesso": "padrinhos",
  "nomes": ["João Silva", "Maria Silva", "Pedro Santos"]
}
```

#### Obter Grupo por Chave de Acesso
```http
GET /v1/acesso-convidado?chave=padrinhos
```

#### Revisar Grupo
```http
PUT /v1/grupos-de-convidados/{idGrupo}
```

**Request Body:**
```json
{
  "nomes": ["João Silva", "Maria Silva", "Ana Santos"]
}
```

### Lista de Presentes

#### Criar Presente
```http
POST /v1/casamentos/{idCasamento}/presentes
```

**Request Body:**
```json
{
  "nome": "Jogo de Panelas",
  "descricao": "Conjunto completo de panelas antiaderentes",
  "preco": 299.99,
  "imagem": "base64-encoded-image"
}
```

#### Finalizar Seleção de Presente
```http
POST /v1/selecoes-de-presente
```

**Request Body:**
```json
{
  "id_presente": "uuid-do-presente",
  "nome_selecionador": "João Silva",
  "email_selecionador": "joao@exemplo.com"
}
```

### Mural de Recados

#### Deixar Recado
```http
POST /v1/recados
```

**Request Body:**
```json
{
  "id_casamento": "uuid-do-casamento",
  "autor": "João Silva",
  "mensagem": "Parabéns pelo casamento! Que sejam muito felizes!"
}
```

#### Listar Recados (Admin)
```http
GET /v1/casamentos/{idCasamento}/recados/admin
```

#### Moderar Recado
```http
PATCH /v1/recados/{idRecado}
```

**Request Body:**
```json
{
  "aprovado": true
}
```

### Galeria de Fotos

#### Fazer Upload de Foto
```http
POST /v1/casamentos/{idCasamento}/fotos
```

**Request Body (multipart/form-data):**
```
file: [arquivo da imagem]
descricao: "Cerimônia no altar"
```

#### Listar Fotos Públicas
```http
GET /v1/casamentos/{idCasamento}/fotos/publico
```

#### Alternar Favorito
```http
POST /v1/fotos/{idFoto}/favoritar
```

#### Adicionar Rótulo
```http
POST /v1/fotos/{idFoto}/rotulos
```

**Request Body:**
```json
{
  "rotulo": "cerimonia"
}
```

#### Remover Rótulo
```http
DELETE /v1/fotos/{idFoto}/rotulos/{nomeDoRotulo}
```

#### Deletar Foto
```http
DELETE /v1/fotos/{idFoto}
```

### Eventos

#### Criar Evento
```http
POST /v1/eventos
```

**Request Body:**
```json
{
  "nome": "Casamento João e Maria",
  "descricao": "Celebração do casamento",
  "data_evento": "2024-06-15T15:00:00Z",
  "local": "Igreja São Pedro",
  "url_slug": "casamento-joao-maria-2024"
}
```

### Billing

#### Criar Assinatura
```http
POST /v1/assinaturas
```

**Request Body:**
```json
{
  "id_plano": "uuid-do-plano",
  "payment_method_id": "pm_stripe_payment_method_id"
}
```

---

## Códigos de Status HTTP

- `200 OK` - Sucesso
- `201 Created` - Recurso criado com sucesso
- `400 Bad Request` - Dados inválidos na requisição
- `401 Unauthorized` - Token de autenticação inválido ou ausente
- `403 Forbidden` - Acesso negado
- `404 Not Found` - Recurso não encontrado
- `409 Conflict` - Conflito (ex: email já existe)
- `422 Unprocessable Entity` - Dados válidos mas não processáveis
- `500 Internal Server Error` - Erro interno do servidor

## Exemplos de Erro

```json
{
  "error": "Descrição do erro",
  "details": "Detalhes adicionais do erro"
}
```