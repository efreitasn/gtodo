FROM golang:1.14
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gtodo cmd/gtodo/main.go

FROM alpine:3.11
WORKDIR /app
COPY --from=0 /app/gtodo ./gtodo
COPY --from=0 /app/.cert ./.cert
COPY --from=0 /app/web ./web
CMD ["./gtodo"]
