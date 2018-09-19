FROM golang:1.10
RUN go get -v github.com/rocymp/atx-server && cd $GOPATH/src/github.com/rocymp/atx-server && go build

FROM debian:stretch
WORKDIR /root/
COPY --from=0 /go/src/github.com/rocymp/atx-server .
ENTRYPOINT ./atx-server --port 8000
