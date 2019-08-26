FROM golang:1.12.9

WORKDIR /go-todo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-todo cmd/go-todo/main.go

EXPOSE 8080

CMD ["./go-todo"]