
.PHONY: build
build:
	go build ./...

EXPORTER_DOCKERFILE ?= Dockerfile
EXPORTER_IMAGE_TAG ?= $(shell cat ${EXPORTER_DOCKERFILE}  | sha256sum | cut -d" " -f 1 | cut -c1-7)
EXPORTER_IMAGE ?= supershal/gha-exporter:$(EXPORTER_IMAGE_TAG)

.PHONY: docker-build
docker-build:
	docker build --file=$(EXPORTER_DOCKERFILE) -t $(EXPORTER_IMAGE) $(dir $(EXPORTER_DOCKERFILE))

.PHONY: docker-push
docker-push: docker-build
	docker push $(EXPORTER_IMAGE)

deploy:
	helm upgrade gha-exporter gha-exporter -f gha-exporter/values.yaml --set github_token=$(GITHUB_TOKEN) --set image.tag=$(EXPORTER_IMAGE_TAG)