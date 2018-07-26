FROM golang:latest 


EXPOSE 8181

WORKDIR /go/src/
RUN mkdir http_server
ADD . /go/src/http_server
WORKDIR /go/src/http_server
CMD ["go", "run", "main.go"]