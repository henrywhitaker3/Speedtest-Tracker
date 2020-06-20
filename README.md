# Speedtest Tracker

[![Docker pulls](https://img.shields.io/docker/pulls/henrywhitaker3/speedtest-tracker)](https://hub.docker.com/r/henrywhitaker3/speedtest-tracker) [![last_commit](https://img.shields.io/github/last-commit/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits) [![issues](https://img.shields.io/github/issues/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/issues) [![commit_freq](https://img.shields.io/github/commit-activity/m/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/commits) ![version](https://img.shields.io/badge/version-v1.5.5-success) [![license](https://img.shields.io/github/license/henrywhitaker3/Speedtest-Tracker)](https://github.com/henrywhitaker3/Speedtest-Tracker/blob/master/LICENSE)

This program runs a speedtest check every hour and graphs the results. The back-end is written in [Laravel](https://laravel.com/) and the front-end uses [React](https://reactjs.org/). It uses [Ookla's Speedtest cli](https://www.speedtest.net/apps/cli) to get the data and uses [Chart.js](https://www.chartjs.org/) to plot the results.

Disclaimer: You will need to accept Ookla's [EULA](https://www.speedtest.net/about/eula) and privacy agreements in order to use this container.

![speedtest](https://user-images.githubusercontent.com/36062479/78822484-a82b8300-79ca-11ea-8525-fdeae496a0bd.gif)

## Features

- Automatically run a speedtest every hour
- Graph of previous speedtests going back x days
- Backup/restore data in JSON/CSV format
- Slack/Discord notifications
- Organizr integration

## Usage

```bash
docker create \
      --name=speedtest \
      -p 8765:80 \
      -v /path/to/data:/config \
      -e SLACK_WEBHOOK=webhook `#optional` \
      -e PUID=uid `#optional` \
      -e PGID=gid `#optional` \
      -e OOKLA_EULA_GDPR=true \
      --restart unless-stopped \
      henrywhitaker3/speedtest-tracker
```

## Parameters

Container images are configured using parameters passed at runtime (such as those above). These parameters are separated by a colon and indicate `<external>:<internal>` respectively. For example, `-p 8080:80` would expose port `80` from inside the container to be accessible from the host's IP on port `8080` outside the container.

|     Parameter             |   Function    |
|     :----:                |   --- |
|     `-p 8765:80`          |   Exposes the webserver on port 8765  |
|     `-v /config`          |   All the config files reside here.   |
|     `-e OOKLA_EULA_GDPR`  |   Set to 'true' to accept the Ookla [EULA](https://www.speedtest.net/about/eula) and privacy agreement. If this is not set, the container will not start   |
|     `-e SLACK_WEBHOOK`    |   Optional. Put a slack webhook here to get slack notifications when a speedtest is run. To use discord webhooks, just append `/slack` to the end of your discord webhook URL   |
|     `-e PUID`             |   Optional. Supply a local user ID for volume permissions   |
|     `-e PGID`             |   Optional. Supply a local group ID for volume permissions  |


## Getting the Image

To get the base image, you have 2 options:

- Use the pre-built image on dockerhub
- Build the image yourself

### Pre-built Image

Run `docker pull henrywhitaker3/speedtest-tracker`

### Dockerfile

Clone the required files from the github repo [here](https://github.com/henrywhitaker3/Speedtest-Tracker/tree/docker) making sure to use the `docker` branch of the repo.

Build the image from the docker file by running (within the cloned git repo):

```bash
docker build . -f Dockerfile --tag=henrywhitaker3/speedtest-tracker:<tag>
```
