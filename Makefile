SHELL := /bin/bash

kuda_root := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
provider_version := $(shell cat $(kuda_root)images/providers/$(provider)/VERSION)

build-provider:
	docker build $(kuda_root)images/providers/$(provider) \
		-t gcr.io/kuda-project/provider-$(provider):$(provider_version)

run:
	go run $(kuda_root) $(cmd)

build-provider-and-run: build-provider run
	