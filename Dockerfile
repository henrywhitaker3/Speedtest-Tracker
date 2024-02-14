FROM linuxserver/nginx:1.20.1-r3-ls142
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='x86_64'

COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
