SHELL := /bin/bash

kuda_root := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
provider_version := $(shell cat $(kuda_root)/providers/$(provider)/VERSION)

build-provider:
	docker build $(kuda_root)/providers/$(provider) \
		-t gcr.io/kuda-project/provider-$(provider):$(provider_version)

run:
	go run \
		-ldflags "-X github.com/cyrildiagne/kuda/cmd.$(provider)ProviderVersion=$(provider_version)" \
		$(kuda_root) $(cmd)

build-provider-and-run: build-provider run
	