FROM golang:1.21 as builder

WORKDIR /app
COPY . .
# this copy command will only work if you're running the docker build command from the root of the repository and specifying -f ./docker/nvidia/Dockerfile
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o main.bin main.go

FROM ghcr.io/jrottenberg/ffmpeg:5.1-alpine

WORKDIR /app
VOLUME /app/database

RUN apk add --no-cache bash

COPY --from=builder /app/main.bin /app/main.bin

EXPOSE 8080

ENTRYPOINT ["./main.bin"]
CMD ["serve"]
