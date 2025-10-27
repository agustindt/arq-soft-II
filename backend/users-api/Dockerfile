# Etapa de build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar y descargar dependencias del microservicio
COPY backend/users-api/go.mod backend/users-api/go.sum ./
RUN go mod download

# Copiar el código fuente
COPY backend/users-api/ ./

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Etapa final (runtime)
FROM alpine:latest

# Instalar certificados SSL para conexiones seguras
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el ejecutable desde la etapa anterior
COPY --from=builder /main .

EXPOSE 8081

CMD ["./main"]
