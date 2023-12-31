HOSTNAME=dantata.com
NAMESPACE=cloud
NAME=cloudfusion
BINARY=terraform-provider-${NAME}
VERSION=1.0.0
GOARCH  := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}

test:
	go test ./... -v

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
