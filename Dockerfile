FROM node:lts AS builder

WORKDIR /build

# can't cache because everything is being copied over
# TODO: move web ui things into seperate folder
COPY . /build

RUN \
npm ci && \
npm run production && \
rm -rf node_modules

FROM linuxserver/nginx:latest
LABEL maintainer=henrywhitaker3@outlook.com

COPY docker/ /
COPY --from=builder /build/ /site/

EXPOSE 80 443

VOLUME [ "/config" ]
