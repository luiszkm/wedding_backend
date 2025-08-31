# HTTP Tests Documentation

Este diretório contém testes HTTP para validação manual e automática da API Wedding Backend.

## Como Executar os Testes

### Pré-requisitos

1. **Servidor Rodando**: O backend deve estar rodando em `localhost:8080`
2. **Banco de Dados**: PostgreSQL deve estar configurado e rodando
3. **Token JWT**: Deve estar válido (atual expira em 2025-01-17)

### Executando o Servidor

```bash
# Opção 1: Direto com Go
go run ./cmd/api/main.go

# Opção 2: Com Docker Compose
docker-compose up --build

# Opção 3: Com binário compilado
go build -o server ./cmd/api/main.go
./server
```

### Executando os Testes

#### Com REST Client (VS Code)

1. Instale a extensão "REST Client" no VS Code
2. Abra o arquivo `guests.http`
3. Execute os requests sequencialmente clicando em "Send Request"

#### Com curl (linha de comando)

```bash
# Testar conectividade
curl http://localhost:8080/health

# Criar evento (substitua o token)
curl -X POST http://localhost:8080/v1/eventos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <seu-token>" \
  -d '{
    "nome": "Teste API",
    "data": "2025-06-15T15:00:00Z",
    "tipo": "CASAMENTO", 
    "urlSlug": "teste-api-2025"
  }'
```

## Arquivos de Teste

### `guests.http`

Testa o módulo completo de convidados com 13 cenários:

1. **Health Check**: Verifica se o servidor está funcionando
2. **Criar Evento**: Cria evento para os testes (captura `eventoId`)
3. **Criar Grupo 1**: Grupo principal (captura `grupoId` e `chaveAcesso`)
4. **Criar Grupo 2**: Segundo grupo para diversidade
5. **Buscar por Chave**: Endpoint público (captura IDs dos convidados)
6. **Confirmar RSVP**: Confirmação de presença pública
7. **Revisar Grupo**: Update administrativo de grupo
8. **Listar Grupos**: Visualização de todos os grupos
9. **Listar com Filtro**: Filtro por status de RSVP
10. **Obter por ID**: Detalhes completos de um grupo
11. **Estatísticas**: Dashboard de métricas RSVP
12. **Criar p/ Remoção**: Grupo separado para teste de DELETE
13. **Remover Válido**: DELETE com validação de negócio
14. **Remover Inválido**: DELETE que deve falhar (teste de validação)

## Variáveis Dinâmicas

Os testes usam JavaScript para capturar e reutilizar dados:

```javascript
// Captura ID do evento criado
> {%
  client.global.set("eventoId", response.body.idEvento);
%}

// Captura IDs dos convidados retornados
> {%
  client.global.set("convidado1Id", response.body.convidados[0].id);
  client.global.set("convidado2Id", response.body.convidados[1].id);
%}
```

## Cenários de Teste

### ✅ Fluxos de Sucesso
- Criação de grupos com dados válidos
- Busca por chave de acesso pública
- Confirmação de RSVP múltiplos convidados
- Listagem com e sem filtros
- Estatísticas agregadas
- Remoção de grupos pendentes

### ❌ Fluxos de Erro
- Remoção de grupos com confirmações (deve falhar com 400)
- Dados inválidos (validação de entrada)
- Recursos não encontrados (404)
- Falta de autenticação (401)

## Códigos de Resposta Esperados

| Endpoint | Método | Sucesso | Erro |
|----------|--------|---------|------|
| `/health` | GET | 200 | - |
| `/eventos` | POST | 201 | 400, 401 |
| `/grupos-de-convidados` | POST | 201 | 400, 401 |
| `/acesso-convidado` | GET | 200 | 404 |
| `/rsvps` | POST | 204 | 400, 404 |
| `/grupos-de-convidados` | PUT | 204 | 400, 401, 404 |
| `/grupos-de-convidados` | GET | 200 | 401 |
| `/rsvp-stats` | GET | 200 | 401 |
| `/grupos-de-convidados/{id}` | DELETE | 204 | 400, 401, 404 |

## Solução de Problemas

### Foreign Key Constraint Error
```
ERROR: insert or update on table "convidados_grupos" violates foreign key constraint
```

**Solução**: Execute primeiro o teste "1. Criar evento" para garantir que o `eventoId` existe.

### Token Expirado
```json
{"error": {"code": "TOKEN_INVALIDO", "message": "..."}}
```

**Solução**: 
1. Faça login na API para obter um novo token
2. Atualize a variável `@authToken` no arquivo
3. Ou use o endpoint de login nos testes

### Servidor Não Responde
```
curl: (7) Failed to connect to localhost:8080
```

**Solução**:
1. Verifique se o servidor está rodando
2. Verifique se as variáveis de ambiente estão configuradas
3. Verifique se o banco PostgreSQL está funcionando