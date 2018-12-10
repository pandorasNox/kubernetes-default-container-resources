 


.PHONY: ci-mini-init
ci-mini-init:
	./hack/minikube/ci-start-minikube.sh
	minikube update-context

.PHONY: ci-mini-await
ci-mini-await:
	./hack/minikube/ci-await-minikube.sh

# .PHONY: logs
# logs: ##@setup Shows logs.
# 	ktail -n container-image-builder
