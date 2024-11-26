FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o main.bin main.go

FROM alpine:latest

WORKDIR /app
VOLUME /app/database
VOLUME /app/imgs

RUN apk add --no-cache ffmpeg bash coreutils

COPY --from=builder /app/main.bin /app/main.bin

EXPOSE 8080

CMD ["./main.bin", "serve"]