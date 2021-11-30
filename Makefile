NAME := prom-journald-exporter
MAINTAINER := Mike Sgarbossa <mikeysky@gmail.com>
DESCRIPTION := Service for exporting metrics derived from journald to Prometheus
LICENSE := MIT

VERSION := $(shell cat VERSION)
OUT := out

clean:
	rm -rf $(OUT)

build:
	rm -rf $(OUT)
	mkdir -p $(OUT)
	GOARCH=amd64 GOOS=linux go build -v -o '$(OUT)/$(NAME)'
	tar -C $(OUT) -zcf $(OUT)/$(NAME)-$(VERSION)-linux-amd64.tar.gz $(NAME)
	rm $(OUT)/$(NAME)
