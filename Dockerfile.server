FROM golang

RUN mkdir -p /go/src/github.com/bfosberry/bonano
ADD . /go/src/github.com/bfosberry/bonano
WORKDIR /go/src/github.com/bfosberry/bonano

RUN go get ./...

RUN go build -o server/server ./server/main.go

ENTRYPOINT server/server
