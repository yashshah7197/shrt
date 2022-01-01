SHELL := /bin/bash

# ======================================================================================================================
# Standard Go stuff
# ======================================================================================================================

run:
	go run main.go

tidy:
	go mod tidy
	go mod vendor

staticcheck:
	staticcheck -checks=all ./...

# ======================================================================================================================
# Docker
# ======================================================================================================================
VERSION := 1.0

docker-image:
	docker build \
			-f zarf/docker/Dockerfile \
			-t shrt-service-arm64:$(VERSION) \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
			--build-arg BUILD_REF=$(VERSION) \
			.

# ======================================================================================================================
# KIND & K8s
# ======================================================================================================================
KIND_CLUSTER := shrt-service-cluster

kind-load:
	kind load docker-image shrt-service-arm64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/shrt-service-pod | kubectl apply -f -

kind-update: docker-image kind-load kind-restart

kind-update-apply: docker-image kind-load kind-apply

kind-up:
	kind create cluster \
			--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
			--name $(KIND_CLUSTER) \
			--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=shrt-service-system

kind-setup: kind-up docker-image kind-load kind-apply

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-restart:
	kubectl rollout restart deployment shrt-service-pod

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-shrt-service:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=shrt-service --all-containers=true -f --tail=100

kind-describe:
	kubectl describe pod -l app=shrt-service
