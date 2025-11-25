FROM golang:1.24.2 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app ./cmd/app/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root

COPY --from=builder /app/app .
COPY --from=builder /app/pkg/db/migrations ./pkg/db/migrations
COPY --from=builder /app/example.env .env

EXPOSE 8080

CMD ["./app"]