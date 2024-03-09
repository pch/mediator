# Mediator

Mediator is a small, standalone service for efficient image processing. Images are processed on the fly with `libvips` and cached by CDN. The service is intended as a replacement for all common image operations in your web apps.

You can use it to:

- **Resize images.** Generate thumbnails and responsive images on the fly.
- **Crop images**. Crop images (currently only smart crop is supported).
- **Apply effects**. Apply filters to images (currently only `pixelate` is supported).
- **Proxy media files**: proxy images and other files using secure, signed URLs.
- **Strip metadata**. Remove metadata from images to reduce file size and protect user privacy.

### Examples

| Image                                                                                                                                                                                                  | URL                                                                                                                                                                                               | Description                       |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------- |
| ![](https://cdn.pixelpeeper.com/image/transform/images/images/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&s=05eaeaef8aac014ab3d0797c031bffadd5ec475a4166a81197e20cd9a84ca070&w=300)                   | https://cdn.pixelpeeper.com/image/transform/images/images/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&s=05eaeaef8aac014ab3d0797c031bffadd5ec475a4166a81197e20cd9a84ca070&w=300                   | Fit to 300x300, keep aspect ratio |
| ![](https://cdn.pixelpeeper.com/image/transform/images/images/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&op=fit%2Cpixelate&s=f461358f29c70c7b189c30de875597dddcd7a7f5ff5dd528f0069c501ce028cd&w=300) | https://cdn.pixelpeeper.com/image/transform/images/images/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&op=fit%2Cpixelate&s=f461358f29c70c7b189c30de875597dddcd7a7f5ff5dd528f0069c501ce028cd&w=300 | Fit to 300x300, pixelate          |
| ![](https://cdn.pixelpeeper.com/image/transform/images/images/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&op=smartcrop&s=80045b3f38f75772a5c78c65e5412db00b749c46cd2eb14d7411b246c2a3f2a7&w=300)      | https://cdn.pixelpeeper.com/image/transform/images/images/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&op=smartcrop&s=80045b3f38f75772a5c78c65e5412db00b749c46cd2eb14d7411b246c2a3f2a7&w=300      | Smart crop                        |

See also: [example code](https://github.com/pch/mediator/tree/main/examples)

## Features

- **Fast**. Written in Go and uses `libvips` for image processing.
- **Simple**. A small web service with minimal configuration.
- **Secure**. Supports signed URLs and authentication tokens.
- **Easy to integrate**. Works with any CDN and storage provider.
- **Easy to deploy**. Just use the provided Docker image.
- **Most common image operations**. Resize, crop, apply effects, and more.
- **Auto WebP**. Auto-convert images to WebP and serve them to browsers that support it.

## Considerations

- **Limited features**. Mediator aims to be cover the most common use cases. If you need more advanced features, consider a) contributing to the project, or b) using an alternative service.
- **New software**. Although Mediator had been running in production without issues, it's still a new project, shaped by limited use cases.

## Installation

Mediator is distributed as a docker image:

```shell
docker run --rm -it --env-file .env --publish 8000:8000 ghcr.io/pch/mediator:latest
```

### Caching

Mediator does not offer a built-in cache and will re-download images on each request. The service is intended to run behind a CDN, which will cache the processed images.

## Usage

### Transforming images

To transform an image, use the `/image/transform/:source/:path` endpoint. The `:source` parameter is the name of the source, and `:path` is the path to the image file (e.g. S3 key/path).

#### Request params

The `/image/transform` endpoint accepts the following query parameters:

| Parameter        | Description                                                                                                                                                                      |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `op`             | Operation names, separated by commas. Supported operations: `fit`, `smartcrop`, `pixelate`. Default: `fit`                                                                       |
| `w`              | Width of the target image.                                                                                                                                                       |
| `h`              | Height of the target image.                                                                                                                                                      |
| `format`         | Output format. Supported values: `jpeg`, `png`, `gif`, `webp`, `auto`. Defaults to `Content-Type` of the requested image. Set `auto` to return WebP to browsers that support it. |
| `strip`          | Strip metadata from the image. Supported values: `true`, `false`. Default: `true`                                                                                                |
| `q`              | Quality of the output image. Supported values: `0-100`. Default: `80`                                                                                                            |
| `pixelatefactor` | Pixelate factor, for example: `1-100`. The smaller the number, the less "pixelized" the result will be. Default: `20`                                                            |
| `s`              | Signature. Required when `MEDIATOR_SECRET_KEY` is set.                                                                                                                           |

#### Operations

Currently, the following operations are supported:

- **Fit**. Resize the image to fit within the specified dimensions, keeping the aspect ratio. The image will be downsized to the largest size that fits within the specified dimensions.
- **Smartcrop**. Crop the image to the specified dimensions, using a smart algorithm to find the most interesting part of the image.
- **Pixelate**. Pixelate the image. The `pixelatefactor` parameter controls the level of pixelation.

### File proxy

To proxy a file from a private source, use the `/proxy/:source/:path` endpoint. The `:source` parameter is the name of the source URL. The `:path` and the query string are passed to the source URL.

**Use case:** this was originally built to proxy screenshots and PDFs generated by a private service, running puppeteer. Mediator forwards the requests to the service, taking care of the URL signing (to prevent spamming the service with automated requests) and caching (via CDN).

**NOTE:** `Cache-Control` is not set for proxied files. The source server should take care of the caching policy.

## URL Signing

The signing mechanism is based on the `MEDIATOR_SECRET_KEY` environment variable. When the key is set, the service will require a `s` parameter in the query string, containing the signature of the request.

The algorithm to generate the signature:

```
path_with_query_params = "/image/transform/images-dev/2024/02/2431ckhy8z9s1tya4trk21pb3f.jpg?h=300&w=300"
signature = HMAC-SHA256(MEDIATOR_SECRET_KEY, path_with_query_params)
signed_url = path_with_query_params + "&s=" + signature
```

See also: [examples](https://github.com/pch/mediator/tree/main/examples)

---

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

## Deployment

### Deploying with Kamal

Example configuration for [Kamal](https://kamal-deploy.org/), with Mediator running as an accessory:

```yaml
service: example
image: pch/example # your main app

servers:
  web:
    hosts:
      - 123.123.123.123
    labels:
      traefik.http.routers.peeper.entrypoints: websecure
      traefik.http.routers.peeper.rule: Host(`example.com`)
      traefik.http.routers.peeper.tls.certresolver: letsencrypt

accessories:
  mediator:
    image: ghcr.io/pch/mediator:latest
    roles:
      - web # will run on the same host as the main app
    port: "8000:8000"
    env:
      clear:
        MEDIATOR_SOURCES: images-dev=https://example-dev.s3.amazonaws.com;images=https://example-prod.s3.amazonaws.com
      secret:
        - MEDIATOR_SECRET_KEY
        - MEDIATOR_AUTH_TOKEN
    labels:
      traefik.http.routers.mediator.entrypoints: websecure
      traefik.http.routers.mediator.rule: Host(`mediator.example.com`)
      traefik.http.routers.mediator.tls.certresolver: letsencrypt
      traefik.tcp.services.mediator.loadbalancer.server.port: 8000
```

### CloudFront setup

If you decide to use CloudFront for CDN, there are only a few considerations to take into account:

- When setting up an origin, make sure to set the `Authorization` header to `Bearer <YOUR_AUTH_TOKEN>` to prevent direct (non-cached) access to the service.
- In the Behavior settings, you have to make sure that the query string is forwarded
- In the "Choose which headers to include in the cache key" part, add the `Accept` header if you want to serve WebP images to browsers that support it.
