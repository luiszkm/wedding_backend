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
      "valor_por_cota": null,
      
    },
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "nome": "Lua de Mel - Passagens",
      "descricao": "Contribuição para passagens da lua de mel",
      "preco": 2000.0,
      "imagem_url": "https://storage.example.com/gifts/viagem.jpg",
      "selecionado": false,
      "tipo": "FRACIONADO",
      "cotas_disponiveis": 10,
      "cotas_selecionadas": 3,
      "valor_por_cota": 200.0
    }
  ]
}
```

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
  "valor_total": 400.0,
  "cotas_selecionadas": 2
}
```

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
  "preco": 1500.0,
  "imagem": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQ...",
  "tipo": "FRACIONADO",
  "numero_cotas": 10,
  "valor_por_cota": 150.0
}
```

**Resposta de Sucesso (201):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000"
}
```

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
      "valor_por_cota": 200.0,
      "cotas_selecionadas": 7,
      "cotas_disponiveis": 3,
      "selecoes": [
        {
          "id": "sel-002",
          "nome_selecionador": "João Silva",
          "quantidade_cotas": 3,
          "valor_total": 600.0,
          "data_selecao": "2024-01-16T14:20:00Z"
        }
      ]
    }
  ]
}
```

### Códigos de Status HTTP

| Status | Descrição             | Quando Ocorre                                     |
| ------ | --------------------- | ------------------------------------------------- |
| 200    | Sucesso               | Lista carregada com sucesso                       |
| 201    | Criado                | Presente criado ou seleção realizada              |
| 400    | Bad Request           | Dados inválidos (preço negativo, cotas inválidas) |
| 401    | Unauthorized          | Token JWT inválido ou ausente                     |
| 404    | Not Found             | Presente ou casamento não encontrado              |
| 409    | Conflict              | Presente já selecionado ou cotas insuficientes    |
| 422    | Unprocessable Entity  | Dados válidos mas não processáveis                |
| 500    | Internal Server Error | Erro interno do servidor                          |

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
