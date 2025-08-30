# Guest API Documentation

## Overview

O módulo Guest gerencia grupos de convidados e confirmações de presença (RSVP) para eventos de casamento. Cada grupo possui uma chave de acesso única que permite aos convidados confirmar sua presença de forma independente.

## Authentication

Todos os endpoints administrativos requerem autenticação JWT via header `Authorization: Bearer <token>`.
Endpoints públicos (RSVP) não requerem autenticação.

---

## Endpoints

### 1. Criar Grupo de Convidados

**POST** `/v1/eventos/{idCasamento}/grupos-de-convidados`

Cria um novo grupo de convidados para um evento.

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Path Parameters:**
- `idCasamento` (string, required): UUID do evento

**Request Body:**
```json
{
  "chaveDeAcesso": "padrinhos123",
  "nomesDosConvidados": [
    "Carlos Silva",
    "Ana Santos"
  ]
}
```

**Response (201 Created):**
```json
{
  "idGrupo": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
```

**Error Responses:**
- `400 Bad Request`: Dados inválidos (chave vazia, sem convidados)
- `401 Unauthorized`: Token JWT inválido
- `500 Internal Server Error`: Erro interno do servidor

---

### 2. Listar Grupos por Evento

**GET** `/v1/eventos/{idEvento}/grupos-de-convidados`

Lista todos os grupos de convidados de um evento com resumo de status.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `idEvento` (string, required): UUID do evento

**Query Parameters:**
- `status` (string, optional): Filtrar por status RSVP (`CONFIRMADO`, `RECUSADO`, `PENDENTE`)

**Response (200 OK):**
```json
{
  "grupos": [
    {
      "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "chaveDeAcesso": "padrinhos123",
      "totalConvidados": 3,
      "convidadosConfirmados": 2,
      "convidadosRecusados": 0,
      "convidadosPendentes": 1
    }
  ],
  "total": 1
}
```

**Error Responses:**
- `401 Unauthorized`: Token JWT inválido
- `400 Bad Request`: ID do evento inválido
- `500 Internal Server Error`: Erro interno do servidor

---

### 3. Obter Grupo por ID (Admin)

**GET** `/v1/grupos-de-convidados/{idGrupo}`

Obtém detalhes completos de um grupo específico.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `idGrupo` (string, required): UUID do grupo

**Response (200 OK):**
```json
{
  "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "idEvento": "b2c3d4e5-f6g7-8901-2345-67890abcdef1",
  "chaveDeAcesso": "padrinhos123",
  "convidados": [
    {
      "id": "c3d4e5f6-g7h8-9012-3456-7890abcdef12",
      "nome": "Carlos Silva",
      "statusRSVP": "CONFIRMADO"
    },
    {
      "id": "d4e5f6g7-h8i9-0123-4567-890abcdef123",
      "nome": "Ana Santos",
      "statusRSVP": "PENDENTE"
    }
  ]
}
```

**Error Responses:**
- `401 Unauthorized`: Token JWT inválido
- `404 Not Found`: Grupo não encontrado
- `400 Bad Request`: ID do grupo inválido
- `500 Internal Server Error`: Erro interno do servidor

---

### 4. Atualizar Grupo de Convidados

**PUT** `/v1/grupos-de-convidados/{idGrupo}`

Atualiza um grupo existente (chave de acesso e lista de convidados).

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Path Parameters:**
- `idGrupo` (string, required): UUID do grupo

**Request Body:**
```json
{
  "chaveDeAcesso": "padrinhos-atualizados",
  "convidados": [
    {
      "id": "c3d4e5f6-g7h8-9012-3456-7890abcdef12",
      "nome": "Carlos Silva Santos"
    },
    {
      "nome": "Pedro Novo Convidado"
    }
  ]
}
```

**Response (204 No Content)**

**Error Responses:**
- `401 Unauthorized`: Token JWT inválido
- `404 Not Found`: Grupo não encontrado
- `400 Bad Request`: Dados inválidos
- `500 Internal Server Error`: Erro interno do servidor

---

### 5. Remover Grupo de Convidados

**DELETE** `/v1/grupos-de-convidados/{idGrupo}`

Remove um grupo de convidados. Só permite remoção se todos os convidados estão com status `PENDENTE`.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `idGrupo` (string, required): UUID do grupo

**Response (204 No Content)**

**Error Responses:**
- `401 Unauthorized`: Token JWT inválido
- `404 Not Found`: Grupo não encontrado
- `400 Bad Request`: Não é possível remover grupo com confirmações de presença
- `500 Internal Server Error`: Erro interno do servidor

---

### 6. Obter Estatísticas de RSVP

**GET** `/v1/eventos/{idEvento}/rsvp-stats`

Retorna estatísticas consolidadas de RSVP para um evento.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `idEvento` (string, required): UUID do evento

**Response (200 OK):**
```json
{
  "totalGrupos": 5,
  "totalConvidados": 24,
  "convidadosConfirmados": 18,
  "convidadosRecusados": 2,
  "convidadosPendentes": 4,
  "percentualConfirmado": 75.0,
  "percentualRecusado": 8.33,
  "percentualPendente": 16.67
}
```

**Error Responses:**
- `401 Unauthorized`: Token JWT inválido
- `400 Bad Request`: ID do evento inválido
- `500 Internal Server Error`: Erro interno do servidor

---

## Endpoints Públicos (RSVP)

### 7. Obter Grupo por Chave de Acesso

**GET** `/v1/acesso-convidado?chave={chave}`

Permite que convidados acessem seu grupo através da chave de acesso para confirmar presença.

**Query Parameters:**
- `chave` (string, required): Chave de acesso do grupo

**Response (200 OK):**
```json
{
  "idGrupo": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "convidados": [
    {
      "id": "c3d4e5f6-g7h8-9012-3456-7890abcdef12",
      "nome": "Carlos Silva",
      "statusRSVP": "PENDENTE"
    },
    {
      "id": "d4e5f6g7-h8i9-0123-4567-890abcdef123",
      "nome": "Ana Santos",
      "statusRSVP": "PENDENTE"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Parâmetro 'chave' obrigatório
- `404 Not Found`: Chave de acesso não encontrada
- `500 Internal Server Error`: Erro interno do servidor

---

### 8. Confirmar Presença (RSVP)

**POST** `/v1/rsvps`

Permite que convidados confirmem ou recusem sua presença.

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "chaveDeAcesso": "padrinhos123",
  "respostas": [
    {
      "idConvidado": "c3d4e5f6-g7h8-9012-3456-7890abcdef12",
      "status": "CONFIRMADO"
    },
    {
      "idConvidado": "d4e5f6g7-h8i9-0123-4567-890abcdef123",
      "status": "RECUSADO"
    }
  ]
}
```

**Response (204 No Content)**

**Error Responses:**
- `400 Bad Request`: Dados inválidos (status inválido, convidado não pertence ao grupo)
- `404 Not Found`: Chave de acesso não encontrada
- `500 Internal Server Error`: Erro interno do servidor

---

## Data Types

### Status RSVP
- `PENDENTE`: Aguardando confirmação
- `CONFIRMADO`: Presença confirmada
- `RECUSADO`: Presença recusada

### Guest Object
```json
{
  "id": "string (UUID)",
  "nome": "string",
  "statusRSVP": "PENDENTE|CONFIRMADO|RECUSADO"
}
```

### Group Summary Object
```json
{
  "id": "string (UUID)",
  "chaveDeAcesso": "string",
  "totalConvidados": "number",
  "convidadosConfirmados": "number", 
  "convidadosRecusados": "number",
  "convidadosPendentes": "number"
}
```

### RSVP Stats Object
```json
{
  "totalGrupos": "number",
  "totalConvidados": "number",
  "convidadosConfirmados": "number",
  "convidadosRecusados": "number", 
  "convidadosPendentes": "number",
  "percentualConfirmado": "number",
  "percentualRecusado": "number",
  "percentualPendente": "number"
}
```

## Error Handling

Todos os endpoints retornam erros no formato:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

### Common Error Codes
- `TOKEN_INVALIDO`: JWT token inválido ou ausente
- `PARAMETRO_INVALIDO`: Parâmetro de URL inválido
- `CORPO_INVALIDO`: JSON do body malformado
- `DADOS_INVALIDOS`: Dados de entrada inválidos
- `NAO_ENCONTRADO`: Recurso não encontrado
- `ERRO_INTERNO`: Erro interno do servidor