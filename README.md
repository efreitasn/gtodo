# gtodo
This is a todo web app written in go with no async calls in the frontend. The purpose of this was to create a UI with CRUD operations by using only HTTP requests through form submissions. It uses MongoDB to store the todos and user accounts and HTTP/2 to serve the files.

## Why
This project was created to learn more about the [net/http](https://godoc.org/net/http) and [html/template](https://godoc.org/html/template) go packages. Besides, it was also used to learn how to use the [MongoDB go driver](https://github.com/mongodb/mongo-go-driver).

## How to run
After cloning this repo, you can use docker-compose to run the project:

```bash
docker-compose up
```

## SSL certificates
As said, this project uses HTTP/2, which means an SSL certificate will be needed to run the project. You can create your certificate using something like [mkcert](https://github.com/FiloSottile/mkcert). It's expected that both the certificate and the private key are in the `./.cert` directory with the respective names `cert.pem` and `key.pem`.