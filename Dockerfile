FROM golang:1.17.1-alpine AS builder

RUN apk --no-cache add gcc musl-dev

WORKDIR ${GOPATH}/src/github.com/PremoWeb/Chronos
COPY . ${GOPATH}/src/github.com/PremoWeb/Chronos

RUN go build -o /go/bin/chronos .

FROM alpine:3.14.2

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /go/bin/chronos /usr/bin/chronos

ENTRYPOINT ["/usr/bin/chronos"]

CMD ["daemon"]
