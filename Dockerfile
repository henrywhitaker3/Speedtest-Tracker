FROM golang:1.27.3 AS gob

ARG VERSION="unknown"

WORKDIR /build

COPY . /build/

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-X main.version=${VERSION}" -a -o speedtest-tracker .

FROM scratch

COPY --from=gob /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=gob /build/speedtest-tracker /prog

ENTRYPOINT [ "/speedtest-tracker" ]
