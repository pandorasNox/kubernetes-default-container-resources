
include ./hack/help.mk

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
ifeq (do,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "do"
  DO_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(DO_ARGS):;@:)
endif
.PHONY: do
do: ##@setup reuse minikube docker env
	@eval $$(minikube docker-env) ;\
	docker $(DO_ARGS)


# If the first argument is "run"...
ifeq (build,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  DO_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(DO_ARGS):;@:)
endif
.PHONY: build
build: ##@setup reuse minikube docker env
	@eval $$(minikube docker-env) ;\
	docker build -t pandorasnox/kubernetes-default-container-resources:1.0.0 .

.PHONY: deploy
deploy:
	kubectl apply -f kubernetes/deploy/namespace.yaml \
		-f kubernetes/deploy/ \
		-f kubernetes/MutatingWebhookConfiguration.yaml

.PHONY: undeploy
undeploy:
	kubectl delete -f kubernetes/deploy/namespace.yaml \
		-f kubernetes/MutatingWebhookConfiguration.yaml
