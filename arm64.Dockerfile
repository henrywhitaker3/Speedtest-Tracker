FROM linuxserver/nginx:arm64v8-latest
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='arm64'

COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
