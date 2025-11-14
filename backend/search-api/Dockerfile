FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar go.mod y go.sum desde el context (raíz del proyecto)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente desde el context
COPY . .

# Construir el binario desde la raíz del módulo
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./backend/search-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8083
CMD ["./main"]
