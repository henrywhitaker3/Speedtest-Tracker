# Speedtest Tracker

This program runs a speedtest check every hour and graphs the results. The back-end is written in [Laravel](https://laravel.com/) and the front-end uses [React](https://reactjs.org/). It uses the [speedtest-cli](https://github.com/sivel/speedtest-cli) package to get the data and uses [Chart.js](https://www.chartjs.org/) to plot the results.

![speedtest](https://user-images.githubusercontent.com/36062479/78822484-a82b8300-79ca-11ea-8525-fdeae496a0bd.gif)

## Installation

### Installing Dependencies

This program has some dependencies, to install them you need to run the following:

```bash
sudo apt update
sudo apt install software-properties-common
sudo add-apt [ONDREJ 7.4]
sudo apt install php7.4 composer python3 python3-pip
sudo pip install speedtest-cli
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

### Scheduling Setup

To get speed test results every hour, you need to add a cronjob, run `sudo crontab -e` and add an entry with the following:

```bash
* * * * * cd /path-to-your-project && php artisan schedule:run >> /dev/null 2>&1
```

### Queue Setup

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
