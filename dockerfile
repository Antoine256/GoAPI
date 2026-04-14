FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8690

CMD ["./api"]