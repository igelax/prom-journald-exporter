NAME := prom-journald-exporter
MAINTAINER := Mike Sgarbossa <mikeysky@gmail.com>
DESCRIPTION := Service for exporting metrics derived from journald to Prometheus
LICENSE := MIT

VERSION := $(shell cat VERSION)
OUT := out

clean:
	rm -rf $(OUT)

build_amd64:
	mkdir -p $(OUT)
	GOARCH=amd64 GOOS=linux go build -v -o '$(OUT)/$(NAME)'
	tar -C $(OUT) -zcf $(OUT)/$(NAME)-$(VERSION)-linux-amd64.tar.gz $(NAME)
	rm $(OUT)/$(NAME)

build_arm64:
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOARCH=arm64 GOOS=linux go build -v -o '$(OUT)/$(NAME)'
	tar -C $(OUT) -zcf $(OUT)/$(NAME)-$(VERSION)-linux-arm64.tar.gz $(NAME)
	rm $(OUT)/$(NAME)

build_arm:
	CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc GOARCH=arm GOOS=linux go build -v -o '$(OUT)/$(NAME)'
	tar -C $(OUT) -zcf $(OUT)/$(NAME)-$(VERSION)-linux-arm.tar.gz $(NAME)
	rm $(OUT)/$(NAME)

all: clean build_amd64 build_arm64 build_arm
