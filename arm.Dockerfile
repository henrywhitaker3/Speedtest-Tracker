FROM linuxserver/nginx:arm32v7-latest
LABEL maintainer=henrywhitaker3@outlook.com

COPY conf/ /

RUN sed -i 's/1\.0\.0-x86_64/1\.0\.0-arm/' /etc/cont-init.d/50-speedtest

EXPOSE 80 443

VOLUME ["/config"]
