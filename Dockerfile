FROM golang:1.10.3

EXPOSE 8088

ENV DEP_VERSION=v0.5.0

RUN apt-get update && \
    apt-get install -y curl && \
    curl -LO https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && \
    chmod a+x dep-linux-amd64 && \
    mv dep-linux-amd64 /usr/local/bin/dep && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/therevels/mixtape

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock
RUN dep ensure -vendor-only

COPY . .

# We don't absolutely NEED the ginkgo binary to run the test suite, but it's
# nice to have
go get -u github.com/onsi/ginkgo/ginkgo
go get -u github.com/onsi/gomega/...

# This is for development purposes only. Eventually we'll want CI/CD to build
# a leaner release image
RUN go build -o dist/mixtape ./...

CMD "dist/mixtape"
