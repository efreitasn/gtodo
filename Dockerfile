FROM golang:1.12.9 as build
WORKDIR /gtodo
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ../gtodo-api cmd/gtodo/main.go
RUN cd .. && go get github.com/efreitasn/wrun

FROM alpine:latest
COPY --from=build /gtodo .
CMD ["../gtodo-api"]