FROM golang:1.12.9

WORKDIR /go-todo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "cmd/go-todo/main.go"]