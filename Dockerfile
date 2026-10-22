FROM --platform=$BUILDPLATFORM ghcr.io/henrywhitaker3/containers/mise:v1.4.1 AS uib

ENV NUXT_PUBLIC_API_URL=""

WORKDIR /build/ui

COPY frontend/mise.toml frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./

RUN apk add --no-cache bash libgcc libstdc++

ENV MISE_ALL_COMPILE=false

RUN mise install

RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    mise exec -- pnpm install --frozen-lockfile --store-dir=/pnpm/store

COPY frontend/ ./
RUN mise exec -- pnpm build

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
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -ldflags="-X main.version=${VERSION}" -o speedtest-tracker .

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=gob /build/speedtest-tracker /speedtest-tracker
COPY --from=uib /build/ui/.output/public /ui
USER 65532:65532

ENTRYPOINT [ "/speedtest-tracker" ]
