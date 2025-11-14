FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar go.mod y go.sum desde la raíz del proyecto
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Construir el binario
WORKDIR /app/backend/activities-api
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Imagen final minimal
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copiar el binario desde builder
COPY --from=builder /main .

EXPOSE 8082
CMD ["./main"]
