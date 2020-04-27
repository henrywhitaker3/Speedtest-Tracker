#!/bin/bash
# This script sets up the

if [ ! -f /config/speed.db ]; then
    echo "Database file not found!"
    echo "Creating the database..."
    touch /config/speed.db
fi

if cat /app/site/.env | grep -E "APP_KEY=base64" > /dev/null; then
    echo "App key found. Skipping app key generation"
else
    echo "No app key found. Generating app key"
    cd /app/site && php artisan key:generate
fi

if [ -z ${SLACK_WEBHOOK+x} ]; then
    echo "Slack webhook is unset"
else
    sed "s,SLACK_WEBHOOK=.*,SLACK_WEBHOOK=$SLACK_WEBHOOK," -i.bak .env
fi

cd /app/site && php artisan migrate

cd /config
mkdir -p logs

chown -R application /config
chmod 775 -R /config

echo "* * * * * cd /app/site/ && php artisan schedule:run >> /dev/null 2>&1" | crontab -
