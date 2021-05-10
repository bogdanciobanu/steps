FROM us-docker.pkg.dev/stackpulse/public/base:1.0.0 AS stepsBase
FROM golang:1.14-alpine

COPY --from=stepsBase /usr/local/bin/ /usr/local/bin/

RUN apk add --no-cache curl
# Locking for current stable version (v1.18.6), to get the latest update run `curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kubectl
RUN mv ./kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl