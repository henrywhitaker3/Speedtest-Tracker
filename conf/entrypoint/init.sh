#!/usr/bin/with-contenv bash
# # This script sets up the speedtest app

# Copy site files to /config
echo "Copying latest site files to config"
cp -rfT /setup/site/ /config/www/
chown -R abc:abc /config/www
chmod -R 755 /config/www/storage
chmod -R 755 /config/www/bootstrap

# Check for DB
if [ ! -f /config/speed.db ]; then
    echo "Database file not found! Creating empty database"
    touch /config/speed.db
    chown abc:abc /config/speed.db
else
    echo "Database file exists"
    chown abc:abc /config/speed.db
fi


# Check for .env
if [ ! -f /config/www/.env ]; then
    echo "Env file not found! Creating .env file"
    cp /setup/site/.env.example /config/www/.env
else
    echo "Env file exists"
fi

sed "s,DB_DATABASE=.*,DB_DATABASE=/config/speed.db," -i.bak /config/www/.env

echo "Running database migrations"
php /config/www/artisan migrate

# Check app key exists
if cat /config/www/.env | grep -E "APP_KEY=[0-9A-Za-z:+\/=]{1,}" > /dev/null ; then
    echo "App key exists"
else
    echo "Generating app key"
    php /config/www/artisan key:generate
fi

# Check JWT secret exists
if cat /config/www/.env | grep -E "JWT_SECRET=[0-9A-Za-z:+\/=]{1,}" > /dev/null ; then
    echo "JWT secret exists"
else
    echo "Generating JWT secret"
    php /config/www/artisan jwt:secret
fi

if [ -z ${SLACK_WEBHOOK+x} ]; then
    echo "Slack webhook is unset"
    sed "s,SLACK_WEBHOOK=.*,SLACK_WEBHOOK=," -i.bak /config/www/.env
else
    echo "Slack webhook set, updating .env"
    sed "s,SLACK_WEBHOOK=.*,SLACK_WEBHOOK=$SLACK_WEBHOOK," -i.bak /config/www/.env
fi

if [ -z ${BASE_PATH+x} ]; then
    echo "Base path is unset"
    sed "s,BASE_PATH=.*,BASE_PATH=," -i.bak /config/www/.env
else
    echo "Base path set, updating .env"
    sed "s,BASE_PATH=.*,BASE_PATH=$BASE_PATH," -i.bak /config/www/.env
fi

mkdir -p /config/log/speedtest
echo "* * * * * php /config/www/artisan schedule:run >> /config/log/speedtest/cron.log" >> /etc/crontabs/root
