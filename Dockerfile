FROM lsiobase/alpine:3.12 as build-stage

# Install apk stuff
RUN apk add --no-cache --upgrade \
        gcc \
        cmake \
        curl-dev \
        libxml2-dev \
        build-base \
        openssl-dev \
        git && \
    cd /tmp && \
    git clone https://github.com/taganaka/SpeedTest && \
    cd SpeedTest && \
    cmake -DCMAKE_BUILD_TYPE=Release . && \
    cd /tmp/SpeedTest && \
    make install

FROM linuxserver/nginx
MAINTAINER henrywhitaker3@outlook.com

# Install apt stuff
# RUN apk add --no-cache --upgrade supervisor

# Copy over static files
COPY conf/ /

# Get SpeedTest++
COPY --from=build-stage /usr/local/bin/SpeedTest /site/app/Bin/

EXPOSE 80 443

VOLUME ["/config"]
