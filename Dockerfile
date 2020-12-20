FROM linuxserver/nginx
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='x86_64'

COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
