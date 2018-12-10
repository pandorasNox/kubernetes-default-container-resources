
PWD:=$(shell pwd)

FOO = $(shell docker build ./hack/docker/lint)
LINT_IMAGE ?= $(shell docker build -q ./hack/docker/lint)

.PHONY: lint
lint: ##@testing lints go code
	docker run -it --rm \
	-v $(PWD):/go/src/github.com/pandorasnox/kubernetes-default-container-resources/ \
	-w /go/src/github.com/pandorasnox/kubernetes-default-container-resources/ \
	--entrypoint "/bin/sh" \
	$(LINT_IMAGE) \
	-c "golangci-lint run ./"
