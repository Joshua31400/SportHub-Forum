FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copier d'abord tout le code source
COPY . .

# Puis gérer les dépendances
RUN go mod tidy
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/github
RUN go get github.com/joho/godotenv
RUN go mod download

# Puis compiler
RUN go build -o main ./cmd

FROM alpine:3.18

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./main"]