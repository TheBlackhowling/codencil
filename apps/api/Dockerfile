# Dev image — build/run via Docker Compose only (see agents/BUILD_ORDER.md P0.2).
FROM golang:1.23-bookworm

WORKDIR /app

# Default command overridden in compose for dev; API entrypoint added in P0.4.
CMD ["go", "version"]
