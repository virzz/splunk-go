LDFLAGS=
PACKAGE_NAME=./cmd
TARGET_NAME=splunk-go
MODULE=github.com/virzz/splunk-go
Version=${DRONE_TAG}

ifeq ($(Version),)
	Version=$(shell git describe --tags --abbrev=0 || echo development )
endif
ifeq ($(Commit),)
	Commit=$(shell git rev-parse HEAD 2>/dev/null )
endif
ifeq ($(TARGET_NAME),unknown)
	$(error TARGET_NAME is not set)
endif

default:
	echo "make target"

linux:
	GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags="-s -w ${LDFLAGS} \
		-X ${MODULE}/config.Version=${Version} \
		-X ${MODULE}/config.Commit=${Commit}" \
		-o build/${TARGET_NAME}_linux_amd64 ./cmd/;

release: linux
	cp build/${TARGET_NAME}_linux_amd64 build/${TARGET_NAME}
	upx -9 build/${TARGET_NAME}
