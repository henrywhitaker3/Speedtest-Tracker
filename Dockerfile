FROM linuxserver/nginx
LABEL maintainer=henrywhitaker3@outlook.com

COPY conf/ /

RUN apk update && apk add ca-certificates wget

EXPOSE 80 443

VOLUME ["/config"]
