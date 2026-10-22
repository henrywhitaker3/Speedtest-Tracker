FROM --platform=$BUILDPLATFORM golang:1.27.3 AS gob

ARG VERSION="unknown"

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY internal/ internal/
COPY main.go main.go
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-X main.version=${VERSION}" -o speedtest-tracker .

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=gob /build/speedtest-tracker /speedtest-tracker
USER 65532:65532

ENTRYPOINT [ "/speedtest-tracker" ]
