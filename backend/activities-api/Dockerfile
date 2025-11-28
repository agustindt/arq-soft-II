FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from activities-api directory
COPY backend/activities-api/go.mod backend/activities-api/go.sum ./
RUN go mod download

# Copy the entire activities-api directory
COPY backend/activities-api/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Imagen final minimal
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copiar el binario desde builder
COPY --from=builder /main .

EXPOSE 8082
CMD ["./main"]
