ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR=$(ROOT_DIR)/bin
VERSION=0.1.0
NAME=kitsu

install:
	@(echo "-> Installing requirements...")
	@(go get -u github.com/golang/dep/...)
	@(dep ensure -v)

lint:
	@(gometalinter \
		--vendor \
		--disable-all \
		--enable=errcheck \
		--enable=ineffassign \
		--enable=unused \
		--enable=staticcheck \
		--enable=goimports \
		--sort=path \
		.)

build:
	@(echo "-> Creating binary...")
	@(mkdir -p $(BIN_DIR))
	@(go build -o $(BIN_DIR)/$(NAME))
