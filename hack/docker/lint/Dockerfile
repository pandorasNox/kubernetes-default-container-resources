FROM golang:1.10.3-alpine3.7 AS compile

# add git
RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

RUN go get -u github.com/alecthomas/gometalinter && \
    gometalinter --install

ENTRYPOINT ["/go/bin/gometalinter"]
CMD ["--help"]
