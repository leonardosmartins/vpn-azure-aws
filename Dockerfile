FROM golang:alpine

RUN apk update && \
    apk add --no-cache git
RUN mkdir -p /golang
ADD connection.go /golang

WORKDIR /golang
RUN go get github.com/go-sql-driver/mysql
RUN go build connection.go

CMD ./connection && while true; do sleep 300; done;