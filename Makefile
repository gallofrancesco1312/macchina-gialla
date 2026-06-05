IMAGE      := macchina-gialla
DEPLOY_DIR := /opt/macchina-gialla

.PHONY: build run test deploy

build:
	docker build -t $(IMAGE) .

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
