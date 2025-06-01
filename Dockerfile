FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN ls -la /app
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./main"]
