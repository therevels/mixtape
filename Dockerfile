FROM golang:1.10.3

EXPOSE 8088

ENV DEP_VERSION=v0.5.0

RUN apt-get update && \
    apt-get install -y curl && \
    curl -LO https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && \
    chmod a+x dep-linux-amd64 && \
    mv dep-linux-amd64 /usr/local/bin/dep && \
    rm -rf /var/lib/apt/lists/*

# We don't absolutely NEED the ginkgo binary to run the test suite, but it's
# nice to have
RUN go get -u github.com/onsi/ginkgo/ginkgo
RUN go get -u github.com/onsi/gomega/...

WORKDIR /go/src/github.com/therevels/mixtape

# Generate a self-signed cert for development
RUN go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock
RUN dep ensure -vendor-only

COPY . .

# This is for development purposes only. Eventually we'll want CI/CD to build
# a leaner release image
RUN go build -o dist/mixtape ./...

CMD "dist/mixtape"
