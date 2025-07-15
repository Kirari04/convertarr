FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add sqlite gcc musl-dev

COPY . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o main.bin main.go

FROM alpine:latest

WORKDIR /app
VOLUME /app/database
VOLUME /app/imgs

RUN apk add --no-cache ffmpeg bash coreutils sqlite

COPY --from=builder /app/main.bin /app/main.bin

EXPOSE 8080

CMD ["./main.bin", "serve"]