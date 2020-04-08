# Speedtest Checker

## Queue Setup

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
