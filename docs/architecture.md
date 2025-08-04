# Arquitetura do Sistema

Esta documentação descreve a arquitetura geral do Wedding Management API, baseada em Clean Architecture e Domain-Driven Design.

## Visão Geral da Arquitetura

O sistema segue os princípios da **Clean Architecture** com **Domain-Driven Design (DDD)**, organizando o código em camadas bem definidas que promovem:

- **Separação de responsabilidades**
- **Independência de frameworks externos**
- **Testabilidade**
- **Manutenibilidade**

```
┌─────────────────────────────────────────────────┐
│                  Interfaces                     │
│            (REST Handlers + DTOs)               │
├─────────────────────────────────────────────────┤
│                 Application                     │
│              (Use Cases/Services)               │
├─────────────────────────────────────────────────┤
│                   Domain                        │
│        (Entities + Business Rules)              │
├─────────────────────────────────────────────────┤
│                Infrastructure                   │
│     (Database + External Services)              │
└─────────────────────────────────────────────────┘
```

---

## Estrutura de Diretórios

```
internal/
├── billing/              # Domínio de Billing (Stripe)
├── event/               # Domínio de Eventos
├── gallery/             # Domínio de Galeria
├── gift/                # Domínio de Presentes
├── guest/               # Domínio de Convidados
├── iam/                 # Domínio de Autenticação
├── messageboard/        # Domínio de Recados
└── platform/            # Serviços de Plataforma
    ├── auth/            # JWT e Middleware
    ├── storage/         # Abstração de Storage
    └── web/             # Utilitários Web
```

---

## Camadas da Arquitetura

### 1. Domain Layer (Núcleo)

**Localização**: `internal/{módulo}/domain/`

**Responsabilidades**:
- Entidades de negócio
- Regras de negócio (business rules)
- Interfaces de repositório
- Value Objects
- Domain Services

**Exemplos**:
```go
// internal/guest/domain/group.go
type GrupoDeConvidados struct {
    id              uuid.UUID
    idCasamento     uuid.UUID
    chaveDeAcesso   string
    convidados      []Convidado
    // ...
}

func NewGrupoDeConvidados(idCasamento uuid.UUID, chave string, nomes []string) (*GrupoDeConvidados, error) {
    if chave == "" {
        return nil, ErrChaveDeAcessoObrigatoria
    }
    // Business logic here...
}
```

**Características**:
- **Não depende de nenhuma camada externa**
- Contém a lógica de negócio pura
- Define interfaces que serão implementadas pelas camadas externas

### 2. Application Layer (Casos de Uso)

**Localização**: `internal/{módulo}/application/`

**Responsabilidades**:
- Orquestração de casos de uso
- Coordenação entre domínio e infraestrutura
- Transações de aplicação
- Validações de aplicação

**Exemplos**:
```go
// internal/guest/application/service.go
type GuestService struct {
    repo domain.GroupRepository
}

func (s *GuestService) CriarGrupoDeConvidados(ctx context.Context, req CreateGroupRequest) error {
    // Orchestrate domain operations
    grupo, err := domain.NewGrupoDeConvidados(req.IDEvento, req.ChaveAcesso, req.Nomes)
    if err != nil {
        return err
    }
    
    return s.repo.Save(ctx, grupo)
}
```

### 3. Infrastructure Layer (Infraestrutura)

**Localização**: `internal/{módulo}/infrastructure/`

**Responsabilidades**:
- Implementação de repositórios
- Integração com banco de dados
- Chamadas para APIs externas
- Implementação de interfaces definidas no domínio

**Exemplos**:
```go
// internal/guest/infrastructure/postgres_repository.go
type PostgresGroupRepository struct {
    db *pgxpool.Pool
}

func (r *PostgresGroupRepository) Save(ctx context.Context, grupo *domain.GrupoDeConvidados) error {
    // Database implementation
}
```

### 4. Interfaces Layer (Camada de Apresentação)

**Localização**: `internal/{módulo}/interfaces/rest/`

**Responsabilidades**:
- Handlers HTTP
- DTOs (Data Transfer Objects)
- Validação de entrada
- Serialização/Deserialização

**Exemplos**:
```go
// internal/guest/interfaces/rest/handler.go
type GuestHandler struct {
    service *application.GuestService
}

func (h *GuestHandler) HandleCriarGrupoDeConvidados(w http.ResponseWriter, r *http.Request) {
    // HTTP handling logic
}
```

---

## Padrões Arquiteturais Utilizados

### 1. Dependency Injection

Todas as dependências são injetadas no `main.go` durante a inicialização da aplicação:

```go
// cmd/api/main.go
func main() {
    // Infrastructure
    guestRepo := guestInfra.NewPostgresGroupRepository(dbpool)
    
    // Application Services
    guestService := guestApp.NewGuestService(guestRepo)
    
    // Handlers
    guestHandler := guestREST.NewGuestHandler(guestService)
}
```

### 2. Repository Pattern

Cada domínio define uma interface de repositório implementada na camada de infraestrutura:

```go
// Domain layer defines interface
type GroupRepository interface {
    Save(ctx context.Context, grupo *GrupoDeConvidados) error
    FindByChaveDeAcesso(ctx context.Context, chave string) (*GrupoDeConvidados, error)
}

// Infrastructure layer implements
type PostgresGroupRepository struct { ... }
```

### 3. Service Layer

Serviços de aplicação orquestram operações entre domínio e infraestrutura:

```go
type GuestService struct {
    repo domain.GroupRepository
}
```

---

## Módulos de Domínio

### Guest (Convidados)
- **Entidades**: GrupoDeConvidados, Convidado
- **Funcionalidades**: Criação de grupos, RSVP, gestão de convidados
- **Business Rules**: Validação de chaves de acesso, controle de confirmações

### Gift (Presentes)
- **Entidades**: Presente, Selecao
- **Funcionalidades**: Lista de presentes, seleção por convidados
- **Business Rules**: Controle de disponibilidade, validação de tipos (PIX/Produto)

### MessageBoard (Recados)
- **Entidades**: Recado
- **Funcionalidades**: Criação de recados, moderação
- **Business Rules**: Aprovação de recados, controle de favoritos

### Gallery (Galeria)
- **Entidades**: Foto
- **Funcionalidades**: Upload de fotos, organização por rótulos
- **Business Rules**: Gestão de favoritos, categorização

### IAM (Identity & Access Management)
- **Entidades**: Usuario
- **Funcionalidades**: Registro, login, autenticação JWT
- **Business Rules**: Validação de credenciais, gestão de tokens

### Event (Eventos)
- **Entidades**: Evento
- **Funcionalidades**: Criação de eventos, gestão de URLs
- **Business Rules**: Validação de slugs únicos, tipos de evento

### Billing (Cobrança)
- **Entidades**: Plano, Assinatura
- **Funcionalidades**: Integração Stripe, gestão de assinaturas
- **Business Rules**: Controle de limites por plano, webhook handling

---

## Serviços de Plataforma

### Auth (`internal/platform/auth/`)
- **JWT Service**: Geração e validação de tokens
- **Middleware**: Autenticação de rotas protegidas
- **Context Management**: Injeção de user ID no contexto

### Storage (`internal/platform/storage/`)
- **Interface**: Abstração para serviços de storage
- **R2 Implementation**: Implementação para Cloudflare R2/AWS S3
- **File Management**: Upload, delete, URL generation

### Web (`internal/platform/web/`)
- **Common Utilities**: Utilitários web compartilhados
- **Error Handling**: Padronização de respostas de erro
- **Response Formatting**: Formatação de respostas JSON

---

## Fluxo de Dados

### Request Flow (Inbound)
```
HTTP Request → Handler → DTO → Application Service → Domain → Repository → Database
```

### Response Flow (Outbound)  
```
Database → Repository → Domain → Application Service → DTO → Handler → HTTP Response
```

### Exemplo Completo
```
POST /v1/grupos-convidados
   ↓
GuestHandler.HandleCriarGrupoDeConvidados()
   ↓
GuestService.CriarGrupoDeConvidados()
   ↓
domain.NewGrupoDeConvidados() (business rules)
   ↓
PostgresGroupRepository.Save()
   ↓
PostgreSQL Database
```

---

## Princípios de Design

### 1. Dependency Rule
- Dependências apontam sempre para dentro (em direção ao domínio)
- Camadas externas conhecem camadas internas, mas não o contrário

### 2. Single Responsibility
- Cada camada tem uma responsabilidade bem definida
- Separação clara entre lógica de negócio e infraestrutura

### 3. Interface Segregation
- Interfaces pequenas e específicas
- Repositórios definem apenas métodos necessários

### 4. Dependency Inversion
- Domínio define interfaces, infraestrutura implementa
- Facilita testing e substituição de implementações

---

## Testing Strategy

### Domain Layer
- **Unit Tests**: Testam business rules e validações
- **Exemplo**: `internal/guest/domain/group_test.go`
- **Foco**: Lógica de negócio pura, sem dependências externas

### Application Layer
- **Integration Tests**: Testam coordenação entre camadas
- **Mocks**: Para repositórios e serviços externos

### Infrastructure Layer
- **Database Tests**: Testam queries e persistência
- **External Service Tests**: Testam integrações (Stripe, S3)

### Interfaces Layer
- **HTTP Tests**: Testam handlers e DTOs
- **End-to-End Tests**: Testam fluxo completo

---

## Configuração e Bootstrapping

### main.go Structure
```go
func main() {
    // 1. Environment & Configuration
    loadEnvironment()
    
    // 2. External Dependencies
    dbpool := initDatabase()
    storageService := initStorage()
    
    // 3. Platform Services
    jwtService := auth.NewJWTService(secret)
    
    // 4. Repositories
    repos := initRepositories(dbpool)
    
    // 5. Application Services
    services := initServices(repos)
    
    // 6. Handlers
    handlers := initHandlers(services)
    
    // 7. Router & Server
    router := setupRouter(handlers, jwtService)
    startServer(router)
}
```

Esta arquitetura promove código limpo, testável e de fácil manutenção, seguindo as melhores práticas de desenvolvimento em Go.