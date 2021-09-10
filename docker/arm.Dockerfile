FROM linuxserver/nginx:arm32v7-latest
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='arm'

COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
