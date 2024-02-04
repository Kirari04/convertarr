FROM golang:1.21 as builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o main.bin main.go

FROM alpine:3.14

WORKDIR /app
VOLUME /app/database

RUN apk add --no-cache ffmpeg bash

COPY --from=builder /app/main.bin /app/main.bin

EXPOSE 8080

CMD ["./main.bin", "serve"]