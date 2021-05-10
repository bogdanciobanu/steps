# Globals
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=get-file
ARG BASE_BRANCH
ARG STEP_BASEPATH

# 1. Base
FROM us-docker.pkg.dev/stackpulse/public/base:$BASE_BRANCH AS base

# 2. Builder
FROM golang:1.15-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BINARY_NAME
ARG BASEDIR

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /${BINARY_NAME}

# 3. Production
FROM alpine

ARG BINARY_NAME
COPY --from=builder /${BINARY_NAME} .
RUN chmod +x ${BINARY_NAME}
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENV STEP_NAME=${BINARY_NAME}
ENTRYPOINT ["/bin/sh", "-c", "set -eo pipefail; /${STEP_NAME} 2>&1 | sp-base-step format"]