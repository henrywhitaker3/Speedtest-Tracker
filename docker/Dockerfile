FROM linuxserver/nginx
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='x86_64'

RUN apk add php82-tokenizer

COPY docker/conf/ /
COPY . /site

EXPOSE 80 443

VOLUME ["/config"]
