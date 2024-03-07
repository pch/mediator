FROM golang:1.22 AS builder

WORKDIR /app

RUN apt-get update && \
  apt-get install -y libvips-dev \
  && apt-get autoremove -y \
  && apt-get clean -y \
  && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
RUN go mod download

COPY . /app
RUN go build -ldflags "-X main.Version=$(cat VERSION | tr -d '\n')" -o ./bin/mediator ./cmd/mediator

# Production stage
FROM ubuntu:latest AS production

RUN apt-get update && \
  apt-get install -y libvips-dev \
  && apt-get autoremove -y \
  && apt-get clean -y \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/bin/mediator .

RUN useradd appuser --create-home --shell /bin/bash && \
  chown -R appuser:appuser .
USER appuser:appuser

EXPOSE 8000

CMD ["./mediator"]

# Development stage
FROM builder AS development

ENV GOCACHE=/go/.cache

RUN useradd appuser --create-home --shell /bin/bash && \
  chown -R appuser:appuser .
USER appuser:appuser

EXPOSE 8000

CMD ["go", "run", "cmd/mediator/main.go"]
