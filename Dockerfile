FROM golang:1.26.7-alpine AS builder

WORKDIR /app

# Copiar go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build API e Worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./cmd/worker

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binários
COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/.env.development .

# Expor porta da API
EXPOSE 8080

# Comando padrão (pode ser sobrescrito)
ENTRYPOINT ["./api"]
