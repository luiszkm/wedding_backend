# 🎁 Documentação de Integração Frontend - Módulo de Lista de Presentes

## Visão Geral

O módulo de Lista de Presentes permite criar e gerenciar uma lista de presentes para o casamento, incluindo suporte a presentes fracionados com sistema de cotas. Oferece interfaces públicas para visualização e seleção de presentes, além de funcionalidades administrativas para gestão completa da lista.

## Endpoints da API

### Base URL

```
http://localhost:8080/v1
```

### 1. 📋 **Listar Presentes Públicos**

**Endpoint:** `GET /eventos/{idEvento}/presentes-publico`

**Descrição:** Retorna lista pública de presentes disponíveis para seleção. Endpoint público - não requer autenticação.

**Parâmetros:**

- `idEvento` (path): UUID do evento

**Resposta de Sucesso (200):**

```json
[
  {
    "id": "6b6b9034-50a9-4dbf-a736-0b4f5d6e8f61",
    "nome": "Cafeteira Elétrica",
    "descricao": "Cafeteira elétrica 12 xícaras",
    "fotoUrl": "",
    "ehFavorito": false,
    "categoria": "COZINHA",
    "detalhes": {
      "tipo": "PRODUTO_EXTERNO",
      "linkDaLoja": "https://exemplo.com/cafeteira"
    },
    "tipo": "INTEGRAL",
    "status": "DISPONIVEL"
  },
  {
    "id": "27c71b89-7c0e-4d16-8d50-c76e8bdb1261",
    "nome": "Geladeira Frost Free",
    "descricao": "Geladeira 2 portas 400L",
    "fotoUrl": "",
    "ehFavorito": true,
    "categoria": "COZINHA",
    "detalhes": {
      "tipo": "PRODUTO_EXTERNO",
      "linkDaLoja": "https://exemplo.com/geladeira"
    },
    "tipo": "FRACIONADO",
    "status": "DISPONIVEL",
    "valorTotal": 2500,
    "valorCota": 250,
    "cotasTotais": 10,
    "cotasDisponiveis": 10,
    "cotasSelecionadas": 0
  }
]
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

**Endpoint:** `POST /eventos/{idEvento}/presentes`

**Descrição:** Cria um novo presente na lista. Requer autenticação JWT.

**Headers:**

```
Content-Type: multipart/form-data
Authorization: Bearer {jwt_token}
```

**Parâmetros:**

- `idEvento` (path): UUID do evento

**Body da Requisição (Multipart Form):**

*Presente Integral:*
```json
{
  "nome": "Cafeteira Elétrica",
  "descricao": "Cafeteira elétrica 12 xícaras",
  "ehFavorito": false,
  "categoria": "COZINHA",
  "tipo": "INTEGRAL",
  "detalhes": {
    "tipo": "PRODUTO_EXTERNO",
    "linkDaLoja": "https://exemplo.com/cafeteira"
  }
}
```

*Presente Fracionado:*
```json
{
  "nome": "Geladeira Frost Free",
  "descricao": "Geladeira 2 portas 400L",
  "ehFavorito": true,
  "categoria": "COZINHA",
  "tipo": "FRACIONADO",
  "valorTotal": 2500.00,
  "numeroCotas": 10,
  "detalhes": {
    "tipo": "PRODUTO_EXTERNO",
    "linkDaLoja": "https://exemplo.com/geladeira"
  }
}
```

**Categorias Disponíveis:**
`MAIN`, `CASAMENTO`, `LUADEMEL`, `HISTORIA`, `FAMILIA`, `OUTROS`, `COZINHA`, `SALA`, `QUARTO`, `BANHEIRO`, `JARDIM`, `DECORACAO`, `ELETRONICOS`, `UTENSILIOS`

**Resposta de Sucesso (201):**

```json
{
  "idPresente": "123e4567-e89b-12d3-a456-426614174000"
}
```

## ⚠️ **Notas de Schema e Migração**

**Hotfixes Aplicados:**
- **Migration 08**: Corrigida coluna `status` removida inadvertidamente pela migration de presentes fracionados
- **Migration 09**: Adicionadas categorias de presente faltantes (`COZINHA`, `SALA`, `QUARTO`, etc.)
- **Migration 10**: Criada tabela `cotas_de_presentes` que estava faltando para suporte a presentes fracionados

**Categorias Completas:**
Após as correções, as seguintes categorias estão disponíveis:
`MAIN`, `CASAMENTO`, `LUADEMEL`, `HISTORIA`, `FAMILIA`, `OUTROS`, `COZINHA`, `SALA`, `QUARTO`, `BANHEIRO`, `JARDIM`, `DECORACAO`, `ELETRONICOS`, `UTENSILIOS`

### 4. 📊 **Listar Presentes Administrativo (Autenticado)**

**Endpoint:** `GET /eventos/{idEvento}/presentes`

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
