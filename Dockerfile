FROM linuxserver/nginx
LABEL maintainer=henrywhitaker3@outlook.com

COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
