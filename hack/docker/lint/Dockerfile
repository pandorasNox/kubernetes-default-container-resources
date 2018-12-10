FROM golang:1.10.3-alpine3.7 AS compile

# add git
# RUN apk --update add git openssh && \
#     rm -rf /var/lib/apt/lists/* && \
#     rm /var/cache/apk/*

# golangci-lint
# In alpine linux (as it does not come with curl by default)
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN wget https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh && \
    chmod +x install.sh && \
    ./install.sh

# ENTRYPOINT ["/go/bin/gometalinter"]
# CMD ["--help"]

ENTRYPOINT ["golangci-lint"]
CMD ["--help"]
