# Globals
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=cmd
ARG CURRENT_BRANCH
ARG STEP_BASEPATH

# 1. Base
FROM us-docker.pkg.dev/stackpulse/public/kubectl/base:1.0.0 AS base

# 2. Builder
FROM golang:1.14-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BINARY_NAME
ARG BASEDIR

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /${BINARY_NAME}

# 3. Production
FROM alpine:3.7
ARG BINARY_NAME
ARG STEP_BASEPATH
COPY --from=builder /${BINARY_NAME} .
RUN chmod +x ${BINARY_NAME}
COPY --from=base /usr/local/bin/ /usr/local/bin/
COPY ${STEP_BASEPATH}/entrypoint.sh /entrypoint.sh

ENV STEP_NAME=${BINARY_NAME}
ENTRYPOINT ["/entrypoint.sh"]