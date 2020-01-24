# Use the official go docker image built on debian.
FROM golang:latest

# Install revel and the revel CLI.
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel

ADD . /go/src/github.com/dsmatilla/just-tit-revel

# Use the revel CLI to start up our application.
ENTRYPOINT revel run github.com/dsmatilla/just-tit-revel

# Open up the port where the app is running.
EXPOSE 9000
