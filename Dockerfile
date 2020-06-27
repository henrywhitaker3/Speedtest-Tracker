FROM linuxserver/nginx
MAINTAINER henrywhitaker3@outlook.com

COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
