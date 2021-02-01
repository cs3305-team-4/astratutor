FROM golang:1.15-alpine as dev

RUN env GO111MODULE=on go get github.com/cortesi/modd/cmd/modd

COPY . /api
WORKDIR /api

RUN go mod download
RUN go build -o server cmd/server/main.go

CMD ["modd"]