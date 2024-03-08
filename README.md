# Mediator

(this readme is a work in progress)

Mediator is a small, standalone service for efficient image processing. Images are processed on the fly with `libvips` and cached by CDN. The service is intended as a replacement for all common image operations in your web apps.

You can use it to:

- **Resize images.** Generate thumbnails and responsive images on the fly.
- **Crop images**. crop images (currently only smart crop is supported).
- **Apply effects**. Apply filters to images (currently only `pixelate` is supported).
- **Proxy media files**: proxy images and other files behind secure, signed URLs.
- **Strip metadata**. Remove metadata from images to reduce file size and protect user privacy.

Examples:

## Features

- **Fast**. Written in Go, uses `libvips`, for image processing.
- **Simple**. A small web service with minimal configuration.
- **Secure**. Supports signed URLs and authentication tokens.
- **Easy to integrate**. Works with any CDN and storage provider.
- **Easy to deploy**. Just use the provided Docker image.
- **Most common image operations**. Resize, crop, apply effects, and more.
- **Auto WebP**. Auto-convert images to WebP and serve them to browsers that support it.

## Installation

Mediator is distributed as a docker image:

```shell
docker run --rm -it --env-file .env --publish 8000:8000 ghcr.io/pch/mediator:latest
```

### Caching

Mediator is does not offer a built-in cache and will re-download images on each request. The service is intended to run behind a CDN that will cache the processed images.

## Usage

## Configuration

Mediator can be configured using `ENV` variables:

| Variable                     | Description                                                                                                                                                                                                                         | Default                    |
| ---------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| `MEDIATOR_SOURCES`           | **Required**. List of supported sources to pull images from. A semicolon-separated list of `key=values`. Example:<br>`mybucket=https://mybucket.s3.amazonaws.com;mybucket-dev=https://mybucket-dev.s3.amazonaws.com`                | —                          |
| `MEDIATOR_SECRET_KEY`        | Optional, but **highly encouraged**. Secret random key, used to generate URL signatures. Can be ganerated with:<br> `dd if=/dev/urandom bs=32 count=1 2>/dev/null \| base64 \| tr -d '='`                                           | `""`                       |
| `MEDIATOR_AUTH_TOKEN`        | Optional. Token for authenticating image requests (for `Authorization: Bearer <AUTH_TOKEN>`). You can set it when configuring your CDN to prevent direct access to the service.                                                     | `""`                       |
| `MEDIATOR_PATH_PREFIX`       | Optional. Prefix for the image processing URL paths. Useful when you have a CDN pulling from multiple sources. <br>Example: `PATH_PREFIX=/my-prefix` will change the `transform` URL to `/my-prefix/image/transform/:source/:path`. | `""`                       |
| `MEDIATOR_CACHE_CONTROL`     | Value for the `Cache-Control` header.                                                                                                                                                                                               | `public, max-age=31536000` |
| `MEDIATOR_DOWNLOAD_MAX_SIZE` | File size download limit, in bytes.                                                                                                                                                                                                 | `50MB`                     |
| `MEDIATOR_DOWNLOAD_TIMEOUT`  | Download timeout, in seconds.                                                                                                                                                                                                       | `10s`                      |
| `MEDIATOR_HTTP_PORT`         | HTTP port for the service.                                                                                                                                                                                                          | `8000`                     |
| `MEDIATOR_LOG_LEVEL`         | Log level. Supported values: `debug`, `info`, `warn`, `error`.                                                                                                                                                                      | `info`                     |
