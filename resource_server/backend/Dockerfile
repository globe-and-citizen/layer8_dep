FROM golang:1.21

WORKDIR /usr/src/resource/backend

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/resource/backend/cmd/main resource/backend/cmd/main.go

CMD ["resource/backend/cmd/main"]
