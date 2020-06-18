FROM linuxserver/nginx
MAINTAINER henrywhitaker3@outlook.com

# Copy over static files
COPY conf/ /

EXPOSE 80 443

VOLUME ["/config"]
