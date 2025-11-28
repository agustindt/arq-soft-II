FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from search-api directory
COPY backend/search-api/go.mod backend/search-api/go.sum ./
RUN go mod download

# Copy the entire search-api directory
COPY backend/search-api/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8083
CMD ["./main"]
