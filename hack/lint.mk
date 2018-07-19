
PWD:=$(shell pwd)

FOO = $(shell docker build ./hack/docker/lint)
LINT_IMAGE ?= $(shell docker build -q ./hack/docker/lint)

.PHONY: lint
lint:
	docker run -it --rm \
	-v $(PWD):/go/src/github.com/pandorasnox/kubernetes-default-container-resources/ \
	-w /go/src/github.com/pandorasnox/kubernetes-default-container-resources/ \
	--entrypoint "/bin/sh" \
	$(LINT_IMAGE) \
	-c "go get -v .; echo ''; echo 'lint results:'; gometalinter ./... | grep -v '../../../../../usr/local'"
