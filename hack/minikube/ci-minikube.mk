 


.PHONY: ci-mini-init
ci-mini-init:
	./hack/minikube/start-minikube.sh
	minikube update-context

.PHONY: ci-mini-await
ci-mini-await:
	./hack/minikube/await-minikube.sh

# .PHONY: logs
# logs: ##@setup Shows logs.
# 	ktail -n container-image-builder
