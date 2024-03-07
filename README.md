# Mediator

(This README is a work in progress)

## Installation

Mediator is distributed as a docker image:

```shell
docker run --rm -it --env-file .env --publish 8000:8000 ghcr.io/pch/mediator:latest
```

## Configuration

Mediator can be configured using `ENV` variables:

| Variable | Description | Default |
| --- | --- | --- |
| `SOURCES` | **Required**. List of supported sources to pull images from. A comma-separated list of `key=values`. Example:<br>`mybucket=https://mybucket.s3.amazonaws.com,mybucket-dev=https://mybucket-dev.s3.amazonaws.com` | — |
| `SECRET_KEY` | Optional, but **highly encouraged**. Secret random key, used to generate URL signatures. Can be ganerated with:<br> `dd if=/dev/urandom bs=32 count=1 2>/dev/null \| base64 \| tr -d '='` | `""` |
| `AUTH_TOKEN` | Optional. Token for authenticating image requests (for `Authorization: Bearer <AUTH_TOKEN>`). You can set it when configuring your CDN to prevent direct access to the service. | `""` |
| `PATH_PREFIX` | Optional. Prefix for the image processing URL paths. Useful when you have a CDN pulling from multiple sources. <br>Example: `PATH_PREFIX=/my-prefix` will change the `transform` URL to `/my-prefix/image/transform/:source/:path`. | `""` |
| `CACHE_CONTROL` | Value for the `Cache-Control` header. | `public, max-age=31536000` |
| `DOWNLOAD_MAX_SIZE` | File size download limit, in bytes. | `50MB` |
| `DOWNLOAD_TIMEOUT` | Download timeout, in seconds. | `10s` |
| `HTTP_PORT` | HTTP port for the service. | `8000` |
