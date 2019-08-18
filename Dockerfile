FROM golang:1.12.9

WORKDIR /go-todo
COPY . .
RUN go get ./...

EXPOSE 8080

CMD ["go", "run", "cmd/go-todo/main.go"]