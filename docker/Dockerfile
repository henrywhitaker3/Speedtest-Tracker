FROM linuxserver/nginx
LABEL maintainer=henrywhitaker3@outlook.com

ENV arch='x86_64'

COPY docker/conf/ /
COPY . /site

RUN cd /site composer install

EXPOSE 80 443

VOLUME ["/config"]
