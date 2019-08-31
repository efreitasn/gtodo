FROM golang:1.12.9

WORKDIR /gtodo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gtodo cmd/gtodo/main.go

EXPOSE 8080

CMD ["./gtodo"]