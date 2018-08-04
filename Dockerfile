FROM golang:1.10.3

EXPOSE 8088

WORKDIR /go/src/github.com/therevels/mixtape
COPY . .

RUN go build -o dist/mixtape ./...

CMD "dist/mixtape"
