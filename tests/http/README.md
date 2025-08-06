# Testes HTTP - Wedding API

Esta pasta contém arquivos `.http` para testar a API usando REST Client (extensão do VS Code) ou similar.

## 📁 Estrutura dos Arquivos

### 🔐 Autenticação
- **`auth.http`** - Registro e login de usuários

### 🎉 Eventos  
- **`events.http`** - Criação e gestão de eventos de casamento

### 👥 Convidados
- **`guests.http`** - Criação de grupos de convidados e confirmação de presença (RSVP)

### 🎁 Presentes
- **`gifts.http`** - Criação de presentes integrais e fracionados (ADR-004)
- **`gift-selections.http`** - Seleção de presentes com suporte a cotas (ADR-004)

### 📢 Comunicados
- **`communications.http`** - Sistema de comunicados (ADR-003)

### 💬 Recados
- **`messages.http`** - Mural de recados com moderação

### 📸 Galeria
- **`gallery.http`** - Upload e gestão de fotos

### 💳 Cobrança
- **`billing.http`** - Planos e assinaturas com Stripe

## 🚀 Como Usar

### 1. Instalar REST Client
No VS Code, instale a extensão **REST Client** by Huachao Mao.

### 2. Executar em Sequência
Para um teste completo, execute os arquivos nesta ordem:

1. **`auth.http`** - Para obter token de autenticação
2. **`events.http`** - Para criar evento
3. **`guests.http`** - Para criar grupos de convidados  
4. **`gifts.http`** - Para criar presentes (integral + fracionado)
5. **`gift-selections.http`** - Para testar seleções de presentes
6. **`communications.http`** - Para testar comunicados
7. **`messages.http`** - Para testar recados
8. **`gallery.http`** - Para testar upload de fotos
9. **`billing.http`** - Para testar cobrança

### 3. Variáveis Automáticas
Os arquivos usam variáveis que são definidas automaticamente pelos responses:
- `{{authToken}}` - Token JWT do login
- `{{eventoId}}` - ID do evento criado
- `{{chaveAcesso}}` - Chave de acesso dos convidados
- `{{presenteIntegralId}}` - ID do presente integral
- `{{presenteFracionadoId}}` - ID do presente fracionado
- etc.

### 4. Configuração do Servidor
Certifique-se de que o servidor esteja rodando em `http://localhost:3000` ou ajuste a variável `@baseUrl` nos arquivos.

## ✨ Recursos Testados

### 🆕 Novos Recursos (ADRs)
- **ADR-003**: Sistema completo de comunicados
- **ADR-004**: Presentes fracionados com cotas
  - Criação de presentes com `valorTotal` e `numeroCotas`
  - Seleção com quantidade variável
  - Status `PARCIALMENTE_SELECIONADO`
  - Informações detalhadas de cotas disponíveis/selecionadas

### 🔧 Recursos Existentes
- Autenticação JWT completa
- CRUD de eventos, convidados, presentes
- Sistema de RSVP
- Mural de recados com moderação
- Galeria de fotos com rótulos
- Integração com Stripe

## 🐛 Troubleshooting

- **Token expirado**: Re-execute `auth.http` para obter novo token
- **IDs não encontrados**: Execute os arquivos em sequência respeitando as dependências
- **Erros 500**: Verifique se o banco de dados está rodando e as migrações foram aplicadas