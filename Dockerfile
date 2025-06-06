# gcr.io/distroless/static-debian11

FROM golang:1.24.4-alpine AS builder

ARG ALPINE_TAG=latest

ARG ARGOBOT_VERSION=dev
ENV ARGOBOT_VERSION=${ARGOBOT_VERSION}
ARG ARGOBOT_COMMIT=none
ENV ARGOBOT_COMMIT=${ARGOBOT_COMMIT}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'main.version=${ARGOBOT_VERSION}' -X 'main.commit=${ARGOBOT_COMMIT}'" -v -o argobot ./cmd/argobot

FROM alpine:latest AS base

RUN addgroup argobot && \
    adduser -S -G argobot argobot && \
    adduser argobot root && \
    chown argobot:root /home/argobot/ && \
    chmod g=u /home/argobot/ && \
    chmod g=u /etc/passwd

RUN mkdir /app && \
    chown argobot:argobot /app

WORKDIR /app

USER root

COPY --from=builder /app/argobot /tmp/argobot
RUN install -m 555 /tmp/argobot /usr/local/bin/argobot && rm -f /tmp/argobot

USER argobot

CMD ["argobot", "start"]
