FROM golang:1.21

WORKDIR /usr/src/layer8

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/layer8 ./main.go

CMD ["layer8", "--server", "auth"]
