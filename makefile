SHELL := /bin/bash
KIND_CLUSTER := shrt-service-cluster

run:
	go run main.go

docker-image:
	docker build \
			-f zarf/docker/Dockerfile \
			-t shrt-service-arm64:$(shell git rev-parse --short HEAD) \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
			--build-arg BUILD_REF=$(shell git rev-parse --short HEAD) \
			.

kind-up:
	kind create cluster \
			--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
			--name $(KIND_CLUSTER) \
			--config zarf/k8s/kind/kind-config.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
