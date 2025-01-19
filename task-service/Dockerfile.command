FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o task-command-service ./cmd/command

CMD ["./task-command-service"]
