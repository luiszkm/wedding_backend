# Estratégia de Testes

Esta documentação descreve a estratégia de testes implementada no Wedding Management API.

## Visão Geral

A estratégia de testes segue a pirâmide de testes, priorizando:

1. **Testes de Unidade** (70%): Lógica de negócio e validações
2. **Testes de Integração** (20%): Interação entre componentes
3. **Testes End-to-End** (10%): Fluxos completos da aplicação

```
      /\
     /  \    E2E Tests
    /____\
   /      \   Integration Tests
  /________\
 /          \  Unit Tests
/____________\
```

---

## Estrutura de Testes

### Localização dos Testes

```
internal/{module}/
├── domain/
│   └── {entity}_test.go     # Testes de domínio
├── application/
│   └── service_test.go      # Testes de aplicação
├── infrastructure/
│   └── repository_test.go   # Testes de repositório
└── interfaces/
    └── rest/
        └── handler_test.go  # Testes de handler
```

### Convenções de Nomenclatura

```go
// Função de teste
func TestFunctionName(t *testing.T) {}

// Subtestes
func TestFunction(t *testing.T) {
    t.Run("should do something when condition", func(t *testing.T) {})
    t.Run("should fail when invalid input", func(t *testing.T) {})
}

// Testes de estrutura
func TestStructMethod(t *testing.T) {}

// Benchmarks
func BenchmarkFunction(b *testing.B) {}
```

---

## Testes de Domínio (Unit Tests)

### Características

- **Rápidos**: Executam em milissegundos
- **Isolados**: Sem dependências externas
- **Determinísticos**: Sempre o mesmo resultado
- **Focados**: Testam uma única funcionalidade

### Exemplo: Testes de Entidade

```go
// internal/guest/domain/group_test.go
package domain

import (
    "testing"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestNewGrupoDeConvidados(t *testing.T) {
    idCasamento := uuid.New()

    t.Run("deve criar grupo com dados válidos", func(t *testing.T) {
        nomes := []string{"João Silva", "Maria Silva"}
        chave := "padrinhos"
        
        grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)

        assert.NoError(t, err)
        assert.NotNil(t, grupo)
        assert.Equal(t, chave, grupo.ChaveDeAcesso())
        assert.Equal(t, idCasamento, grupo.IDCasamento())
        assert.Len(t, grupo.Convidados(), 2)
        assert.Equal(t, "João Silva", grupo.Convidados()[0].Nome())
        assert.Equal(t, "Maria Silva", grupo.Convidados()[1].Nome())
    })

    t.Run("deve retornar erro se chave de acesso for vazia", func(t *testing.T) {
        nomes := []string{"João Silva"}
        
        _, err := NewGrupoDeConvidados(idCasamento, "", nomes)

        assert.Error(t, err)
        assert.Equal(t, ErrChaveDeAcessoObrigatoria, err)
    })

    t.Run("deve retornar erro se não houver convidados", func(t *testing.T) {
        chave := "padrinhos"
        nomes := []string{}

        _, err := NewGrupoDeConvidados(idCasamento, chave, nomes)

        assert.Error(t, err)
        assert.Equal(t, ErrPeloMenosUmConvidado, err)
    })
}

func TestGrupoDeConvidados_ConfirmarPresenca(t *testing.T) {
    grupo := createValidGroup(t)

    t.Run("deve confirmar presença de convidado", func(t *testing.T) {
        err := grupo.ConfirmarPresenca("João Silva", true)

        assert.NoError(t, err)
        convidado := grupo.encontrarConvidado("João Silva")
        assert.Equal(t, StatusConfirmado, convidado.Status())
    })

    t.Run("deve retornar erro para convidado inexistente", func(t *testing.T) {
        err := grupo.ConfirmarPresenca("Pedro Santos", true)

        assert.Error(t, err)
        assert.Equal(t, ErrConvidadoNaoEncontrado, err)
    })
}

// Helper function
func createValidGroup(t *testing.T) *GrupoDeConvidados {
    idCasamento := uuid.New()
    nomes := []string{"João Silva", "Maria Silva"}
    chave := "padrinhos"
    
    grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)
    assert.NoError(t, err)
    
    return grupo
}
```

### Testando Value Objects

```go
func TestEmail(t *testing.T) {
    t.Run("deve criar email válido", func(t *testing.T) {
        email, err := NewEmail("joao@exemplo.com")
        
        assert.NoError(t, err)
        assert.Equal(t, "joao@exemplo.com", email.String())
    })

    t.Run("deve falhar com email inválido", func(t *testing.T) {
        testCases := []string{
            "",
            "invalid-email",
            "@exemplo.com",
            "joao@",
            "joao..silva@exemplo.com",
        }

        for _, tc := range testCases {
            t.Run(fmt.Sprintf("email: %s", tc), func(t *testing.T) {
                _, err := NewEmail(tc)
                assert.Error(t, err)
            })
        }
    })
}
```

---

## Testes de Aplicação (Service Tests)

### Características

- **Testam orquestração**: Coordenação entre domínio e infraestrutura
- **Usam mocks**: Para repositórios e serviços externos
- **Validam casos de uso**: Fluxos completos de negócio

### Criando Mocks

```go
// internal/guest/application/mocks_test.go
package application

import (
    "context"
    "github.com/google/uuid"
    "github.com/luiszkm/wedding_backend/internal/guest/domain"
)

type MockGroupRepository struct {
    groups   map[uuid.UUID]*domain.GrupoDeConvidados
    findByChaveResult *domain.GrupoDeConvidados
    findByChaveError  error
}

func NewMockGroupRepository() *MockGroupRepository {
    return &MockGroupRepository{
        groups: make(map[uuid.UUID]*domain.GrupoDeConvidados),
    }
}

func (m *MockGroupRepository) Save(ctx context.Context, grupo *domain.GrupoDeConvidados) error {
    m.groups[grupo.ID()] = grupo
    return nil
}

func (m *MockGroupRepository) FindByChaveDeAcesso(ctx context.Context, chave string) (*domain.GrupoDeConvidados, error) {
    if m.findByChaveError != nil {
        return nil, m.findByChaveError
    }
    return m.findByChaveResult, nil
}

func (m *MockGroupRepository) SetFindByChaveResult(grupo *domain.GrupoDeConvidados, err error) {
    m.findByChaveResult = grupo
    m.findByChaveError = err
}
```

### Exemplo de Teste de Serviço

```go
// internal/guest/application/service_test.go
package application

import (
    "context"
    "testing"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/luiszkm/wedding_backend/internal/guest/domain"
)

func TestGuestService_CriarGrupoDeConvidados(t *testing.T) {
    t.Run("deve criar grupo com sucesso", func(t *testing.T) {
        repo := NewMockGroupRepository()
        service := NewGuestService(repo)

        req := CreateGroupRequest{
            IDEvento:    uuid.New(),
            ChaveAcesso: "padrinhos",
            Nomes:       []string{"João Silva", "Maria Silva"},
        }

        err := service.CriarGrupoDeConvidados(context.Background(), req)

        assert.NoError(t, err)
        assert.Len(t, repo.groups, 1)
        
        // Verificar se o grupo foi salvo corretamente
        var grupo *domain.GrupoDeConvidados
        for _, g := range repo.groups {
            grupo = g
            break
        }
        
        assert.Equal(t, req.ChaveAcesso, grupo.ChaveDeAcesso())
        assert.Len(t, grupo.Convidados(), 2)
    })

    t.Run("deve falhar com dados inválidos", func(t *testing.T) {
        repo := NewMockGroupRepository()
        service := NewGuestService(repo)

        req := CreateGroupRequest{
            IDEvento:    uuid.New(),
            ChaveAcesso: "", // Inválido
            Nomes:       []string{"João Silva"},
        }

        err := service.CriarGrupoDeConvidados(context.Background(), req)

        assert.Error(t, err)
        assert.Equal(t, domain.ErrChaveDeAcessoObrigatoria, err)
        assert.Len(t, repo.groups, 0)
    })
}

func TestGuestService_ConfirmarPresenca(t *testing.T) {
    t.Run("deve confirmar presença com sucesso", func(t *testing.T) {
        repo := NewMockGroupRepository()
        service := NewGuestService(repo)

        // Setup: criar grupo existente
        grupo := createValidGroup(t)
        repo.SetFindByChaveResult(grupo, nil)

        req := RSVPRequest{
            ChaveAcesso: "padrinhos",
            Confirmacoes: []ConfirmacaoPresenca{
                {Nome: "João Silva", Confirmado: true},
                {Nome: "Maria Silva", Confirmado: false},
            },
        }

        err := service.ConfirmarPresenca(context.Background(), req)

        assert.NoError(t, err)
        
        // Verificar se as confirmações foram processadas
        joao := grupo.encontrarConvidado("João Silva")
        maria := grupo.encontrarConvidado("Maria Silva")
        
        assert.Equal(t, StatusConfirmado, joao.Status())
        assert.Equal(t, StatusRecusado, maria.Status())
    })

    t.Run("deve falhar se grupo não existir", func(t *testing.T) {
        repo := NewMockGroupRepository()
        service := NewGuestService(repo)

        // Setup: grupo não encontrado
        repo.SetFindByChaveResult(nil, domain.ErrGrupoNaoEncontrado)

        req := RSVPRequest{
            ChaveAcesso: "inexistente",
            Confirmacoes: []ConfirmacaoPresenca{
                {Nome: "João Silva", Confirmado: true},
            },
        }

        err := service.ConfirmarPresenca(context.Background(), req)

        assert.Error(t, err)
        assert.Equal(t, domain.ErrGrupoNaoEncontrado, err)
    })
}
```

---

## Testes de Integração

### Características

- **Testam componentes reais**: Database, HTTP, etc.
- **Ambiente controlado**: Database de teste
- **Dados isolados**: Cada teste limpa/prepara seus dados

### Setup para Testes de Database

```go
// internal/guest/infrastructure/repository_test.go
package infrastructure

import (
    "context"
    "testing"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
    // Usar database de teste
    dbURL := "postgres://user:password@localhost:5432/wedding_test_db"
    
    pool, err := pgxpool.New(context.Background(), dbURL)
    require.NoError(t, err)
    
    // Cleanup no final do teste
    t.Cleanup(func() {
        pool.Close()
    })
    
    return pool
}

func cleanupTables(t *testing.T, pool *pgxpool.Pool) {
    tables := []string{
        "convidados",
        "convidados_grupos",
        "eventos",
        "usuarios",
    }
    
    for _, table := range tables {
        _, err := pool.Exec(context.Background(), "DELETE FROM "+table)
        require.NoError(t, err)
    }
}

func TestPostgresGroupRepository_Save(t *testing.T) {
    pool := setupTestDB(t)
    repo := NewPostgresGroupRepository(pool)
    
    t.Run("deve salvar grupo com sucesso", func(t *testing.T) {
        cleanupTables(t, pool)
        
        // Setup: criar evento primeiro
        eventoID := createTestEvento(t, pool)
        
        grupo := createValidGroup(t)
        grupo.SetIDEvento(eventoID) // Método helper ou ajuste na entidade
        
        err := repo.Save(context.Background(), grupo)
        
        assert.NoError(t, err)
        
        // Verificar se foi salvo no banco
        var count int
        err = pool.QueryRow(context.Background(), 
            "SELECT COUNT(*) FROM convidados_grupos WHERE id = $1", 
            grupo.ID()).Scan(&count)
        
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    })
}

func TestPostgresGroupRepository_FindByChaveDeAcesso(t *testing.T) {
    pool := setupTestDB(t)
    repo := NewPostgresGroupRepository(pool)
    
    t.Run("deve encontrar grupo por chave", func(t *testing.T) {
        cleanupTables(t, pool)
        
        // Setup: criar e salvar grupo
        eventoID := createTestEvento(t, pool)
        grupo := createValidGroup(t)
        grupo.SetIDEvento(eventoID)
        
        err := repo.Save(context.Background(), grupo)
        require.NoError(t, err)
        
        // Test: buscar por chave
        found, err := repo.FindByChaveDeAcesso(context.Background(), "padrinhos")
        
        assert.NoError(t, err)
        assert.NotNil(t, found)
        assert.Equal(t, grupo.ID(), found.ID())
        assert.Equal(t, grupo.ChaveDeAcesso(), found.ChaveDeAcesso())
    })
    
    t.Run("deve retornar erro se não encontrar", func(t *testing.T) {
        cleanupTables(t, pool)
        
        _, err := repo.FindByChaveDeAcesso(context.Background(), "inexistente")
        
        assert.Error(t, err)
        assert.Equal(t, domain.ErrGrupoNaoEncontrado, err)
    })
}

// Helper functions
func createTestEvento(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
    eventoID := uuid.New()
    userID := createTestUser(t, pool)
    
    _, err := pool.Exec(context.Background(),
        "INSERT INTO eventos (id, id_usuario, nome, tipo, url_slug) VALUES ($1, $2, $3, $4, $5)",
        eventoID, userID, "Casamento Teste", "CASAMENTO", "casamento-teste")
    require.NoError(t, err)
    
    return eventoID
}

func createTestUser(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
    userID := uuid.New()
    
    _, err := pool.Exec(context.Background(),
        "INSERT INTO usuarios (id, nome, email, password_hash) VALUES ($1, $2, $3, $4)",
        userID, "Teste User", "teste@exemplo.com", "hash")
    require.NoError(t, err)
    
    return userID
}
```

---

## Testes de Handler (HTTP Tests)

### Características

- **Testam camada HTTP**: Request/Response, status codes
- **Isolam dependências**: Mocks para services
- **Validam contratos**: DTOs, headers, etc.

### Exemplo de Teste de Handler

```go
// internal/guest/interfaces/rest/handler_test.go
package rest

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockGuestService struct {
    mock.Mock
}

func (m *MockGuestService) CriarGrupoDeConvidados(ctx context.Context, req application.CreateGroupRequest) error {
    args := m.Called(ctx, req)
    return args.Error(0)
}

func (m *MockGuestService) ConfirmarPresenca(ctx context.Context, req application.RSVPRequest) error {
    args := m.Called(ctx, req)
    return args.Error(0)
}

func TestGuestHandler_HandleCriarGrupoDeConvidados(t *testing.T) {
    t.Run("deve criar grupo com sucesso", func(t *testing.T) {
        service := &MockGuestService{}
        handler := NewGuestHandler(service)

        reqBody := CreateGroupDTO{
            ChaveAcesso: "padrinhos",
            Nomes:       []string{"João Silva", "Maria Silva"},
        }
        
        service.On("CriarGrupoDeConvidados", mock.Anything, mock.MatchedBy(func(req application.CreateGroupRequest) bool {
            return req.ChaveAcesso == "padrinhos" && len(req.Nomes) == 2
        })).Return(nil)

        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequest(http.MethodPost, "/v1/casamentos/123/grupos-de-convidados", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        // Adicionar parâmetro de rota
        rctx := chi.NewRouteContext()
        rctx.URLParams.Add("idCasamento", "123")
        req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

        w := httptest.NewRecorder()

        handler.HandleCriarGrupoDeConvidados(w, req)

        assert.Equal(t, http.StatusCreated, w.Code)
        service.AssertExpectations(t)
    })

    t.Run("deve retornar 400 para JSON inválido", func(t *testing.T) {
        service := &MockGuestService{}
        handler := NewGuestHandler(service)

        req := httptest.NewRequest(http.MethodPost, "/v1/casamentos/123/grupos-de-convidados", 
                                  bytes.NewReader([]byte("invalid json")))
        req.Header.Set("Content-Type", "application/json")
        
        w := httptest.NewRecorder()

        handler.HandleCriarGrupoDeConvidados(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
    })

    t.Run("deve retornar 400 para dados inválidos", func(t *testing.T) {
        service := &MockGuestService{}
        handler := NewGuestHandler(service)

        reqBody := CreateGroupDTO{
            ChaveAcesso: "", // Inválido
            Nomes:       []string{"João Silva"},
        }
        
        service.On("CriarGrupoDeConvidados", mock.Anything, mock.Anything).
                Return(domain.ErrChaveDeAcessoObrigatoria)

        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequest(http.MethodPost, "/v1/casamentos/123/grupos-de-convidados", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        rctx := chi.NewRouteContext()
        rctx.URLParams.Add("idCasamento", "123")
        req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

        w := httptest.NewRecorder()

        handler.HandleCriarGrupoDeConvidados(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        
        var response ErrorResponse
        err := json.NewDecoder(w.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Contains(t, response.Error, "chave de acesso")
    })
}
```

---

## Testes End-to-End

### Características

- **Testam fluxos completos**: Da requisição HTTP ao banco de dados
- **Ambiente real**: Servidor HTTP rodando
- **Dados persistentes**: Database real de teste

### Setup E2E

```go
// test/e2e/setup_test.go
package e2e

import (
    "context"
    "net/http/httptest"
    "testing"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/luiszkm/wedding_backend/cmd/api"
)

type TestSuite struct {
    server *httptest.Server
    db     *pgxpool.Pool
    client *http.Client
}

func setupE2E(t *testing.T) *TestSuite {
    // Database de teste
    dbURL := "postgres://user:password@localhost:5432/wedding_e2e_db"
    pool, err := pgxpool.New(context.Background(), dbURL)
    require.NoError(t, err)
    
    // Setup da aplicação
    app := api.NewApp(pool, testConfig())
    server := httptest.NewServer(app.Handler())
    
    suite := &TestSuite{
        server: server,
        db:     pool,
        client: server.Client(),
    }
    
    t.Cleanup(func() {
        server.Close()
        pool.Close()
    })
    
    return suite
}

func (s *TestSuite) cleanupDB(t *testing.T) {
    tables := []string{
        "convidados", "convidados_grupos", "presentes", 
        "recados", "fotos", "eventos", "usuarios",
    }
    
    for _, table := range tables {
        _, err := s.db.Exec(context.Background(), "DELETE FROM "+table)
        require.NoError(t, err)
    }
}
```

### Exemplo E2E Test

```go
// test/e2e/guest_flow_test.go
package e2e

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGuestFlow(t *testing.T) {
    suite := setupE2E(t)
    suite.cleanupDB(t)

    // 1. Registrar usuário
    user := registerUser(t, suite)
    token := loginUser(t, suite, user.Email, "password123")

    // 2. Criar evento
    evento := createEvento(t, suite, token)

    // 3. Criar grupo de convidados
    grupo := createGrupoConvidados(t, suite, token, evento.ID)

    // 4. Confirmar presença (público)
    confirmPresence(t, suite, grupo.ChaveAcesso)

    // 5. Verificar confirmações (privado)
    verifyConfirmations(t, suite, token, grupo.ID)
}

func registerUser(t *testing.T, suite *TestSuite) User {
    reqBody := map[string]string{
        "nome":  "João Silva",
        "email": "joao@exemplo.com", 
        "senha": "password123",
    }

    body, _ := json.Marshal(reqBody)
    resp, err := suite.client.Post(
        suite.server.URL+"/v1/usuarios/registrar",
        "application/json",
        bytes.NewReader(body),
    )
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    var user User
    err = json.NewDecoder(resp.Body).Decode(&user)
    require.NoError(t, err)

    return user
}

func loginUser(t *testing.T, suite *TestSuite, email, password string) string {
    reqBody := map[string]string{
        "email": email,
        "senha": password,
    }

    body, _ := json.Marshal(reqBody)
    resp, err := suite.client.Post(
        suite.server.URL+"/v1/usuarios/login",
        "application/json", 
        bytes.NewReader(body),
    )
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var loginResp LoginResponse
    err = json.NewDecoder(resp.Body).Decode(&loginResp)
    require.NoError(t, err)

    return loginResp.Token
}

func createEvento(t *testing.T, suite *TestSuite, token string) Evento {
    reqBody := map[string]interface{}{
        "nome":      "Casamento João e Maria",
        "tipo":      "CASAMENTO",
        "url_slug":  "casamento-joao-maria-2024",
    }

    body, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest(http.MethodPost, suite.server.URL+"/v1/eventos", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := suite.client.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    var evento Evento
    err = json.NewDecoder(resp.Body).Decode(&evento)
    require.NoError(t, err)

    return evento
}

func createGrupoConvidados(t *testing.T, suite *TestSuite, token, eventoID string) GrupoConvidados {
    reqBody := map[string]interface{}{
        "chave_acesso": "padrinhos",
        "nomes":        []string{"João Silva", "Maria Silva"},
    }

    body, _ := json.Marshal(reqBody)
    url := fmt.Sprintf("%s/v1/casamentos/%s/grupos-de-convidados", suite.server.URL, eventoID)
    req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := suite.client.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    var grupo GrupoConvidados
    err = json.NewDecoder(resp.Body).Decode(&grupo)
    require.NoError(t, err)

    return grupo
}

func confirmPresence(t *testing.T, suite *TestSuite, chaveAcesso string) {
    reqBody := map[string]interface{}{
        "chave_acesso": chaveAcesso,
        "confirmacoes": []map[string]interface{}{
            {"nome": "João Silva", "confirmado": true},
            {"nome": "Maria Silva", "confirmado": false},
        },
    }

    body, _ := json.Marshal(reqBody)
    resp, err := suite.client.Post(
        suite.server.URL+"/v1/rsvps",
        "application/json",
        bytes.NewReader(body),
    )
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

---

## Executando Testes

### Comandos Básicos

```bash
# Todos os testes
go test ./...

# Testes com verbose
go test -v ./...

# Testes de um módulo específico
go test ./internal/guest/...

# Teste específico
go test -run TestNewGrupoDeConvidados ./internal/guest/domain

# Testes em paralelo
go test -parallel 4 ./...

# Testes com timeout
go test -timeout 30s ./...
```

### Coverage

```bash
# Coverage básico
go test -cover ./...

# Coverage detalhado
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Coverage por módulo
go test -coverprofile=coverage.out ./internal/guest/...
go tool cover -func=coverage.out
```

### Benchmarks

```bash
# Executar benchmarks
go test -bench=. ./...

# Benchmark específico
go test -bench=BenchmarkFunction ./internal/guest/domain

# Benchmark com memory profiling
go test -bench=. -memprofile=mem.prof ./...
```

---

## Ferramentas de Teste

### Bibliotecas Utilizadas

```go
import (
    "testing"                          // Testing padrão do Go
    "github.com/stretchr/testify/assert" // Assertions
    "github.com/stretchr/testify/require" // Requirements (para erro)
    "github.com/stretchr/testify/mock"    // Mocking
    "github.com/stretchr/testify/suite"   // Test suites
)
```

### Configuração de CI/CD

```yaml
# .github/workflows/tests.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_USER: user
          POSTGRES_DB: wedding_test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.23

    - name: Run tests
      run: |
        go test -v -cover ./...
      env:
        DATABASE_URL: postgres://user:password@localhost:5432/wedding_test_db

    - name: Run integration tests
      run: |
        go test -v -tags=integration ./...
      env:
        DATABASE_URL: postgres://user:password@localhost:5432/wedding_test_db
```

---

## Boas Práticas

### 1. Naming

```go
// ✅ Bom
func TestUserService_CreateUser_ShouldReturnErrorWhenEmailExists(t *testing.T) {}

// ❌ Ruim  
func TestCreateUser(t *testing.T) {}
```

### 2. Arrange-Act-Assert

```go
func TestFunction(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := Function(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### 3. Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 4. Test Helpers

```go
// Helpers em arquivos separados
func createTestUser(t *testing.T) *User {
    user, err := domain.NewUser("test@example.com", "Test User")
    require.NoError(t, err)
    return user
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
    // ... setup
    t.Cleanup(func() {
        pool.Close()
    })
    return pool
}
```

### 5. Error Testing

```go
func TestFunction_Error(t *testing.T) {
    _, err := Function("invalid input")
    
    // Verificar que erro ocorreu
    assert.Error(t, err)
    
    // Verificar tipo específico de erro
    assert.Equal(t, domain.ErrInvalidInput, err)
    
    // Verificar se é um tipo de erro
    assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}
```

A estratégia de testes garante qualidade, confiabilidade e facilita refatorações seguras do código.