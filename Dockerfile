FROM golang:1.24.2 as builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

ENV APP_ENV=production
#TODO: наверное стоит заменить на development env, чтобы сделать авто подргузку example.env

CMD["./app/cmd/app"]

