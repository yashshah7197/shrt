SHELL := /bin/bash

run:
	go run main.go

docker-image:
	docker build \
			-f zarf/docker/Dockerfile \
			-t shrt-arm64:$(shell git rev-parse --short HEAD) \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
			--build-arg BUILD_REF=$(shell git rev-parse --short HEAD) \
			.
