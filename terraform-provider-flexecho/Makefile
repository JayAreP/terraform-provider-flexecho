# flexecho provider build vars. the build-local scripts read these so nothing drifts
NAME=flexecho
BINARY=terraform-provider-${NAME}
VERSION=0.1.1
HOSTNAME=localdomain
NAMESPACE=provider
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: build

# host build
build:
	go build -o ${BINARY}

# install into the local plugin dir so terraform PIcks it up via localdomain/provider/flexecho
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test ./...

vet:
	go vet ./...
