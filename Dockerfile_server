FROM golang:1.21.4-alpine3.17

RUN mkdir /build

COPY ./proxy /build

WORKDIR /build

RUN go mod tidy

RUN go build -o main .

EXPOSE 5001

RUN chmod +x ./main

ENTRYPOINT ["./main", "--server=proxy", "--port=5001"]