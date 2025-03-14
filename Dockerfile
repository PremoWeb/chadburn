FROM golang:1.19-alpine AS builder

LABEL org.opencontainers.image.description="Cron alternative for Docker Swarm enviornments."

RUN apk --no-cache add gcc musl-dev

WORKDIR ${GOPATH}/src/github.com/PremoWeb/Chadburn
COPY . ${GOPATH}/src/github.com/PremoWeb/Chadburn

RUN go build -o /go/bin/chadburn .

FROM alpine:3.20.3

RUN apk --update --no-cache add ca-certificates tzdata

COPY --from=builder /go/bin/chadburn /usr/bin/chadburn

ENTRYPOINT ["/usr/bin/chadburn"]

CMD ["daemon"]