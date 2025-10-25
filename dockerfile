# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go

# Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
  && adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=builder /app/server .
ENV PORT=8080
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/server"]
