FROM golang:1.15-alpine as dev

RUN env GO111MODULE=on go get github.com/cortesi/modd/cmd/modd

COPY . /signalling
WORKDIR /signalling

RUN go mod download
RUN go build -o signalling cmd/signalling/main.go

CMD ["modd"]