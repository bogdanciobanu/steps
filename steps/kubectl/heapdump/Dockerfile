# Globals
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=heapdump
ARG CURRENT_BRANCH
ARG STEP_BASEPATH

# 1. Base
FROM us-docker.pkg.dev/stackpulse/public/kubectl/base:1.0.0 AS base

# 2. Builder
FROM golang:1.16-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BINARY_NAME
ARG BASEDIR

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /${BINARY_NAME}

RUN apk add --no-cache curl
RUN curl -L -o /jattach https://github.com/apangin/jattach/releases/download/v1.5/jattach

# 3. Production
FROM alpine:3.7
ARG BINARY_NAME
COPY --from=builder /${BINARY_NAME} .
COPY --from=builder /jattach .
RUN chmod +x ${BINARY_NAME} /jattach
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENV STEP_NAME=${BINARY_NAME}
ENTRYPOINT ["/bin/sh", "-c", "set -eo pipefail; /${STEP_NAME} 2>&1 | sp-base-step format"]