IMAGE          := macchina-gialla
REGISTRY_HOST  := forgejo.home.potato
REGISTRY_USER  ?= ciccio
REGISTRY_TOKEN ?=
REGISTRY       := $(REGISTRY_HOST)/ciccio/macchina-gialla
TAG            := latest
DEPLOY_DIR := /opt/services/macchina-gialla

.PHONY: build run test deploy push login

login:
	@echo "$(REGISTRY_TOKEN)" | docker login $(REGISTRY_HOST) -u $(REGISTRY_USER) --password-stdin

build:
	docker build -t $(IMAGE) .

push: build
	docker tag $(IMAGE) $(REGISTRY):$(TAG)
	docker push $(REGISTRY):$(TAG)

run:
	docker run --env-file .env -v $(PWD)/counter.json:/app/counter.json $(IMAGE)

test:
	go test ./...

deploy: build
	docker stop $(IMAGE) 2>/dev/null || true
	docker rm   $(IMAGE) 2>/dev/null || true
	docker run -d \
		--name $(IMAGE) \
		--restart unless-stopped \
		--env-file $(DEPLOY_DIR)/.env \
		-v $(DEPLOY_DIR)/counter.json:/app/counter.json \
		$(IMAGE)
