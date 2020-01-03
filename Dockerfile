FROM golang:1.13 as base
WORKDIR /gtodo
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base as build
RUN CGO_ENABLED=0 go build -o gtodo-api cmd/gtodo/main.go

FROM alpine:latest as prod
WORKDIR /gtodo
COPY --from=build /gtodo/gtodo-api .
COPY --from=build /gtodo/.cert ./.cert
COPY --from=build /gtodo/web ./web
CMD ["./gtodo-api"]