FROM golang:1.21.4-alpine3.17

RUN mkdir /app

COPY ./proxy /app

WORKDIR /app

RUN go mod tidy

RUN go build -o main .

EXPOSE 5001

CMD ["./main --server=proxy --port=5001"]