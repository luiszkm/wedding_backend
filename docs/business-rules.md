# Regras de Negócio

Esta documentação descreve as regras de negócio implementadas na Wedding Management API, organizadas por domínio.

---

## 👥 Guest (Convidados)

### Grupos de Convidados

**Regras de Criação:**
- Cada grupo deve ter uma chave de acesso única por evento
- A chave de acesso é obrigatória e não pode ser vazia
- Deve haver pelo menos um convidado no grupo
- A chave de acesso é case-sensitive

**Regras de RSVP:**
- Convidados podem confirmar presença usando apenas a chave de acesso
- Status possíveis: `PENDENTE`, `CONFIRMADO`, `RECUSADO`
- Todos os convidados do grupo podem ter status diferentes
- Uma vez confirmado, o status pode ser alterado

**Validações:**
```go
// Chave de acesso obrigatória
if chaveDeAcesso == "" {
    return ErrChaveDeAcessoObrigatoria
}

// Pelo menos um convidado
if len(nomes) == 0 {
    return ErrPeloMenosUmConvidado
}
```

**Casos de Uso:**
- Famílias com múltiplos membros
- Grupos de amigos
- Padrinhos e madrinhas
- Convidados individuais

---

## 🎁 Gift (Presentes)

### Lista de Presentes

**Tipos de Presente:**
1. **INTEGRAL**: Presente completo selecionado por um único convidado
   - Status: `DISPONIVEL` ou `SELECIONADO`
   - Uma vez selecionado, fica indisponível para outros
   
2. **FRACIONADO**: Presente dividido em cotas que podem ser selecionadas por múltiplos convidados
   - Status: `DISPONIVEL`, `PARCIALMENTE_SELECIONADO`, ou `SELECIONADO`
   - Valor total dividido em número específico de cotas
   - Cada cota tem valor individual e status próprio

**Modalidades de Detalhes:**
1. **PRODUTO_EXTERNO**: Link para loja externa
   - Campo `detalhes_link_loja` é obrigatório
   - Campo `detalhes_chave_pix` deve ser nulo

2. **PIX**: Doação via PIX
   - Campo `detalhes_chave_pix` é obrigatório
   - Campo `detalhes_link_loja` deve ser nulo

**Regras de Seleção:**

*Para Presentes Integrais:*
- Convidados selecionam o presente completo usando grupo de convidados
- Presente passa de `DISPONIVEL` para `SELECIONADO`
- Uma seleção é registrada com data/hora

*Para Presentes Fracionados:*
- Convidados podem selecionar uma ou múltiplas cotas
- Cada cota selecionada é vinculada ao grupo que selecionou
- Status do presente atualiza automaticamente:
  - `DISPONIVEL`: Todas as cotas disponíveis
  - `PARCIALMENTE_SELECIONADO`: Algumas cotas selecionadas
  - `SELECIONADO`: Todas as cotas selecionadas

**Categorização:**
- Presentes podem ter categoria (rótulo)
- Categorias disponíveis: `MAIN`, `CASAMENTO`, `LUADEMEL`, `HISTORIA`, `FAMILIA`, `OUTROS`, `COZINHA`, `SALA`, `QUARTO`, `BANHEIRO`, `JARDIM`, `DECORACAO`, `ELETRONICOS`, `UTENSILIOS`
- Casais podem marcar presentes como favoritos

**Sistema de Cotas:**
- Cada cota tem valor individual calculado automaticamente (valor_total ÷ numero_cotas)
- Cotas são numeradas sequencialmente (1, 2, 3...)
- Não é possível selecionar cotas parciais (ex: 0.5 cota)
- Sistema garante integridade através de constraints de banco

**Validações:**
```sql
-- Constraint de banco para validar tipos
CONSTRAINT chk_detalhes CHECK (
    (detalhes_tipo = 'PRODUTO_EXTERNO' AND detalhes_link_loja IS NOT NULL) OR
    (detalhes_tipo = 'PIX' AND detalhes_chave_pix IS NOT NULL)
)
```

---

## 💬 MessageBoard (Recados)

### Sistema de Recados

**Regras de Moderação:**
- Status possíveis: `PENDENTE`, `APROVADO`, `REJEITADO`
- Recados iniciam como `PENDENTE`
- Apenas recados `APROVADO` aparecem na listagem pública
- Administradores podem aprovar/rejeitar recados

**Regras de Criação:**
- Recados devem estar vinculados a um evento
- Autor (nome) é obrigatório
- Texto não pode ser vazio
- Data de criação é automaticamente definida

**Controle de Qualidade:**
- Sistema de moderação manual
- Casais podem marcar recados como favoritos
- Recados rejeitados ficam ocultos do público

**Permissões:**
- **Público**: Pode criar recados, ver aprovados
- **Admin/Casal**: Pode moderar, ver todos, favoritar

---

## 📸 Gallery (Galeria)

### Sistema de Fotos

**Regras de Upload:**
- Fotos são armazenadas em storage externo (S3/R2)
- Cada foto gera uma URL pública
- Storage key é mantido para gerenciamento

**Sistema de Rótulos:**
- Múltiplos rótulos por foto (relação N:N)
- Rótulos disponíveis: `MAIN`, `CASAMENTO`, `LUADEMEL`, `HISTORIA`, `FAMILIA`, `OUTROS`
- Facilita organização e filtros

**Controle de Favoritos:**
- Casais podem marcar fotos como favoritas
- Sistema de destaque para fotos importantes

**Gerenciamento:**
- Deleção remove foto do storage e banco
- Órfãos de storage são evitados através de transações

---

## 👤 IAM (Identity & Access Management)

### Sistema de Usuários

**Regras de Registro:**
- Email deve ser único na plataforma
- Senha deve atender critérios mínimos de segurança
- Nome é obrigatório
- Telefone é opcional

**Autenticação:**
- Login via email/senha
- JWT tokens com expiração de 7 dias
- Tokens incluem user ID para identificação

**Segurança:**
- Senhas são hasheadas (nunca armazenadas em plain text)
- JWT secret deve ser complexo e secreto
- Tokens podem ser revogados (através de blacklist se necessário)

---

## 📅 Event (Eventos)

### Gestão de Eventos

**Regras de Criação:**
- Cada usuário pode ter múltiplos eventos
- URL slug deve ser único globalmente
- Nome do evento é obrigatório
- Data é opcional (pode ser definida depois)

**Tipos de Evento:**
- `CASAMENTO`: Evento principal
- `ANIVERSARIO`: Comemorações
- `CHA_DE_BEBE`: Baby shower
- `OUTRO`: Eventos diversos

**URL Slug:**
- Deve ser único na plataforma
- Usado para URLs públicas amigáveis
- Formato sugerido: `casamento-joao-maria-2024`

---

## 💳 Billing (Cobrança)

### Sistema de Assinaturas

**Planos Disponíveis:**
1. **Mensal**: R$ 99,90 - 1 evento - 30 dias
2. **Trimestral**: R$ 279,90 - 3 eventos - 90 dias
3. **Semestral**: R$ 539,90 - 5 eventos - 180 dias

**Regras de Assinatura:**
- Usuário pode ter apenas uma assinatura ativa
- Limite de eventos por plano é respeitado
- Assinatura expira automaticamente na data_fim

**Status de Assinatura:**
- `PENDENTE`: Aguardando pagamento
- `ATIVA`: Pagamento confirmado, pode usar recursos
- `EXPIRADA`: Plano vencido
- `CANCELADA`: Cancelada pelo usuário ou falta de pagamento

**Integração Stripe:**
- Webhooks processam eventos de pagamento
- Subscription ID do Stripe é mantido para referência
- Pagamentos são processados via Stripe

**Regras de Negócio:**
- Usuário sem assinatura ativa não pode criar eventos
- Eventos existentes ficam inacessíveis se assinatura expirar
- Renovação reativa eventos suspensos

---

## 🔐 Segurança e Permissões

### Controle de Acesso

**Rotas Públicas:**
- Registro e login de usuários
- RSVP por chave de acesso
- Visualização de presentes públicos
- Visualização de recados aprovados
- Listagem de planos
- Webhooks do Stripe

**Rotas Protegidas (JWT Required):**
- Todas as operações de gerenciamento
- Criação e edição de conteúdo
- Moderação de recados
- Upload de fotos
- Gestão de assinaturas

**Isolamento de Dados:**
- Usuários só acessam seus próprios eventos
- Validação de ownership em todas as operações
- Eventos são isolados por usuário

---

## 🔄 Regras de Integridade

### Cascata de Deleção

**Ao deletar usuário:**
- Eventos são deletados
- Assinaturas são canceladas

**Ao deletar evento:**
- Grupos de convidados são deletados
- Presentes são deletados
- Recados são deletados
- Fotos são deletadas
- Seleções são deletadas

**Ao deletar grupo de convidados:**
- Convidados são deletados
- Seleções vinculadas são removidas

### Validações de Integridade

**Referências obrigatórias:**
- Todo evento deve ter um usuário
- Todo convidado deve ter um grupo
- Toda seleção deve referenciar presente válido

**Unicidade:**
- Email de usuário
- Chave de acesso por evento
- URL slug de evento
- Nome de plano

---

## 📊 Regras de Auditoria

### Timestamps

**Criação automática:**
- `created_at` em usuários, eventos, grupos, recados, fotos
- Fuso horário: America/Sao_Paulo

**Atualização:**
- `updated_at` em grupos de convidados
- Atualizado em modificações

### Logs de Atividade

**Eventos importantes:**
- Criação de usuário
- Login
- Criação de evento
- Confirmação de RSVP
- Seleção de presente
- Aprovação de recado

---

## 🚫 Regras de Validação

### Validações de Entrada

**Strings obrigatórias:**
- Não podem ser vazias ou apenas espaços
- Tamanhos máximos respeitados

**UUIDs:**
- Formato válido obrigatório
- Referências devem existir

**Emails:**
- Formato válido
- Únicos na plataforma

**URLs:**
- Formato válido para links externos
- Acessibilidade verificada quando possível

### Regras de Negócio por Campo

**Chave de Acesso:**
- Mínimo 3 caracteres
- Máximo 255 caracteres
- Sem espaços ou caracteres especiais

**Nomes:**
- Mínimo 2 caracteres
- Máximo 255 caracteres
- Apenas letras, espaços e acentos

**Preços:**
- Sempre em centavos (inteiros)
- Valores positivos
- Máximo reasonable (ex: 999999999 centavos)

Essas regras garantem a consistência e integridade dos dados, além de uma boa experiência do usuário.