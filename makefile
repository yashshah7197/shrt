SHELL := /bin/bash

# ======================================================================================================================
# Standard Go stuff
# ======================================================================================================================

run shrt-api:
	go run app/services/shrt-api/main.go | go run app/tooling/logfmt/main.go

shrt-admin:
	go run app/tooling/admin/main.go

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
			-f zarf/docker/dockerfile.shrt-api \
			-t shrt-api-arm64:$(VERSION) \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
			--build-arg BUILD_REF=$(VERSION) \
			.

# ======================================================================================================================
# KIND & K8s
# ======================================================================================================================
KIND_CLUSTER := shrt-api-cluster

kind-load:
	cd zarf/k8s/kind/shrt-api-pod; kustomize edit set image shrt-api-image=shrt-api-arm64:$(VERSION)
	kind load docker-image shrt-api-arm64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/shrt-api-pod | kubectl apply -f -

kind-update: docker-image kind-load kind-restart

kind-update-apply: docker-image kind-load kind-apply

kind-up:
	kind create cluster \
			--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
			--name $(KIND_CLUSTER) \
			--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=shrt-api-system

kind-setup: kind-up docker-image kind-load kind-apply

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-restart:
	kubectl rollout restart deployment shrt-api-pod

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-shrt-api:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=shrt-api --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-describe:
	kubectl describe pod -l app=shrt-api

# ==================================================================================================
# ExpvarMon
# ==================================================================================================

expvarmon:
	~/go/bin/expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# ==================================================================================================
# Hey
# ==================================================================================================

hey:
	hey -m GET -c 100 -n 10000 http://localhost:3000/test

# ==================================================================================================
# Generate a private/public PEM keypair.
# ==================================================================================================

genkeypair:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private.pem -out public.pem
