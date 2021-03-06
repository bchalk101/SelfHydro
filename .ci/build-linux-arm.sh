#!/bin/bash

set -e -x

GO_ENV=(
	CGO_ENABLED=1
)

GO_CROSS_ENV=(
  GOOS=linux
	GOARCH=arm
	GOARM=7
	CGO_ENABLED=1
	CC=arm-linux-gnueabihf-gcc
)

apt-get update
apt-get install crossbuild-essential-armhf -y

export GOPATH=$PWD

echo "v$(cat version/version)" > release/name
echo "v$(cat version/version)" > release/tag
echo "{ \"TAG\": \"v$(cat version/version)\" }" > release/docker-args

cat > release/body <<EOF
Selfhydro release
EOF

mkdir -p src/github.com/selfhydro/

cp -R ./selfhydro src/github.com/selfhydro/.

cd src/github.com/selfhydro/selfhydro

go get
env ${GO_CROSS_ENV[@]} go build -o $GOPATH/release/selfhydro



