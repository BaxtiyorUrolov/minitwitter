FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY . /app

RUN go build -o main cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .

COPY .env /app

COPY ./migrations /app/migrations

CMD ["/app/main"]
