FROM golang:1.17 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /user/go/src
COPY . /user/go/src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-w -s" -o main

#RUN go build .
FROM busybox:1.28.4

WORKDIR /user/go/src
ENV GIN_MODE=release

EXPOSE 8000

ENTRYPOINT ["./main"]