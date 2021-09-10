FROM linuxserver/nginx:arm32v7-latest
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='arm'

COPY docker/conf/ /
COPY . /site

EXPOSE 80 443

VOLUME ["/config"]
