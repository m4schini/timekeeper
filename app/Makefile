CONTAINER_TOOL ?= podman
IMAGE           = codeberg.org/aur0ra/raumzeitalpaka:latest
CONTAINERFILE   = Containerfile

.PHONY: build push all

all: build push

build: generate
	$(CONTAINER_TOOL) build -t $(IMAGE) -f $(CONTAINERFILE) .

push:
	$(CONTAINER_TOOL) push $(IMAGE)

ports/api/api.go: raumzeitalpaka.openapi.yaml ports/api/config.yaml
	go generate ./ports/api/...
	@touch $@

.PHONY: generate
generate: ports/api/api.go

.PHONY: deps
deps:
	$(GO) install tool