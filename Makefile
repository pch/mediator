IMAGE := ghcr.io/pch/mediator
IMAGE_DEV := $(IMAGE):dev
VERSION ?= $(shell cat VERSION | tr -d '\n')
PORT ?= 8000

build:
	docker buildx build --target production --tag ${IMAGE}:${VERSION} .

build-dev:
	docker buildx build --target development --tag ${IMAGE_DEV} --progress plain .

push:
	docker push ${IMAGE}:${VERSION}

release:
	git tag -a "v${VERSION}" -m "Release ${VERSION}"
	docker pull ${IMAGE}:${VERSION}
	docker tag  ${IMAGE}:${VERSION} ${IMAGE}:latest
	docker push ${IMAGE}:latest

run-dev:
	docker run --rm -it --env-file .env \
	  --volume .:/app \
		--volume gocache:/go \
		--publish ${PORT}:${PORT} ${IMAGE_DEV}

run:
	docker run --rm -it --env-file .env --publish ${PORT}:${PORT} ${IMAGE}:${VERSION}
