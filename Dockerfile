FROM golang:1.21.0

WORKDIR /lamoda-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build cmd/main.go

CMD ["./main"]