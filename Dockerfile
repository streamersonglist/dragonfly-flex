ARG DRAGONFLY_VERSION=1.23.0 
ARG REDIS_VERSION=7.4.0

FROM golang:1.22.6 AS builder

WORKDIR /go/src/github.com/streamersonglist/dragonfly-flex
COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
  go build -o /fly/bin/start ./cmd/start && \
  go build -o /fly/bin/monitor ./cmd/monitor && \
  go build -o /fly/bin/admin_server ./cmd/admin_server

FROM docker.dragonflydb.io/dragonflydb/dragonfly:v${DRAGONFLY_VERSION} AS dragonfly

FROM redis:${REDIS_VERSION} AS redis

FROM otel/opentelemetry-collector-contrib:latest as collector

FROM ubuntu:24.04

ARG VERSION=custom
ARG HAPROXY_VERSION=2.8

LABEL dragonfly-flex.version=${VERSION}
LABEL dragonfly-flex.dragonfly.version=${DRAGONFLY_VERSION}
LABEL dragonfly-flex.redis.version=${REDIS_VERSION}
LABEL dragonfly-flex.haproxy.version=${HAPROXY_VERSION}

# Install DragonflyDB
COPY --from=dragonfly /usr/local/bin/dragonfly /usr/local/bin/dragonfly

# Install Redis Sentinel
COPY --from=redis /usr/local/bin/redis-sentinel /usr/local/bin/redis-sentinel

# Install HAProxy
RUN apt-get update && apt-get install --no-install-recommends -y \
  haproxy=$HAPROXY_VERSION.\* \
  && apt autoremove -y && apt clean

# Install OpenTelemetry Collector
# TODO: need to figure out a better OTEL strategy, this is a hefty binary
COPY --from=collector /otelcol-contrib /otelcol-contrib

# Copy Go binaries from the builder stage
COPY --from=builder /fly/bin/* /usr/local/bin

ADD config/sentinel.conf /fly/sentinel.conf
ADD config/haproxy.cfg /fly/haproxy.cfg
RUN mkdir -p /run/haproxy/

EXPOSE 6379

CMD ["start"]
