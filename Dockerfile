# gcr.io/distroless/static-debian11

FROM golang:1.21.0-alpine AS builder

ARG ALPINE_TAG=latest

ARG ARGOBOT_VERSION=dev
ENV ARGOBOT_VERSION=${ARGOBOT_VERSION}
ARG ARGOBOT_COMMIT=none
ENV ARGOBOT_COMMIT=${ARGOBOT_COMMIT}

WORKDIR /app

RUN apk add --no-cache \
        bash~=5.2 \
        curl

RUN curl -sSL -o /tmp/argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64

COPY go.mod go.sum ./
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

# RUN --mount=type=cache,target=/go/pkg/mod \
#     --mount=type=cache,target=/root/.cache/go-build \
#     CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'main.version=${ARGOBOT_VERSION}' -X 'main.commit=${ARGOBOT_COMMIT}'" -v -o argobot .
COPY . /app

RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'main.version=${ARGOBOT_VERSION}' -X 'main.commit=${ARGOBOT_COMMIT}'" -v -o argobot ./cmd/argobot

FROM alpine:latest AS alpine

RUN addgroup argobot && \
    adduser -S -G argobot argobot && \
    adduser argobot root && \
    chown argobot:root /home/argobot/ && \
    chmod g=u /home/argobot/ && \
    chmod g=u /etc/passwd

COPY --from=builder /app/argobot /usr/local/bin/argobot
# COPY ./bin/argocd-linux-amd64 /usr/local/bin/argocd
COPY --from=builder /tmp/argocd-linux-amd64 /tmp/argocd-linux-amd64
RUN install -m 555 /tmp/argocd-linux-amd64 /usr/local/bin/argocd && rm -f /tmp/argocd-linux-amd64
# RUN curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64 \
#     sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd \
#     rm argocd-linux-amd64

RUN apk add --no-cache \
        ca-certificates~=20230506 \
        git~=2.40

RUN mkdir /app && \
    chown argobot:argobot /app

WORKDIR /app

USER argobot

CMD ["argobot", "start"]
