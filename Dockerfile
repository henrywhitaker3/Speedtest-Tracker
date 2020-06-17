FROM linuxserver/nginx
MAINTAINER henrywhitaker3@outlook.com

# Install apt stuff
RUN apk add --no-cache --upgrade \
        gcc \
        cmake \
        curl-dev \
        libxml2-dev \
        build-base \
        openssl-dev \
        supervisor

# Copy over static files
COPY conf/ /setup/

# Get and compile SpeedTest++
RUN cd /tmp && \
    git clone https://github.com/taganaka/SpeedTest && \
    cd SpeedTest && \
    cmake -DCMAKE_BUILD_TYPE=Release . && \
    cd /tmp/SpeedTest && \
    make install && \
    mv /usr/local/bin/SpeedTest /setup/site/app/Bin/

# Setup new init script
RUN cp /setup/entrypoint/init.sh /etc/cont-init.d/50-speedtest

# Update webroot
RUN cp /setup/default /defaults/default

RUN mkdir -p /etc/services.d/supervisord/ && \
    cp /setup/supervisor-service.sh /etc/services.d/supervisord/run && \
    mkdir -p /etc/supervisor.d/ && \
    cp /setup/laravel-worker.conf /etc/supervisor.d/laravel-worker.ini

VOLUME ["/config"]
