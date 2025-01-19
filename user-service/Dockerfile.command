FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o user-command-service ./cmd/command

CMD ["./user-command-service"]
