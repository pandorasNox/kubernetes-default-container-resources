
include ./hack/help.mk
include ./hack/lint.mk

UID:=$(shell id -u)
GID:=$(shell id -g)
PWD:=$(shell pwd)

# .PHONY: setup
# setup: ##@setup builds the docker image for the cli make cmd
#  docker ...

.PHONY: cli
cli: ##@setup set up a docker container with mounted source where you can execute all go commands
	# docker run -it --rm -u $(UID):$(GID) -v $(PWD):/source -w /source golang:1.10.3 bash
	docker run -it --rm -v $(PWD):/go/src/k8s-resource-mutator -w /go/src/k8s-resource-mutator -v $(PWD)/certs:/certs -p 8083:8083 golang:1.10.3 bash

# If the first argument is "do"...
ifeq (mini-do,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "do"
  DO_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(DO_ARGS):;@:)
endif
.PHONY: mini-do
mini-do: ##@minikube reuse minikube docker env to run docker cmd's | e.g. `make mini-do ps`
	@eval $$(minikube docker-env) ;\
	docker $(DO_ARGS)

# If the first argument is "run"...
ifeq (build,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  DO_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(DO_ARGS):;@:)
endif
.PHONY: mini-build
mini-build: ##@minikube build docker img in minikube env
	@eval $$(minikube docker-env) ;\
	docker build -t pandorasnox/kubernetes-default-container-resources:1.1.0 .

.PHONY: deploy
deploy: ##@setup deploys the webhook server + mutate config to the current kubernetes ctx
	kubectl apply -f kubernetes/deploy/namespace.yaml \
		-f kubernetes/deploy/ \
		-f kubernetes/MutatingWebhookConfiguration.yaml

.PHONY: undeploy
undeploy: ##@setup undeploy the mutate server webhook (based on current kubernetes ctx)
	kubectl delete -f kubernetes/deploy/namespace.yaml \
		-f kubernetes/MutatingWebhookConfiguration.yaml
	./hack/wait-until-ns-deleted.sh

.PHONY: mini-clear-intermediate
mini-clear-intermediate: ##@minikube deletes all intermediate docker images on minikube k8s cluster
	@eval $$(minikube docker-env) ;\
	docker rmi -f $$(docker images --filter dangling=true -q)

.PHONY: test-deploy (based on current kubernetes ctx)
test-deploy: ##@dev deploys a test example
	kubectl apply -f kubernetes/example/namespace.yaml \
		-f kubernetes/example/

.PHONY: test-undeploy (based on current kubernetes ctx)
test-undeploy: ##@dev unddeploys a test example
	kubectl delete -f kubernetes/example/namespace.yaml

.PHONY: test
test: ##@testing runs all go tests
	go test -v ./pkg/

.PHONY: bats
bats: ##@testing runs all bats tests
	bats ./tests/bats/apply-pod-without-resources.bats

.PHONY: run
run: ##@dev
	go run main.go -sslCert "./certs/ssl-cert.pem" -sslKey "./certs/ssl-key.pem"

.PHONY: run-no-tls
run-no-tls: ##@dev
	go run main.go -tlsDisabled true

.PHONY: certs
certs: ##@dev
	rm -rf certs/*
	docker run --rm -v $(PWD)/certs:/certs -e SSL_SUBJECT=default-container-resources.mutating-webhook.svc -e SSL_KEY="ssl-key.pem" -e SSL_CSR="ssl-key.csr" -e SSL_CERT="ssl-cert.pem" -e K8S_NAME="tls-cert-default-container-resources" -e K8S_NAMESPACE="mutating-webhook" -e CA_EXPIRE=3650 paulczar/omgwtfssl
	docker run --rm -v $(PWD)/certs:/certs alpine:3.8 sh -c "chmod -R a+rw /certs"
