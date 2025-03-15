FROM golang:1.23.2-alpine AS builder

LABEL org.opencontainers.image.description="Cron alternative for Docker Swarm enviornments."

RUN apk --no-cache add gcc musl-dev

WORKDIR ${GOPATH}/src/github.com/PremoWeb/Chadburn
COPY . ${GOPATH}/src/github.com/PremoWeb/Chadburn

RUN go build -o /go/bin/chadburn .

FROM alpine:3.20.3

RUN apk --update --no-cache add ca-certificates tzdata

# Add docker group with same GID as host
RUN addgroup -S -g 969 docker && \
    adduser -S -D -H -h /app -s /sbin/nologin -G docker -u 1000 appuser

COPY --from=builder /go/bin/chadburn /usr/bin/chadburn

# Set permissions
RUN chmod +x /usr/bin/chadburn

# Use the appuser
USER appuser

ENTRYPOINT ["/usr/bin/chadburn"]

CMD ["daemon"]