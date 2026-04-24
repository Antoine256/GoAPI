FROM golang:1.25.9-alpine AS builder

RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api .

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /app/logs && chmod 777 /app/logs

COPY --from=builder /app/api .

EXPOSE 8690

CMD ["./api"]