FROM golang:latest 
WORKDIR /go/src/
RUN mkdir http_server
ADD . /go/src/http_server
WORKDIR /go/src/http_server