# Deployment e Configuração de Ambiente

Esta documentação descreve como configurar e fazer deploy da Wedding Management API.

## Pré-requisitos

- **Go 1.23+**
- **Docker e Docker Compose**
- **PostgreSQL 16+**
- **Conta Stripe** (para pagamentos)
- **AWS S3 ou Cloudflare R2** (para storage de arquivos)

---

## Configuração de Ambiente

### Variáveis de Ambiente Obrigatórias

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/wedding_db

# JWT Authentication
JWT_SECRET=seu-jwt-secret-super-seguro-aqui

# Stripe Payment
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# Cloudflare R2 / AWS S3 Storage
R2_ACCOUNT_ID=seu-account-id
R2_ACCESS_KEY_ID=sua-access-key
R2_SECRET_ACCESS_KEY=sua-secret-key
R2_BUCKET_NAME=wedding-photos
R2_PUBLIC_URL=https://pub-hash.r2.dev
```

### Exemplo de .env para Desenvolvimento

```bash
# Development Environment
DATABASE_URL=postgres://user:password@localhost:5432/wedding_db
JWT_SECRET=development-jwt-secret-change-in-production
STRIPE_SECRET_KEY=sk_test_51...
STRIPE_WEBHOOK_SECRET=whsec_...
R2_ACCOUNT_ID=your-cloudflare-account-id
R2_ACCESS_KEY_ID=your-r2-access-key
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_BUCKET_NAME=wedding-photos-dev
R2_PUBLIC_URL=https://pub-yourhash.r2.dev
```

---

## Setup Local com Docker Compose

### 1. Clone e Configure

```bash
git clone <repository-url>
cd wedding_backend
cp .env.example .env  # Configure suas variáveis
```

### 2. Execute com Docker Compose

```bash
# Inicia toda a stack (PostgreSQL + API)
docker-compose up --build

# Ou execute apenas o banco de dados
docker-compose up db
```

### 3. Verificação

A API estará disponível em: `http://localhost:3000`

O PostgreSQL estará disponível em: `localhost:5432`

---

## Desenvolvimento Local (Sem Docker)

### 1. Configurar PostgreSQL

```bash
# Instalar PostgreSQL
# Ubuntu/Debian
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql

# Criar banco de dados
createdb wedding_db
```

### 2. Executar Migrações

```sql
-- Execute os arquivos na pasta db/init/ em ordem:
psql -d wedding_db -f db/init/01-init.sql
psql -d wedding_db -f db/init/02-seed-plans.sql
```

### 3. Executar a Aplicação

```bash
# Download de dependências
go mod download

# Executar em modo desenvolvimento
go run ./cmd/api/main.go

# Ou compilar e executar
go build -o server ./cmd/api/main.go
./server
```

---

## Configuração de Serviços Externos

### Stripe Setup

1. **Criar conta no Stripe**
   - Acesse [stripe.com](https://stripe.com)
   - Configure sua conta

2. **Configurar Products e Prices**
   ```bash
   # Os IDs devem coincidir com os no 02-seed-plans.sql
   Mensal: price_1RYz6wE2D6lLQS0txXt4h65V
   Trimestral: price_1RYz82E2D6lLQS0t7tmdebs3
   Semestral: price_1RYz8ME2D6lLQS0tsJNujNxN
   ```

3. **Configurar Webhook**
   - URL: `https://seu-dominio.com/v1/webhooks/stripe`
   - Eventos: `customer.subscription.created`, `customer.subscription.updated`, `customer.subscription.deleted`

### Cloudflare R2 Setup

1. **Criar bucket R2**
   ```bash
   # Via Cloudflare Dashboard
   1. Acesse R2 Object Storage
   2. Criar novo bucket
   3. Configurar público se necessário
   ```

2. **Gerar API Token**
   ```bash
   1. Acesse API Tokens
   2. Criar token com permissões R2
   3. Anotar Account ID, Access Key, Secret Key
   ```

### AWS S3 (Alternativa ao R2)

```bash
# Se preferir usar AWS S3 em vez de Cloudflare R2
R2_ACCOUNT_ID=        # deixe vazio
R2_ACCESS_KEY_ID=your-aws-access-key
R2_SECRET_ACCESS_KEY=your-aws-secret-key  
R2_BUCKET_NAME=your-s3-bucket
R2_PUBLIC_URL=https://your-bucket.s3.region.amazonaws.com
```

---

## Deploy em Produção

### Docker Production Build

```dockerfile
# Dockerfile otimizado para produção
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["/app/server"]
```

### Build e Deploy

```bash
# Build da imagem
docker build -t wedding-api:latest .

# Deploy (exemplo com Docker)
docker run -d \
  --name wedding-api \
  --env-file .env \
  -p 3000:3000 \
  wedding-api:latest
```

### Deploy com Docker Compose (Produção)

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "3000:3000"
    env_file:
      - .env.production
    depends_on:
      - db
    restart: always

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/init:/docker-entrypoint-initdb.d:ro
    restart: always

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - api
    restart: always

volumes:
  postgres_data:
```

---

## Configuração de Proxy Reverso (Nginx)

### nginx.conf

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api {
        server api:3000;
    }

    server {
        listen 80;
        server_name seu-dominio.com;

        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name seu-dominio.com;

        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;

        location / {
            proxy_pass http://api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Configuração especial para uploads grandes
        client_max_body_size 50M;
    }
}
```

---

## Healthcheck e Monitoramento

### Endpoint de Health

Adicione um endpoint básico de health:

```go
// Adicionar em main.go
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().ISO8601() + `"}`))
})
```

### Docker Healthcheck

```dockerfile
# Adicionar ao Dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1
```

---

## Backup e Segurança

### Backup do Banco de Dados

```bash
# Backup automático
docker exec -t wedding_db pg_dumpall -c -U user > dump_`date +%d-%m-%Y"_"%H_%M_%S`.sql

# Restore
cat dump_file.sql | docker exec -i wedding_db psql -U user -d wedding_db
```

### Variáveis de Segurança

```bash
# Produção - NUNCA use valores de desenvolvimento
JWT_SECRET=$(openssl rand -base64 32)
DB_PASSWORD=$(openssl rand -base64 16)
```

### SSL/TLS

Configure certificados SSL usando Let's Encrypt:

```bash
# Instalação do Certbot
sudo apt install certbot python3-certbot-nginx

# Obter certificado
sudo certbot --nginx -d seu-dominio.com

# Renovação automática
sudo crontab -e
# Adicione: 0 12 * * * /usr/bin/certbot renew --quiet
```

---

## Troubleshooting

### Problemas Comuns

1. **Erro de conexão com banco**
   ```bash
   # Verificar se PostgreSQL está rodando
   docker-compose logs db
   
   # Testar conexão
   psql -h localhost -U user -d wedding_db
   ```

2. **Erro de upload de arquivo**
   ```bash
   # Verificar configuração R2/S3
   # Verificar permissões do bucket
   # Verificar variáveis de ambiente
   ```

3. **Webhook Stripe não funciona**
   ```bash
   # Verificar endpoint público
   # Verificar STRIPE_WEBHOOK_SECRET
   # Verificar logs do Stripe Dashboard
   ```

### Logs

```bash
# Logs do Docker Compose
docker-compose logs -f

# Logs específicos da API
docker-compose logs -f api

# Logs do banco
docker-compose logs -f db
```

---

## Performance e Otimização

### Configurações de Produção

```bash
# Variáveis de otimização Go
GOGC=100
GOMAXPROCS=4

# PostgreSQL tuning
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB
```

### Monitoramento

Considere usar ferramentas como:
- **Prometheus + Grafana** para métricas
- **Sentry** para error tracking  
- **New Relic** ou **DataDog** para APM