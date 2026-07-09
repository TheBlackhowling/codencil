# Dev image — build/run via Docker Compose only (see agents/BUILD_ORDER.md P0.2).
FROM golang:1.23-bookworm AS dev

WORKDIR /app

CMD ["go", "version"]

FROM golang:1.23-bookworm AS build

WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api/ ./
RUN CGO_ENABLED=0 go build -o /codencil ./cmd/codencil

FROM gcr.io/distroless/static-debian12:nonroot AS prod

COPY --from=build /codencil /codencil
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/codencil"]
