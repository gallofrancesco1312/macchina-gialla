IMAGE := macchina-gialla

.PHONY: build run test

build:
	docker build -t $(IMAGE) .

run:
	docker run --env-file .env -v $(PWD)/counter.json:/app/counter.json $(IMAGE)

test:
	go test ./...
