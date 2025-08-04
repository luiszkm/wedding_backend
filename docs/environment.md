# Variáveis de Ambiente

Esta documentação descreve todas as variáveis de ambiente necessárias para configurar a Wedding Management API.

## Variáveis Obrigatórias

### Database

```bash
DATABASE_URL=postgres://user:password@localhost:5432/wedding_db
```

**Descrição**: URL de conexão com o banco PostgreSQL.

**Formato**: `postgres://[user[:password]@][host][:port][/dbname][?param1=value1&...]`

**Exemplos**:
- Local: `postgres://user:password@localhost:5432/wedding_db`
- Docker: `postgres://user:password@db:5432/wedding_db`
- Cloud: `postgres://user:pass@host.com:5432/wedding_db?sslmode=require`

---

### JWT Authentication

```bash
JWT_SECRET=seu-jwt-secret-muito-seguro-aqui-com-pelo-menos-32-caracteres
```

**Descrição**: Chave secreta para assinar tokens JWT.

**Requisitos**:
- Mínimo 32 caracteres
- Caracteres aleatórios e seguros
- **NUNCA** usar em produção o mesmo valor de desenvolvimento

**Geração Segura**:
```bash
# Linux/macOS
openssl rand -base64 32

# Ou usando Go
go run -c 'fmt.Println(base64.StdEncoding.EncodeToString(make([]byte, 32)))'
```

---

### Stripe Payment

```bash
STRIPE_SECRET_KEY=sk_test_51...
STRIPE_WEBHOOK_SECRET=whsec_...
```

**STRIPE_SECRET_KEY**:
- Chave secreta da API Stripe
- Formato test: `sk_test_...`
- Formato produção: `sk_live_...`

**STRIPE_WEBHOOK_SECRET**:
- Segredo do webhook endpoint
- Obtido no Stripe Dashboard
- Formato: `whsec_...`

**Como Obter**:
1. Acesse [Stripe Dashboard](https://dashboard.stripe.com)
2. API Keys para `STRIPE_SECRET_KEY`
3. Webhooks → Create Endpoint para `STRIPE_WEBHOOK_SECRET`

---

### Storage (Cloudflare R2 / AWS S3)

```bash
R2_ACCOUNT_ID=seu-cloudflare-account-id
R2_ACCESS_KEY_ID=sua-r2-access-key-id
R2_SECRET_ACCESS_KEY=sua-r2-secret-key
R2_BUCKET_NAME=wedding-photos
R2_PUBLIC_URL=https://pub-hashdorepositorio.r2.dev
```

**R2_ACCOUNT_ID**: 
- ID da conta Cloudflare
- Encontrado no dashboard R2

**R2_ACCESS_KEY_ID** / **R2_SECRET_ACCESS_KEY**:
- Credenciais de API
- Criadas em "Manage R2 API tokens"

**R2_BUCKET_NAME**:
- Nome do bucket para armazenar arquivos
- Deve existir antes de usar a aplicação

**R2_PUBLIC_URL**:
- URL pública do bucket
- Usado para gerar URLs de download
- Formato: `https://pub-[hash].r2.dev`

### Alternativa AWS S3

Se preferir usar AWS S3 em vez de Cloudflare R2:

```bash
R2_ACCOUNT_ID=              # deixe vazio
R2_ACCESS_KEY_ID=AKIA...    # AWS Access Key
R2_SECRET_ACCESS_KEY=...    # AWS Secret Key  
R2_BUCKET_NAME=my-s3-bucket
R2_PUBLIC_URL=https://my-s3-bucket.s3.us-east-1.amazonaws.com
```

---

## Configuração por Ambiente

### Desenvolvimento (.env)

```bash
# Development Environment
DATABASE_URL=postgres://user:password@localhost:5432/wedding_db
JWT_SECRET=development-jwt-secret-change-in-production-32chars
STRIPE_SECRET_KEY=sk_test_51ABCDEFghijklmnop...
STRIPE_WEBHOOK_SECRET=whsec_ABC123DEF456...
R2_ACCOUNT_ID=your-cloudflare-account-id
R2_ACCESS_KEY_ID=your-r2-access-key
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_BUCKET_NAME=wedding-photos-dev
R2_PUBLIC_URL=https://pub-yourhash.r2.dev
```

### Testing (.env.test)

```bash
# Test Environment
DATABASE_URL=postgres://user:password@localhost:5432/wedding_test_db
JWT_SECRET=test-jwt-secret-32-characters-long
STRIPE_SECRET_KEY=sk_test_dummy_key_for_testing
STRIPE_WEBHOOK_SECRET=whsec_test_dummy_secret
R2_ACCOUNT_ID=test-account-id
R2_ACCESS_KEY_ID=test-access-key
R2_SECRET_ACCESS_KEY=test-secret-key
R2_BUCKET_NAME=wedding-test-bucket
R2_PUBLIC_URL=https://test.r2.dev
```

### Staging (.env.staging)

```bash
# Staging Environment  
DATABASE_URL=postgres://user:securepass@staging-db:5432/wedding_staging
JWT_SECRET=staging-jwt-secret-generated-secure-key
STRIPE_SECRET_KEY=sk_test_staging_key...
STRIPE_WEBHOOK_SECRET=whsec_staging_secret...
R2_ACCOUNT_ID=staging-account-id
R2_ACCESS_KEY_ID=staging-access-key
R2_SECRET_ACCESS_KEY=staging-secret-key
R2_BUCKET_NAME=wedding-staging-bucket
R2_PUBLIC_URL=https://pub-staging.r2.dev
```

### Produção (.env.production)

```bash
# Production Environment
DATABASE_URL=postgres://prod_user:very_secure_password@prod-db:5432/wedding_prod?sslmode=require
JWT_SECRET=production-super-secure-jwt-secret-at-least-32-chars
STRIPE_SECRET_KEY=sk_live_production_key...
STRIPE_WEBHOOK_SECRET=whsec_production_webhook_secret...
R2_ACCOUNT_ID=production-account-id
R2_ACCESS_KEY_ID=production-access-key
R2_SECRET_ACCESS_KEY=production-secret-key
R2_BUCKET_NAME=wedding-production-bucket
R2_PUBLIC_URL=https://pub-production.r2.dev
```

---

## Variáveis Opcionais

### Server Configuration

```bash
# Porta do servidor (padrão: 3000)
PORT=3000

# Timeout de requisições (padrão: 30s)
REQUEST_TIMEOUT=30s

# Modo debug (padrão: false)
DEBUG=true
```

### Database Pool

```bash
# Máximo de conexões no pool (padrão: 10)
DB_MAX_CONNECTIONS=20

# Timeout de conexão (padrão: 5s)
DB_CONNECTION_TIMEOUT=10s

# Tempo de vida da conexão (padrão: 1h)
DB_MAX_LIFETIME=2h
```

### CORS Configuration

```bash
# Origens permitidas (padrão: *)
CORS_ALLOWED_ORIGINS=https://meusite.com,https://app.meusite.com

# Headers permitidos
CORS_ALLOWED_HEADERS=Authorization,Content-Type,X-Requested-With

# Métodos permitidos
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
```

### Rate Limiting

```bash
# Requests por minuto por IP (padrão: 100)
RATE_LIMIT_RPM=200

# Burst size (padrão: 10)
RATE_LIMIT_BURST=20
```

---

## Validação de Variáveis

### Checklist de Configuração

**Development**:
- [ ] `DATABASE_URL` aponta para DB local
- [ ] `JWT_SECRET` tem pelo menos 32 caracteres
- [ ] `STRIPE_SECRET_KEY` usa `sk_test_`
- [ ] Storage configurado para ambiente de teste

**Production**:
- [ ] `DATABASE_URL` usa SSL (`sslmode=require`)
- [ ] `JWT_SECRET` gerado com `openssl rand -base64 32`
- [ ] `STRIPE_SECRET_KEY` usa `sk_live_`
- [ ] Storage configurado para produção
- [ ] Todas as chaves são únicas e seguras

### Script de Validação

```bash
#!/bin/bash
# validate-env.sh

echo "Validating environment variables..."

# Check required variables
required_vars=(
    "DATABASE_URL"
    "JWT_SECRET" 
    "STRIPE_SECRET_KEY"
    "STRIPE_WEBHOOK_SECRET"
    "R2_ACCESS_KEY_ID"
    "R2_SECRET_ACCESS_KEY"
    "R2_BUCKET_NAME"
    "R2_PUBLIC_URL"
)

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "❌ Missing required variable: $var"
        exit 1
    else
        echo "✅ $var is set"
    fi
done

# Validate JWT_SECRET length
if [ ${#JWT_SECRET} -lt 32 ]; then
    echo "❌ JWT_SECRET must be at least 32 characters"
    exit 1
fi

# Validate Stripe key format
if [[ ! $STRIPE_SECRET_KEY =~ ^sk_(test|live)_ ]]; then
    echo "❌ STRIPE_SECRET_KEY has invalid format"
    exit 1
fi

# Validate webhook secret format
if [[ ! $STRIPE_WEBHOOK_SECRET =~ ^whsec_ ]]; then
    echo "❌ STRIPE_WEBHOOK_SECRET has invalid format"  
    exit 1
fi

echo "✅ All environment variables are valid!"
```

---

## Docker e Docker Compose

### docker-compose.yml com Environment

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/wedding_db
      - JWT_SECRET=${JWT_SECRET}
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
      - STRIPE_WEBHOOK_SECRET=${STRIPE_WEBHOOK_SECRET}
      - R2_ACCOUNT_ID=${R2_ACCOUNT_ID}
      - R2_ACCESS_KEY_ID=${R2_ACCESS_KEY_ID}
      - R2_SECRET_ACCESS_KEY=${R2_SECRET_ACCESS_KEY}
      - R2_BUCKET_NAME=${R2_BUCKET_NAME}
      - R2_PUBLIC_URL=${R2_PUBLIC_URL}
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: wedding_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/init:/docker-entrypoint-initdb.d:ro

volumes:
  postgres_data:
```

### Usando env_file

```yaml
services:
  api:
    build: .
    env_file:
      - .env
    depends_on:
      - db
```

---

## Segurança

### Proteção de Secrets

**❌ NUNCA faça**:
```bash
# Não comitar no Git
git add .env

# Não logar secrets  
log.Printf("JWT Secret: %s", jwtSecret)

# Não usar valores padrão em produção
JWT_SECRET=default-secret
```

**✅ Boas práticas**:
```bash
# Adicione ao .gitignore
echo ".env*" >> .gitignore

# Use gerenciadores de secrets em produção
# - AWS Secrets Manager
# - HashiCorp Vault  
# - Kubernetes Secrets

# Rotacione secrets regularmente
# - JWT secrets a cada deploy
# - Database passwords mensalmente
# - API keys conforme política
```

### Exemplo de .env.example

```bash
# Copy this file to .env and fill with your values

# Database
DATABASE_URL=postgres://user:password@localhost:5432/wedding_db

# JWT (generate with: openssl rand -base64 32)  
JWT_SECRET=your-jwt-secret-here

# Stripe (get from dashboard.stripe.com)
STRIPE_SECRET_KEY=sk_test_your_stripe_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Storage (Cloudflare R2 or AWS S3)
R2_ACCOUNT_ID=your-r2-account-id
R2_ACCESS_KEY_ID=your-r2-access-key  
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_BUCKET_NAME=your-bucket-name
R2_PUBLIC_URL=https://pub-yourhash.r2.dev
```

---

## Troubleshooting

### Erros Comuns

**Erro**: `pq: SSL is not enabled on the server`
```bash
# Solução: Adicionar sslmode=disable para desenvolvimento
DATABASE_URL=postgres://user:password@localhost:5432/wedding_db?sslmode=disable
```

**Erro**: `Stripe webhook signature verification failed`
```bash
# Verificar se STRIPE_WEBHOOK_SECRET está correto
# Verificar se endpoint está configurado no Stripe Dashboard
```

**Erro**: `S3 Access Denied`
```bash
# Verificar credenciais R2_ACCESS_KEY_ID e R2_SECRET_ACCESS_KEY
# Verificar permissões do bucket
# Verificar se bucket existe
```

**Erro**: `JWT token invalid`
```bash
# Verificar se JWT_SECRET é o mesmo usado na geração
# Verificar se token não expirou (7 dias padrão)
```