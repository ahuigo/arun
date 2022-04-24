FROM golang:1.14

MAINTAINER Ahuigo <ahuigo@qq.com>

ENV GOPATH /go
ENV GO111MODULE on

COPY . /go/src/github.com/ahuigo/arun
WORKDIR /go/src/github.com/ahuigo/arun
RUN make install

ENTRYPOINT ["/go/bin/arun"]
