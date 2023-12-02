FROM golang:1.21.4-alpine3.17

RUN mkdir /app

COPY ./proxy /app/proxy

WORKDIR /app/proxy

RUN go build -o main .

EXPOSE 5001

CMD ["./main"]