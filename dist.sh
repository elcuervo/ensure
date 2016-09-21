#!/usr/bin/env sh

for os in $(echo darwin linux); do
  for arch in $(echo amd64 386); do
    echo "==> Creating ${os}_${arch} build."

    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o dist/ensure_${os}_${arch} .
  done
done
