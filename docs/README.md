# Wedding Management API - Documentação

Esta pasta contém toda a documentação do projeto Wedding Management API.

## Índice

### 📚 Documentação Principal
- [**API Endpoints**](./api-endpoints.md) - Documentação completa de todos os endpoints da API
- [**Arquitetura**](./architecture.md) - Visão geral da arquitetura e padrões utilizados
- [**Banco de Dados**](./database.md) - Schema do banco de dados e relacionamentos

### 🚀 Guias de Configuração
- [**Deployment**](./deployment.md) - Guia de deploy e configuração de ambiente
- [**Environment**](./environment.md) - Variáveis de ambiente necessárias

### 💼 Regras de Negócio
- [**Business Rules**](./business-rules.md) - Regras de negócio e lógica do domínio
- [**Domain Models**](./domain-models.md) - Modelos de domínio e suas responsabilidades

### 🔧 Desenvolvimento
- [**Development Guide**](./development.md) - Guia de desenvolvimento local
- [**Testing**](./testing.md) - Estratégias e padrões de teste

## Sobre o Projeto

O Wedding Management API é uma aplicação Go que fornece funcionalidades completas para gestão de casamentos, incluindo:

- 👥 **Gestão de Convidados**: Grupos de convidados com confirmação de presença
- 🎁 **Lista de Presentes**: Registro e seleção de presentes pelos convidados
- 💬 **Mural de Recados**: Sistema de mensagens com moderação
- 📸 **Galeria de Fotos**: Upload e organização de fotos do evento
- 👤 **Autenticação**: Sistema de usuários com JWT
- 📅 **Eventos**: Gestão de eventos de casamento
- 💳 **Billing**: Integração com Stripe para planos de assinatura

## Tecnologias Utilizadas

- **Backend**: Go 1.23 com Chi Router
- **Banco de Dados**: PostgreSQL com pgx
- **Autenticação**: JWT
- **Pagamentos**: Stripe
- **Storage**: AWS S3/Cloudflare R2
- **Containerização**: Docker e Docker Compose