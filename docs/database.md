# Banco de Dados

Esta documentação descreve o schema do banco de dados PostgreSQL utilizado pela Wedding Management API.

## Visão Geral

O banco de dados é estruturado para suportar múltiplos usuários, cada um podendo ter múltiplos eventos (casamentos), com funcionalidades completas de gestão de convidados, presentes, recados e galeria de fotos.

## Tipos ENUM

O sistema utiliza vários tipos ENUM para garantir consistência de dados:

```sql
-- Status de confirmação de presença
CREATE TYPE status_rsvp AS ENUM ('PENDENTE', 'CONFIRMADO', 'RECUSADO');

-- Tipos de detalhes para presentes
CREATE TYPE tipo_detalhe_presente AS ENUM ('PRODUTO_EXTERNO', 'PIX');

-- Status dos presentes
CREATE TYPE status_presente AS ENUM ('DISPONIVEL', 'SELECIONADO');

-- Status dos recados
CREATE TYPE status_recado AS ENUM ('PENDENTE', 'APROVADO', 'REJEITADO');

-- Rótulos para fotos
CREATE TYPE nome_rotulo_enum AS ENUM ('MAIN', 'CASAMENTO', 'LUADEMEL', 'HISTORIA', 'FAMILIA', 'OUTROS');

-- Status das assinaturas
CREATE TYPE status_assinatura AS ENUM ('PENDENTE', 'ATIVA', 'EXPIRADA', 'CANCELADA');

-- Tipos de eventos
CREATE TYPE tipo_evento AS ENUM ('CASAMENTO', 'ANIVERSARIO', 'CHA_DE_BEBE', 'OUTRO');
```

---

## Tabelas de Plataforma

### usuarios
Tabela principal de usuários da plataforma.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| nome | VARCHAR(255) | Nome completo do usuário |
| email | VARCHAR(255) | Email único do usuário |
| telefone | VARCHAR(20) | Telefone (opcional) |
| password_hash | TEXT | Hash da senha |
| created_at | TIMESTAMP | Data de criação |

### planos
Define os planos de assinatura disponíveis.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| nome | VARCHAR(100) | Nome do plano (único) |
| preco_em_centavos | INTEGER | Preço em centavos |
| numero_maximo_eventos | INTEGER | Máximo de eventos permitidos |
| duracao_em_dias | INTEGER | Duração do plano em dias |
| id_stripe_price | VARCHAR(255) | ID do preço no Stripe |

**Planos Padrão:**
- **Mensal**: R$ 99,90 - 1 evento - 30 dias
- **Trimestral**: R$ 279,90 - 3 eventos - 90 dias  
- **Semestral**: R$ 539,90 - 5 eventos - 180 dias

### assinaturas
Controla as assinaturas ativas dos usuários.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_usuario | UUID | FK para usuarios |
| id_plano | UUID | FK para planos |
| data_inicio | TIMESTAMP | Data de início da assinatura |
| data_fim | TIMESTAMP | Data de fim da assinatura |
| status | status_assinatura | Status atual da assinatura |

### eventos
Representa os eventos (casamentos) criados pelos usuários.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_usuario | UUID | FK para usuarios |
| nome | VARCHAR(255) | Nome do evento |
| data | DATE | Data do evento |
| tipo | tipo_evento | Tipo do evento |
| url_slug | VARCHAR(255) | Slug único para URLs |
| created_at | TIMESTAMP | Data de criação |

---

## Tabelas de Evento

### convidados_grupos
Grupos de convidados com chave de acesso única.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_evento | UUID | FK para eventos |
| chave_de_acesso | VARCHAR(255) | Chave única para RSVP |
| created_at | TIMESTAMP | Data de criação |
| updated_at | TIMESTAMP | Data de atualização |

**Constraint:** `UNIQUE(id_evento, chave_de_acesso)`

### convidados
Convidados individuais dentro de grupos.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_grupo | UUID | FK para convidados_grupos |
| nome | VARCHAR(255) | Nome do convidado |
| status_rsvp | status_rsvp | Status de confirmação |

### presentes_selecoes
Registra as seleções de presentes pelos convidados.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_evento | UUID | FK para eventos |
| id_grupo_de_convidados | UUID | FK para convidados_grupos |
| data_da_selecao | TIMESTAMP | Data da seleção |

### presentes
Lista de presentes do evento, suportando presentes integrais e fracionados.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_evento | UUID | FK para eventos |
| nome | VARCHAR(255) | Nome do presente |
| descricao | TEXT | Descrição detalhada |
| foto_url | TEXT | URL da foto do presente |
| eh_favorito | BOOLEAN | Se é favorito do casal |
| status | status_presente | Status atual (DISPONIVEL, PARCIALMENTE_SELECIONADO, SELECIONADO) |
| categoria | nome_rotulo_enum | Categoria do presente |
| detalhes_tipo | tipo_detalhe_presente | Tipo de detalhe |
| detalhes_link_loja | TEXT | Link para loja externa |
| detalhes_chave_pix | VARCHAR(255) | Chave PIX para doação |
| id_selecao | UUID | FK para presentes_selecoes |
| tipo | tipo_presente | Tipo do presente (INTEGRAL, FRACIONADO) |
| valor_total_presente | NUMERIC(10,2) | Valor total (para presentes fracionados) |

**Business Rules:**
- Se `detalhes_tipo = 'PRODUTO_EXTERNO'`, então `detalhes_link_loja` deve estar preenchido
- Se `detalhes_tipo = 'PIX'`, então `detalhes_chave_pix` deve estar preenchida
- Se `tipo = 'FRACIONADO'`, então `valor_total_presente` deve estar preenchido
- Presentes fracionados têm cotas associadas na tabela `cotas_de_presentes`

### cotas_de_presentes
Sistema de cotas para presentes fracionados.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_presente | UUID | FK para presentes |
| numero_cota | INTEGER | Número sequencial da cota (1, 2, 3...) |
| valor_cota | NUMERIC(10,2) | Valor individual da cota |
| status | status_presente | Status da cota (DISPONIVEL, SELECIONADO) |
| id_selecao | UUID | FK para presentes_selecoes (quando selecionada) |

**Business Rules:**
- `numero_cota` deve ser positivo e único por presente
- `valor_cota` deve ser positivo
- Para um presente fracionado: valor_total = SUM(valor_cota) de todas as cotas
- Cotas só podem ser selecionadas individualmente
- Constraint única: `(id_presente, numero_cota)`

### recados
Mural de recados dos convidados.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_evento | UUID | FK para eventos |
| id_grupo_de_convidados | UUID | FK para convidados_grupos |
| nome_do_autor | VARCHAR(255) | Nome do autor |
| texto | TEXT | Conteúdo do recado |
| status | status_recado | Status de moderação |
| eh_favorito | BOOLEAN | Se é favorito do casal |
| created_at | TIMESTAMP | Data de criação |

### fotos
Galeria de fotos do evento.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | Chave primária |
| id_evento | UUID | FK para eventos |
| storage_key | TEXT | Chave no storage (S3/R2) |
| url_publica | TEXT | URL pública da foto |
| eh_favorito | BOOLEAN | Se é favorita do casal |
| created_at | TIMESTAMP | Data de upload |

### fotos_rotulos
Rótulos das fotos (relação N:N).

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id_foto | UUID | FK para fotos |
| nome_rotulo | nome_rotulo_enum | Rótulo da foto |

**Chave Primária:** `(id_foto, nome_rotulo)`

---

## Relacionamentos

### Hierarquia Principal
```
usuarios (1) → (N) eventos
eventos (1) → (N) convidados_grupos
convidados_grupos (1) → (N) convidados
```

### Funcionalidades por Evento
```
eventos (1) → (N) presentes
eventos (1) → (N) recados  
eventos (1) → (N) fotos
eventos (1) → (N) presentes_selecoes
```

### Sistema de Assinaturas
```
usuarios (1) → (N) assinaturas
planos (1) → (N) assinaturas
```

---

## Índices Recomendados

Para otimizar performance, considere criar os seguintes índices:

```sql
-- Índices para queries frequentes
CREATE INDEX idx_eventos_usuario ON eventos(id_usuario);
CREATE INDEX idx_convidados_grupos_evento ON convidados_grupos(id_evento);
CREATE INDEX idx_convidados_grupo ON convidados(id_grupo);
CREATE INDEX idx_presentes_evento ON presentes(id_evento);
CREATE INDEX idx_recados_evento ON recados(id_evento);
CREATE INDEX idx_fotos_evento ON fotos(id_evento);

-- Índices para chaves de acesso
CREATE INDEX idx_convidados_grupos_chave ON convidados_grupos(chave_de_acesso);
CREATE INDEX idx_eventos_slug ON eventos(url_slug);

-- Índices para status
CREATE INDEX idx_recados_status ON recados(status);
CREATE INDEX idx_presentes_status ON presentes(status);
CREATE INDEX idx_assinaturas_status ON assinaturas(status);
```

---

## Migração e Seed

### Arquivos de Inicialização
1. **01-init.sql**: Schema completo das tabelas
2. **02-seed-plans.sql**: Planos padrão de assinatura
3. **03-alter-subscriptions.sql**: Alterações de schema (se houver)

### Docker Compose
Os arquivos SQL são executados automaticamente na inicialização do container PostgreSQL através do volume:
```yaml
volumes:
  - ./db/init:/docker-entrypoint-initdb.d:ro
```