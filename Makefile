# Makefile

export GOPATH := $(shell pwd)

all:
	echo $$GOPATH

build:
	go build -o threader

install:
	cp threader /usr/bin
