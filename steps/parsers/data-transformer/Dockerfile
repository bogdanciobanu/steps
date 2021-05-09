FROM us-docker.pkg.dev/stackpulse/public/parsers/base:1.0.0

RUN apk add --no-cache unzip curl wget

# Install JQ, the JSON Parser
RUN apk add --no-cache jq

# Install pup, the XML/HTML Parser
ARG PUP_VERSION=v0.4.0
RUN wget "https://github.com/ericchiang/pup/releases/download/$PUP_VERSION/pup_${PUP_VERSION}_linux_amd64.zip" && \
    unzip "pup_${PUP_VERSION}_linux_amd64.zip" && \
    rm "pup_${PUP_VERSION}_linux_amd64.zip" && \
    chmod +x ./pup && \
    mv ./pup /usr/local/bin/pup

# Install faq, the format-agnostic libjq-based parser tool - 
ARG LATEST_RELEASE=0.0.6
RUN curl -Lo /usr/local/bin/faq https://github.com/jzelinskie/faq/releases/download/$LATEST_RELEASE/faq-linux-amd64
RUN chmod +x /usr/local/bin/faq

# Start Entry Point
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT /entrypoint.sh
