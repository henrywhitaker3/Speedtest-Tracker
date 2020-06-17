FROM linuxserver/nginx
MAINTAINER henrywhitaker3@outlook.com

# Install apt stuff
RUN apk add --no-cache --upgrade \
        python3 \
        py-pip \
        supervisor

# Install speedtest-cli
RUN pip3 install speedtest-cli

# Copy over static files
COPY conf/ /setup/

# Setup new init script
RUN cp /setup/entrypoint/init.sh /etc/cont-init.d/50-speedtest

# Update webroot
RUN cp /setup/default /defaults/default

RUN mkdir -p /etc/services.d/supervisord/ && \
    cp /setup/supervisor-service.sh /etc/services.d/supervisord/run && \
    mkdir -p /etc/supervisor.d/ && \
    cp /setup/laravel-worker.conf /etc/supervisor.d/laravel-worker.ini

VOLUME ["/config"]
