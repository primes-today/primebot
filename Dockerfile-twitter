FROM golang:alpine

RUN apk --no-cache add ca-certificates && update-ca-certificates

RUN mkdir -p /go/src/github.com/fardog/primebot
COPY . /go/src/github.com/fardog/primebot

WORKDIR /go/src/github.com/fardog/primebot/cmd/primebot-twitter
RUN go install -v
WORKDIR /go/bin
RUN rm -rf /go/src/github.com/fardog/primebot

CMD primebot-twitter --interval 1h
