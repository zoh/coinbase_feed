FROM golang:1.14


RUN go version
# Create app directory
ADD . /usr/src/app

WORKDIR /usr/src/app

RUN go build -o feeder ./cmd/feeder.go  && ls -lah