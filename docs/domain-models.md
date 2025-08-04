# Modelos de Domínio

Esta documentação descreve os modelos de domínio (entidades) implementados na Wedding Management API, suas responsabilidades e relacionamentos.

---

## 👥 Guest Domain

### GrupoDeConvidados

**Responsabilidades**:
- Representar um grupo de convidados com acesso único
- Controlar confirmações de presença (RSVP)
- Validar regras de criação e modificação

**Atributos**:
```go
type GrupoDeConvidados struct {
    id              uuid.UUID
    idCasamento     uuid.UUID  
    chaveDeAcesso   string
    convidados      []Convidado
    createdAt       time.Time
    updatedAt       time.Time
}
```

**Business Rules**:
- Chave de acesso é obrigatória e única por evento
- Deve ter pelo menos um convidado
- Chave de acesso é case-sensitive

**Métodos Principais**:
```go
func NewGrupoDeConvidados(idCasamento uuid.UUID, chave string, nomes []string) (*GrupoDeConvidados, error)
func (g *GrupoDeConvidados) ConfirmarPresenca(nome string, confirmado bool) error
func (g *GrupoDeConvidados) AdicionarConvidado(nome string) error
func (g *GrupoDeConvidados) RemoverConvidado(nome string) error
```

### Convidado

**Responsabilidades**:
- Representar um convidado individual
- Controlar status de confirmação

**Atributos**:
```go
type Convidado struct {
    id       uuid.UUID
    nome     string
    status   StatusRSVP
}

type StatusRSVP string
const (
    StatusPendente   StatusRSVP = "PENDENTE"
    StatusConfirmado StatusRSVP = "CONFIRMADO" 
    StatusRecusado   StatusRSVP = "RECUSADO"
)
```

**Business Rules**:
- Nome é obrigatório
- Status inicial é sempre `PENDENTE`
- Status pode ser alterado quantas vezes necessário

---

## 🎁 Gift Domain

### Presente

**Responsabilidades**:
- Representar um item da lista de presentes
- Controlar disponibilidade e seleção
- Validar tipos de presente (PIX vs Produto)

**Atributos**:
```go
type Presente struct {
    id          uuid.UUID
    idEvento    uuid.UUID
    nome        string
    descricao   string
    fotoURL     string
    ehFavorito  bool
    status      StatusPresente
    categoria   Categoria
    detalhes    DetalhesPresente
    idSelecao   *uuid.UUID
}

type StatusPresente string
const (
    StatusDisponivel  StatusPresente = "DISPONIVEL"
    StatusSelecionado StatusPresente = "SELECIONADO"
)

type DetalhesPresente struct {
    tipo       TipoDetalhe
    linkLoja   string
    chavePIX   string
}

type TipoDetalhe string
const (
    TipoProdutoExterno TipoDetalhe = "PRODUTO_EXTERNO"
    TipoPIX           TipoDetalhe = "PIX"
)
```

**Business Rules**:
- Nome é obrigatório
- Se tipo é `PRODUTO_EXTERNO`, link da loja é obrigatório
- Se tipo é `PIX`, chave PIX é obrigatória
- Apenas presentes `DISPONIVEL` podem ser selecionados
- Uma vez selecionado, fica indisponível para outros

**Métodos Principais**:
```go
func NewPresente(idEvento uuid.UUID, nome, descricao string, detalhes DetalhesPresente) (*Presente, error)
func (p *Presente) Selecionar(idSelecao uuid.UUID) error
func (p *Presente) Desselecionar() error
func (p *Presente) MarcarComoFavorito() 
func (p *Presente) DesmarcarComoFavorito()
```

### Selecao

**Responsabilidades**:
- Registrar a seleção de um presente por um grupo
- Manter histórico de seleções

**Atributos**:
```go
type Selecao struct {
    id                    uuid.UUID
    idEvento              uuid.UUID
    idGrupoDeConvidados   uuid.UUID
    dataDaSelecao        time.Time
}
```

**Business Rules**:
- Uma seleção por presente
- Vinculada ao grupo que selecionou
- Data é automaticamente definida

---

## 💬 MessageBoard Domain

### Recado

**Responsabilidades**:
- Representar uma mensagem de convidado
- Controlar moderação e aprovação
- Permitir favoritação pelo casal

**Atributos**:
```go
type Recado struct {
    id                    uuid.UUID
    idEvento              uuid.UUID
    idGrupoDeConvidados   uuid.UUID
    nomeDoAutor          string
    texto                string
    status               StatusRecado
    ehFavorito           bool
    createdAt            time.Time
}

type StatusRecado string
const (
    StatusPendente  StatusRecado = "PENDENTE"
    StatusAprovado  StatusRecado = "APROVADO"
    StatusRejeitado StatusRecado = "REJEITADO"
)
```

**Business Rules**:
- Autor (nome) é obrigatório
- Texto não pode ser vazio
- Status inicial é `PENDENTE`
- Apenas recados `APROVADO` aparecem publicamente

**Métodos Principais**:
```go
func NewRecado(idEvento, idGrupo uuid.UUID, autor, texto string) (*Recado, error)
func (r *Recado) Aprovar() error
func (r *Recado) Rejeitar() error
func (r *Recado) MarcarComoFavorito()
func (r *Recado) DesmarcarComoFavorito()
```

---

## 📸 Gallery Domain

### Foto

**Responsabilidades**:
- Representar uma foto da galeria
- Controlar armazenamento e URLs
- Permitir organização por rótulos

**Atributos**:
```go
type Foto struct {
    id          uuid.UUID
    idEvento    uuid.UUID
    storageKey  string
    urlPublica  string
    ehFavorito  bool
    rotulos     []Rotulo
    createdAt   time.Time
}

type Rotulo string
const (
    RotuloPrincipal  Rotulo = "MAIN"
    RotuloCasamento  Rotulo = "CASAMENTO"
    RotuloPré        Rotulo = "PRE_CASAMENTO"
    RotuloPós        Rotulo = "POS_CASAMENTO"
    RotuloPré        Rotulo = "ENSAIO"
    RotuloChamorro   Rotulo = "CHURRASCO"
    // ... outros rótulos
)
```

**Business Rules**:
- Storage key é único e imutável
- URL pública é gerada automaticamente
- Múltiplos rótulos permitidos por foto
- Deleção remove do storage e banco

**Métodos Principais**:
```go
func NewFoto(idEvento uuid.UUID, storageKey, urlPublica string) *Foto
func (f *Foto) AdicionarRotulo(rotulo Rotulo) error
func (f *Foto) RemoverRotulo(rotulo Rotulo) error
func (f *Foto) MarcarComoFavorito()
func (f *Foto) DesmarcarComoFavorito()
```

---

## 👤 IAM Domain

### Usuario

**Responsabilidades**:
- Representar um usuário da plataforma
- Controlar autenticação e autorização
- Gerenciar dados pessoais

**Atributos**:
```go
type Usuario struct {
    id           uuid.UUID
    nome         string
    email        string
    telefone     string
    passwordHash string
    createdAt    time.Time
}
```

**Business Rules**:
- Email é único na plataforma
- Nome é obrigatório
- Senha deve ser hasheada (nunca plain text)
- Telefone é opcional

**Métodos Principais**:
```go
func NewUsuario(nome, email, senha string) (*Usuario, error)
func (u *Usuario) VerificarSenha(senha string) bool
func (u *Usuario) AlterarSenha(novaSenha string) error
func (u *Usuario) AtualizarDados(nome, telefone string) error
```

---

## 📅 Event Domain

### Evento

**Responsabilidades**:
- Representar um evento (casamento, aniversário, etc.)
- Controlar dados do evento e URLs
- Validar tipos e regras específicas

**Atributos**:
```go
type Evento struct {
    id        uuid.UUID
    idUsuario uuid.UUID
    nome      string
    data      time.Time
    tipo      TipoEvento
    urlSlug   string
    createdAt time.Time
}

type TipoEvento string
const (
    TipoCasamento    TipoEvento = "CASAMENTO"
    TipoAniversario  TipoEvento = "ANIVERSARIO" 
    TipoChaDeberbe   TipoEvento = "CHA_DE_BEBE"
    TipoOutro        TipoEvento = "OUTRO"
)
```

**Business Rules**:
- Nome é obrigatório
- URL slug deve ser único globalmente
- Data é opcional (pode ser definida depois)
- Tipo determina funcionalidades disponíveis

**Métodos Principais**:
```go
func NewEvento(idUsuario uuid.UUID, nome string, tipo TipoEvento, urlSlug string) (*Evento, error)
func (e *Evento) DefinirData(data time.Time) error
func (e *Evento) AtualizarNome(nome string) error
func (e *Evento) AtualizarSlug(slug string) error
```

---

## 💳 Billing Domain

### Plano

**Responsabilidades**:
- Representar um plano de assinatura
- Definir limites e características
- Integrar com Stripe

**Atributos**:
```go
type Plano struct {
    id                  uuid.UUID
    nome                string
    precoEmCentavos     int
    numeroMaximoEventos int
    duracaoEmDias       int
    idStripePrice       string
}
```

**Business Rules**:
- Nome deve ser único
- Preço sempre em centavos (para evitar problemas de precisão)
- Limites devem ser positivos
- ID Stripe Price deve ser válido

**Métodos Principais**:
```go
func NewPlano(nome string, preco, maxEventos, duracao int, stripeID string) (*Plano, error)
func (p *Plano) AtualizarPreco(novoPreco int) error
func (p *Plano) AtualizarLimites(maxEventos, duracao int) error
```

### Assinatura

**Responsabilidades**:
- Representar uma assinatura ativa
- Controlar período e status
- Integrar com webhooks Stripe

**Atributos**:
```go
type Assinatura struct {
    id                   uuid.UUID
    idUsuario            uuid.UUID
    idPlano              uuid.UUID
    idStripeSubscription string
    dataInicio           time.Time
    dataFim              time.Time
    status               StatusAssinatura
}

type StatusAssinatura string
const (
    StatusPendente     StatusAssinatura = "PENDENTE"
    StatusAtiva        StatusAssinatura = "ATIVA"
    StatusExpirada     StatusAssinatura = "EXPIRADA"
    StatusCancelada    StatusAssinatura = "CANCELADA"
    StatusInadimplente StatusAssinatura = "INADIMPLENTE"
)
```

**Business Rules**:
- Usuário pode ter apenas uma assinatura ativa
- Data fim é calculada baseada no plano
- Status controla acesso aos recursos
- Integração com Stripe para pagamentos

**Métodos Principais**:
```go
func NewAssinatura(idUsuario, idPlano uuid.UUID, stripeID string) (*Assinatura, error)
func (a *Assinatura) Ativar() error
func (a *Assinatura) Cancelar() error
func (a *Assinatura) Renovar(novaDataFim time.Time) error
func (a *Assinatura) EstaAtiva() bool
func (a *Assinatura) PermiteEventos(quantidade int) bool
```

---

## 🔗 Relacionamentos Entre Domínios

### Hierarquia Principal
```
Usuario (1) ──→ (N) Evento
    ↓               ↓
Assinatura     GrupoDeConvidados (N)
    ↓               ↓
 Plano         Convidado (N)
```

### Funcionalidades por Evento
```
Evento (1) ──→ (N) Presente
    ↓              ↓
    ├─→ (N) Recado ├─→ Selecao
    ├─→ (N) Foto
    └─→ (N) GrupoDeConvidados
```

### Fluxo de Dados
```
1. Usuario registra → cria Assinatura
2. Assinatura ativa → permite criar Evento
3. Evento criado → permite:
   - Criar GrupoDeConvidados
   - Criar Presente
   - Receber Recado
   - Upload Foto
```

---

## 📋 Invariants e Validações

### Invariants Globais
- Todos os IDs são UUIDs v4
- Timestamps usam timezone America/Sao_Paulo
- Strings obrigatórias não podem ser vazias
- Referências devem existir (integridade referencial)

### Validações por Entidade

**GrupoDeConvidados**:
- ChaveDeAcesso: 3-255 caracteres, sem espaços
- Convidados: pelo menos 1, máximo 50
- Nomes: 2-255 caracteres cada

**Presente**:
- Nome: 3-255 caracteres
- Preço: se informado, deve ser positivo
- URL loja: formato URL válido
- Chave PIX: formato válido (email, telefone, CPF/CNPJ)

**Recado**:
- NomeAutor: 2-255 caracteres
- Texto: 10-2000 caracteres
- Sem palavrões ou conteúdo impróprio

**Foto**:
- Arquivo: JPG, PNG, WebP aceitos
- Tamanho: máximo 10MB
- Dimensões: mínimo 200x200px

**Usuario**:
- Email: formato válido, único
- Senha: mínimo 8 caracteres, com letra e número
- Nome: 2-255 caracteres

**Evento**:
- Nome: 3-255 caracteres
- URLSlug: único, formato slug válido
- Data: se informada, não pode ser no passado

---

## 🧪 Testing dos Modelos

### Estratégia de Teste

**Unit Tests** para cada entidade:
- Criação válida
- Validações de business rules
- Métodos de comportamento
- Edge cases

**Exemplo de Estrutura**:
```go
func TestNewGrupoDeConvidados(t *testing.T) {
    t.Run("deve criar com dados válidos", func(t *testing.T) {})
    t.Run("deve falhar com chave vazia", func(t *testing.T) {})
    t.Run("deve falhar sem convidados", func(t *testing.T) {})
}

func TestGrupoDeConvidados_ConfirmarPresenca(t *testing.T) {
    t.Run("deve confirmar convidado existente", func(t *testing.T) {})
    t.Run("deve falhar para convidado inexistente", func(t *testing.T) {})
}
```

### Property-Based Testing

Para validações complexas:
```go
func TestValidateEmail_Properties(t *testing.T) {
    // Propriedade: email válido sempre contém @ e .
    // Propriedade: email inválido sempre retorna erro
    // etc.
}
```

Os modelos de domínio são o coração da aplicação, concentrando toda a lógica de negócio e garantindo consistência através de validações rigorosas.