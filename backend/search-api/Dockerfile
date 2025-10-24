# Search API Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project (needed for monorepo imports)
COPY . .

# Build the specific service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./backend/search-api/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8083
CMD ["./main"]