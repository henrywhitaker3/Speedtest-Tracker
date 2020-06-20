# Speedtest Tracker

[![Docker pulls](https://img.shields.io/docker/pulls/henrywhitaker3/speedtest-tracker)](https://hub.docker.com/r/henrywhitaker3/speedtest-tracker) [![last_commit](https://img.shields.io/github/last-commit/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits) [![issues](https://img.shields.io/github/issues/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/issues) ![version](https://img.shields.io/badge/version-v1.5.6-success) [![license](https://img.shields.io/github/license/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/blob/master/LICENSE)

This program runs a speedtest check every hour and graphs the results. The back-end is written in [Laravel](https://laravel.com/) and the front-end uses [React](https://reactjs.org/). It uses the [speedtest-cli](https://www.speedtest.net/apps/cli) package to get the data and uses [Chart.js](https://www.chartjs.org/) to plot the results.

![speedtest](https://user-images.githubusercontent.com/36062479/78822484-a82b8300-79ca-11ea-8525-fdeae496a0bd.gif)

## Features

- Automatically run a speedtest every hour
- Graph of previous speedtests going back x days
- Backup/restore data in JSON/CSV format
- Slack/Discord notifications
- Organizr integration

## Installation & Setup

### Using Docker

A docker image is available [here](https://hub.docker.com/r/henrywhitaker3/speedtest-tracker), you can create a new conatiner by running:

```bash
docker create \
      --name=speedtest \
      -p 8765:80 \
      -v /path/to/data:/config \
      -e SLACK_WEBHOOK=webhook `#optional` \
      -e PUID=uid `#optional` \
      -e PGID=gid `#optional` \
      --restart unless-stopped \
      henrywhitaker3/speedtest-tracker
```

### Manual Install

#### Installing Dependencies

This program has some dependencies, to install them you need to run the following:

```bash
sudo apt update
sudo apt update
sudo apt install php-common php7.2 php7.2-cli php7.2-common php7.2-json php7.2-opcache php7.2-readline php-xml php-sqlite3 php-zip composer python3 python3-pip git
```
```bash
sudo apt install curl
curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
sudo apt install nodejs
```

```bash
sudo pip3 install speedtest-cli
```

Then, download the code by running:

```bash
git clone https://github.com/henrywhitaker3/Speedtest-Tracker.git
```

Install the composer and npm dependencies:

```bash
composer install
npm install && npm run production
```

#### Setting up the database

Run the following to set your database variables:

```bash
cp .env.example .env
```

Then update the `DB_DATABASE` value with the absolute path of your install, followed by `/database/speed.db`.

Finally, run the following to setup the tables in the database:

```bash
php artisan key:generate
php artisan migrate
```

Now run the following to make sure everything has been setup properly (it should output a speedtest result):

```bash
php artisan speedtest:run
```

#### Scheduling Setup

To get speed test results every hour, you need to add a cronjob, run `sudo crontab -e` and add an entry with the following (with the path you your install):

```bash
* * * * * cd /path-to-your-project && php artisan schedule:run >> /dev/null 2>&1
```

#### Queue Setup

```bash
sudo apt install supervisor
```

```bash
sudo vim /etc/supervisor/conf.d/laravel-worker.conf
```

Add the following, updating the `command` and user `values`:

```bash
[program:laravel-worker]
process_name=%(program_name)s_%(process_num)02d
command=php /path/to/project/artisan queue:work
autostart=true
autorestart=true
user=<user>
numprocs=8
redirect_stderr=true
```

Then run:

```bash
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl restart all
```
