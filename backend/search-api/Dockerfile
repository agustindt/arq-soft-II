FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/backend/search-api
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8083
CMD ["./main"]
