FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from users-api directory
COPY backend/users-api/go.mod backend/users-api/go.sum ./
RUN go mod download

# Copy the entire users-api directory
COPY backend/users-api/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8081
CMD ["./main"]
