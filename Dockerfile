FROM golang:1.15.6-alpine AS builder

RUN apk --no-cache add gcc musl-dev

WORKDIR ${GOPATH}/src/github.com/PremoWeb/Chronos
COPY . ${GOPATH}/src/github.com/PremoWeb/Chronos

RUN go build -o /go/bin/chronos .

FROM alpine:3.12

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /go/bin/chronos /usr/bin/chronos

ENTRYPOINT ["/usr/bin/chronos"]

CMD ["daemon", "--config", "/etc/chronos/config.ini"]
