FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
  build-base \
  vips-dev

COPY go.mod ./
RUN go mod download

COPY . /app
RUN go build -ldflags "-X main.Version=$(cat VERSION | tr -d '\n')" -o ./bin/mediator ./cmd/mediator

# Production stage
FROM alpine:latest AS production

RUN apk add --no-cache \
  vips-dev

WORKDIR /app
COPY --from=builder /app/bin/mediator .

RUN adduser -D appuser && \
  chown -R appuser:appuser .
USER appuser:appuser

EXPOSE 8000

CMD ["./mediator"]

# Development stage
FROM builder AS development

ENV GOCACHE=/go/.cache

RUN adduser -D appuser && \
  chown -R appuser:appuser .
USER appuser:appuser

EXPOSE 8000

CMD ["go", "run", "cmd/mediator/main.go"]
