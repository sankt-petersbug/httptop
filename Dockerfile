FROM    golang:1.10.3

RUN     go get github.com/sankt-petersbug/httptop/cmd/httptop
WORKDIR /go/src/github.com/sankt-petersbug/httptop
cmd httptop
