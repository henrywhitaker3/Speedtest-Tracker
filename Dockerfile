FROM webdevops/php-nginx:ubuntu-18.04
MAINTAINER henrywhitaker3@outlook.com

ENV WEB_DOCUMENT_ROOT  /app/site/public
ARG DEBIAN_FRONTEND=noninteractive

# Install apt stuff
RUN apt-get update
RUN apt-get install \
            python3 \
            python3-pip -y

# Install speedtest-cli
RUN pip3 install speedtest-cli

# Copy over static files
RUN mkdir /setup
COPY conf/ /setup/

# Copy over Speedtest site files
RUN cp -r /setup/site/ /app/

# Setup env file
RUN cd /app/site \
    && cp .env.example .env \
    && sed 's/DB_DATABASE=.*/DB_DATABASE=\/config\/speed.db/' -i.bak .env

# Copy supervisor queue worker config
RUN cp /setup/laravel-worker.conf /opt/docker/etc/supervisor.d/

# Copy entry script
RUN cp /setup/entrypoint/init.sh /entrypoint.d \
    && chmod +x /setup/entrypoint/init.sh

# Set permissions for files
RUN chown -R application:application /app/site/storage \
    && chmod -R 775 /app/site/storage \
    && chown -R application:application /app/site/bootstrap \
    && chmod -R 775 /app/site/bootstrap

# Cleanup setup files
RUN rm -r /setup

VOLUME ["/config"]
