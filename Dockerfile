FROM golang:1.9.2

ADD *.go /go/src/github.com/papaschloss/plexheadend/
RUN go get -v github.com/papaschloss/plexheadend
CMD plexheadend
