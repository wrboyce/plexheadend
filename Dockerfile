FROM golang:1.9.2

ADD *.go /go/src/github.com/wrboyce/plexheadend/
RUN go get -v github.com/wrboyce/plexheadend
CMD plexheadend
