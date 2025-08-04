# Guia de Desenvolvimento

Este guia fornece informações essenciais para desenvolvedores trabalhando no Wedding Management API.

## Configuração do Ambiente de Desenvolvimento

### Pré-requisitos

```bash
# Go 1.23+
go version

# Docker e Docker Compose
docker --version
docker-compose --version

# Git
git --version

# Opcional: Air para hot reload
go install github.com/cosmtrek/air@latest
```

### Setup Inicial

```bash
# 1. Clone o repositório
git clone <repository-url>
cd wedding_backend

# 2. Configure variáveis de ambiente
cp .env.example .env
# Edite o .env com suas configurações

# 3. Inicie o banco de dados
docker-compose up db

# 4. Execute as migrações
psql -h localhost -U user -d wedding_db -f db/init/01-init.sql
psql -h localhost -U user -d wedding_db -f db/init/02-seed-plans.sql

# 5. Execute a aplicação
go run ./cmd/api/main.go
```

---

## Estrutura do Projeto

### Organização de Módulos

Cada módulo de domínio segue a mesma estrutura:

```
internal/{module}/
├── application/
│   └── service.go          # Use cases e orquestração
├── domain/
│   ├── {entity}.go         # Entidades de negócio
│   ├── repository.go       # Interfaces de repositório
│   └── {entity}_test.go    # Testes de domínio
├── infrastructure/
│   └── postgres_repository.go  # Implementação PostgreSQL
└── interfaces/
    └── rest/
        ├── dto.go          # Data Transfer Objects
        └── handler.go      # HTTP Handlers
```

### Convenções de Nomenclatura

**Arquivos:**
- `snake_case` para nomes de arquivos
- `_test.go` sufixo para testes
- `postgres_` prefixo para implementações PostgreSQL

**Go Code:**
- `PascalCase` para tipos exportados
- `camelCase` para variáveis e funções privadas
- `UPPER_CASE` para constantes

**Database:**
- `snake_case` para tabelas e colunas
- `plural` para nomes de tabelas
- `id` prefixo para chaves estrangeiras

---

## Desenvolvimento de Features

### 1. Criando Nova Entidade de Domínio

```go
// internal/newmodule/domain/entity.go
package domain

import (
    "errors"
    "github.com/google/uuid"
)

var (
    ErrEntityNotFound = errors.New("entity not found")
    ErrInvalidData    = errors.New("invalid data")
)

type Entity struct {
    id    uuid.UUID
    name  string
    // ... outros campos
}

func NewEntity(name string) (*Entity, error) {
    if name == "" {
        return nil, ErrInvalidData
    }
    
    return &Entity{
        id:   uuid.New(),
        name: name,
    }, nil
}

// Getters
func (e *Entity) ID() uuid.UUID { return e.id }
func (e *Entity) Name() string  { return e.name }
```

### 2. Definindo Interface de Repositório

```go
// internal/newmodule/domain/repository.go
package domain

import "context"

type EntityRepository interface {
    Save(ctx context.Context, entity *Entity) error
    FindByID(ctx context.Context, id uuid.UUID) (*Entity, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 3. Implementando Repositório

```go
// internal/newmodule/infrastructure/postgres_repository.go
package infrastructure

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/luiszkm/wedding_backend/internal/newmodule/domain"
)

type PostgresEntityRepository struct {
    db *pgxpool.Pool
}

func NewPostgresEntityRepository(db *pgxpool.Pool) *PostgresEntityRepository {
    return &PostgresEntityRepository{db: db}
}

func (r *PostgresEntityRepository) Save(ctx context.Context, entity *domain.Entity) error {
    query := `INSERT INTO entities (id, name) VALUES ($1, $2)`
    _, err := r.db.Exec(ctx, query, entity.ID(), entity.Name())
    return err
}
```

### 4. Criando Serviço de Aplicação

```go
// internal/newmodule/application/service.go
package application

import (
    "context"
    "github.com/luiszkm/wedding_backend/internal/newmodule/domain"
)

type EntityService struct {
    repo domain.EntityRepository
}

func NewEntityService(repo domain.EntityRepository) *EntityService {
    return &EntityService{repo: repo}
}

func (s *EntityService) CreateEntity(ctx context.Context, name string) error {
    entity, err := domain.NewEntity(name)
    if err != nil {
        return err
    }
    
    return s.repo.Save(ctx, entity)
}
```

### 5. Implementando Handler HTTP

```go
// internal/newmodule/interfaces/rest/handler.go
package rest

import (
    "encoding/json"
    "net/http"
    "github.com/luiszkm/wedding_backend/internal/newmodule/application"
)

type EntityHandler struct {
    service *application.EntityService
}

func NewEntityHandler(service *application.EntityService) *EntityHandler {
    return &EntityHandler{service: service}
}

func (h *EntityHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
    var req CreateEntityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := h.service.CreateEntity(r.Context(), req.Name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
}
```

---

## Padrões de Desenvolvimento

### Error Handling

```go
// Erros de domínio como variáveis
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidEmail     = errors.New("invalid email")
    ErrDuplicateEmail   = errors.New("email already exists")
)

// Wrapping erros para contexto
func (s *UserService) CreateUser(ctx context.Context, email string) error {
    if err := s.validateEmail(email); err != nil {
        return fmt.Errorf("creating user: %w", err)
    }
    // ...
}

// Tratamento em handlers
func (h *UserHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
    err := h.service.CreateUser(r.Context(), req.Email)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrInvalidEmail):
            http.Error(w, err.Error(), http.StatusBadRequest)
        case errors.Is(err, domain.ErrDuplicateEmail):
            http.Error(w, err.Error(), http.StatusConflict)
        default:
            http.Error(w, "internal error", http.StatusInternalServerError)
        }
        return
    }
}
```

### Context Usage

```go
// Sempre passe context como primeiro parâmetro
func (s *Service) ProcessData(ctx context.Context, data Data) error {
    // Verifique cancelamento
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Passe context para chamadas downstream
    return s.repo.Save(ctx, data)
}
```

### Validation Patterns

```go
// Validação no domínio
func NewUser(email, name string) (*User, error) {
    if email == "" {
        return nil, ErrEmailRequired
    }
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    if name == "" {
        return nil, ErrNameRequired
    }
    
    return &User{
        id:    uuid.New(),
        email: email,
        name:  name,
    }, nil
}

// Validação adicional na aplicação
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    // Validação de unicidade
    exists, err := s.repo.ExistsByEmail(ctx, req.Email)
    if err != nil {
        return err
    }
    if exists {
        return ErrDuplicateEmail
    }
    
    user, err := domain.NewUser(req.Email, req.Name)
    if err != nil {
        return err
    }
    
    return s.repo.Save(ctx, user)
}
```

---

## Testing

### Testes de Domínio

```go
// internal/guest/domain/group_test.go
func TestNewGrupoDeConvidados(t *testing.T) {
    t.Run("deve criar grupo válido", func(t *testing.T) {
        idCasamento := uuid.New()
        nomes := []string{"João", "Maria"}
        chave := "padrinhos"
        
        grupo, err := NewGrupoDeConvidados(idCasamento, chave, nomes)
        
        assert.NoError(t, err)
        assert.NotNil(t, grupo)
        assert.Equal(t, chave, grupo.ChaveDeAcesso())
        assert.Len(t, grupo.Convidados(), 2)
    })
    
    t.Run("deve falhar com chave vazia", func(t *testing.T) {
        idCasamento := uuid.New()
        nomes := []string{"João"}
        
        _, err := NewGrupoDeConvidados(idCasamento, "", nomes)
        
        assert.Error(t, err)
        assert.Equal(t, ErrChaveDeAcessoObrigatoria, err)
    })
}
```

### Testes de Serviço com Mocks

```go
// Criando mock do repositório
type MockGroupRepository struct {
    groups map[uuid.UUID]*domain.GrupoDeConvidados
}

func (m *MockGroupRepository) Save(ctx context.Context, grupo *domain.GrupoDeConvidados) error {
    m.groups[grupo.ID()] = grupo
    return nil
}

// Teste do serviço
func TestGuestService_CriarGrupo(t *testing.T) {
    repo := &MockGroupRepository{groups: make(map[uuid.UUID]*domain.GrupoDeConvidados)}
    service := NewGuestService(repo)
    
    req := CreateGroupRequest{
        IDEvento:     uuid.New(),
        ChaveAcesso:  "padrinhos",
        Nomes:        []string{"João", "Maria"},
    }
    
    err := service.CriarGrupoDeConvidados(context.Background(), req)
    
    assert.NoError(t, err)
    assert.Len(t, repo.groups, 1)
}
```

### Executando Testes

```bash
# Todos os testes
go test ./...

# Testes com verbose
go test -v ./...

# Testes de um módulo específico
go test ./internal/guest/domain

# Testes com coverage
go test -cover ./...

# Testes específicos
go test -run TestNewGrupoDeConvidados ./internal/guest/domain
```

---

## Debugging

### Usando Delve

```bash
# Instalar delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug da aplicação
dlv debug ./cmd/api/main.go

# Debug de testes
dlv test ./internal/guest/domain
```

### Logs para Debug

```go
import "log"

// Log simples
log.Printf("Processing user: %s", userID)

// Log com contexto
func (s *Service) ProcessUser(ctx context.Context, userID string) error {
    log.Printf("Starting process for user: %s", userID)
    defer log.Printf("Finished process for user: %s", userID)
    
    // ... lógica
}
```

---

## Hot Reload com Air

### Instalar Air

```bash
go install github.com/cosmtrek/air@latest
```

### Configurar .air.toml

```toml
# .air.toml
root = "."
cmd = "go run ./cmd/api/main.go"
bin = "tmp/main"

[build]
  args_bin = []
  bin = "tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "docs"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

### Executar com Hot Reload

```bash
# Na raiz do projeto
air
```

---

## Database Migrations

### Criar Nova Migration

```sql
-- db/migrations/004-add-new-table.sql
CREATE TABLE new_table (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('America/Sao_Paulo', now())
);
```

### Aplicar Migration

```bash
# Aplicar manualmente
psql -h localhost -U user -d wedding_db -f db/migrations/004-add-new-table.sql

# Ou usar ferramenta de migration (futuro)
migrate -path db/migrations -database "postgres://user:password@localhost:5432/wedding_db?sslmode=disable" up
```

---

## Performance Tips

### Database Queries

```go
// Use prepared statements para queries repetidas
const getUserQuery = "SELECT id, name, email FROM users WHERE id = $1"

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    var user domain.User
    err := r.db.QueryRow(ctx, getUserQuery, id).Scan(&user.id, &user.name, &user.email)
    return &user, err
}

// Use batch operations quando possível
func (r *PostgresUserRepository) SaveMany(ctx context.Context, users []*domain.User) error {
    batch := &pgx.Batch{}
    for _, user := range users {
        batch.Queue("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", 
                   user.ID(), user.Name(), user.Email())
    }
    
    results := r.db.SendBatch(ctx, batch)
    defer results.Close()
    
    for range users {
        _, err := results.Exec()
        if err != nil {
            return err
        }
    }
    return nil
}
```

### Memory Management

```go
// Use context para cancelamento
func (s *Service) ProcessLargeDataset(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Process chunk
        }
    }
}

// Limite tamanho de slices grandes
func (s *Service) GetUsers(ctx context.Context, limit int) ([]*User, error) {
    if limit > 1000 {
        limit = 1000 // Máximo razoável
    }
    return s.repo.FindMany(ctx, limit)
}
```

---

## Troubleshooting Comum

### Erro de Importação Circular

```bash
# Erro: import cycle not allowed
# Solução: Reorganizar dependências, usar interfaces
```

### Deadlock no Banco

```go
// Sempre use context com timeout para queries
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := repo.Save(ctx, data)
```

### Memory Leaks

```go
// Sempre feche recursos
func (r *Repository) Query(ctx context.Context) error {
    rows, err := r.db.Query(ctx, "SELECT ...")
    if err != nil {
        return err
    }
    defer rows.Close() // Importante!
    
    // ... processar rows
}
```

---

## Code Review Checklist

### Antes de Commit

- [ ] Testes passando: `go test ./...`
- [ ] Código formatado: `go fmt ./...`
- [ ] Análise estática: `go vet ./...`
- [ ] Não há logs de debug esquecidos
- [ ] Documentação atualizada se necessário

### Durante Review

- [ ] Lógica de negócio no domínio
- [ ] Validações adequadas
- [ ] Error handling apropriado
- [ ] Context passado corretamente
- [ ] Testes cobrem casos importantes
- [ ] Performance considerada
- [ ] Segurança verificada