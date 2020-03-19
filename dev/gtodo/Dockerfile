FROM golang:1.14-alpine3.11

WORKDIR /app
RUN wget https://github.com/efreitasn/wrun/releases/download/v4.0.0/wrun-v4.0.0-linux-amd64 -O /bin/wrun
RUN chmod u+x /bin/wrun

CMD wrun start