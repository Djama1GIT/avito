FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o avito-microservice ./cmd/main.go

CMD ["./avito-microservice"]
