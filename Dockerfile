FROM golang:1.24.2 as builder

WORKDIR /app
COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app ./cmd/app/main.go

FROM debian:bookworm-slim

WORKDIR /root
COPY --from=builder /app/app .
COPY --from=builder /app/example.env .env
#TODO: проследить, что копируется энв норм
CMD ["./app"]