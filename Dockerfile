FROM golang:1.21.4-alpine3.17

RUN mkdir /app

COPY ./proxy /app/proxy

WORKDIR /app

CMD ["go run main.go --server=proxy --port=5001"]