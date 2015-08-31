FROM golang

RUN mkdir -p /go/src/github.com/bfosberry/bonano
ADD . /go/src/github.com/bfosberry/bonano
WORKDIR /go/src/github.com/bfosberry/bonano

RUN go get ./...

RUN go build -o client/client ./client/main.go

ENTRYPOINT client/client
