IMAGE := ghcr.io/pch/mediator
IMAGE_DEV := $(IMAGE):dev
VERSION ?= $(shell cat VERSION | tr -d '\n')
PORT ?= 8000

build:
	docker buildx build --push --target production --platform linux/amd64,linux/arm64 --tag ${IMAGE}:${VERSION} --tag ${IMAGE}:latest .

build-dev:
	docker buildx build --target development --tag ${IMAGE_DEV} --progress plain .

run-dev:
	docker run --rm -it --env-file .env \
	  --volume .:/app \
		--volume gocache:/go \
		--publish ${PORT}:${PORT} ${IMAGE_DEV}

run:
	docker run --rm -it --env-file .env --publish ${PORT}:${PORT} ${IMAGE}:${VERSION}
