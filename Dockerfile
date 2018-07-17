FROM golang:1.10.3-alpine3.7 AS compile

# add git
RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY . /go/src/github.com/pandorasnox/kubernetes-default-container-resources/
WORKDIR /go/src/github.com/pandorasnox/kubernetes-default-container-resources/
RUN go get .
RUN go install .

# ============================================================
# ============================================================

FROM alpine:3.8
COPY --from=compile /go/bin/kubernetes-default-container-resources /kubernetes-default-container-resources
ENTRYPOINT ["/kubernetes-default-container-resources"]
