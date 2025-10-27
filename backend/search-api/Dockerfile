FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar archivos desde la ruta correcta
COPY backend/search-api/go.mod backend/search-api/go.sum ./
RUN go mod download

# Copiar el código fuente
COPY backend/search-api/ ./

# Compilar
RUN go build -o search-api .

# Imagen final
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/search-api .

EXPOSE 8083

CMD ["./search-api"]