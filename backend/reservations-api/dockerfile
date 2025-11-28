FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from reservations-api directory
COPY backend/reservations-api/go.mod backend/reservations-api/go.sum ./
RUN go mod download

# Copy the entire reservations-api directory
COPY backend/reservations-api/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8080
CMD ["./main"]
