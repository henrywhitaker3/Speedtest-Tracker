# Speedtest Tracker

This program runs a speedtest check every hour and graphs the results. The back-end is written in [Laravel](https://laravel.com/) and the front-end uses [React](https://reactjs.org/). It uses the [speedtest-cli](https://github.com/sivel/speedtest-cli) package to get the data and uses [Chart.js](https://www.chartjs.org/) to plot the results.

![speedtest](https://user-images.githubusercontent.com/36062479/78811223-0a7b8800-79b9-11ea-8259-4fb7890c8a14.gif)

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
git clone [URL]
```

Install the composer and npm dependencies:

```bash
composer install
npm install && npm run production
```

### Queue Setup

`sudo apt install supervisor`

`sudo vim /etc/supervisor/conf.d/laravel-worker.conf`

Add the following:

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
