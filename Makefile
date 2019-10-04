SHELL := /bin/bash

kuda_root := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build-provider-and-run:
	docker build $(kuda_root)images/providers/$(provider) \
		-t gcr.io/kuda-project/provider-$(provider)
	go run $(kuda_root) $(cmd)