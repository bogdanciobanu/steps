#!/usr/bin/env sh

if ! hash goimports 2>/dev/null; then
  echo "Installing goimports..."
  go get -u golang.org/x/tools/...
fi

if ! hash gotestsum 2>/dev/null; then
  echo "Installing gotestsum..."
  go install gotest.tools/gotestsum@latest
fi

if ! [ -x "$(command -v $TARGET_DIR/golangci-lint)" ]; then
  echo 'Installing golangci-lint' >&2
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $TARGET_DIR v1.39.0
fi