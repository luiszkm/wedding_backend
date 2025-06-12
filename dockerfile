# Estágio 1: O Compilador (Builder)
# Usamos uma imagem oficial do Go com Alpine Linux, que é leve.
FROM golang:1.23-alpine AS builder

# Define o diretório de trabalho dentro do container.
WORKDIR /app

# Copia os arquivos de gerenciamento de módulos primeiro.
# Isso aproveita o cache do Docker: as dependências só serão baixadas novamente se o go.mod/sum mudar.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código-fonte da nossa aplicação.
COPY . .

# Compila a aplicação.
# - CGO_ENABLED=0: Compila um binário estático, sem depender de bibliotecas C do sistema. Essencial para rodar em imagens mínimas.
# - -ldflags="-w -s": Remove informações de debug e a tabela de símbolos, reduzindo o tamanho do executável.
# - -o /app/server: Define o nome e o local do arquivo de saída.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go


# Estágio 2: A Imagem Final de Produção
# Começamos com uma imagem Alpine Linux limpa, que é muito pequena e segura.
FROM alpine:latest

# Define o diretório de trabalho.
WORKDIR /app

# Copia APENAS o executável compilado do estágio 'builder'.
# Nenhum código-fonte ou ferramenta de compilação é incluído na imagem final.
COPY --from=builder /app/server .

# Expõe a porta que nossa aplicação usa para o mundo exterior do container.
EXPOSE 3000

# Define o comando que será executado quando o container iniciar.
CMD ["/app/server"]