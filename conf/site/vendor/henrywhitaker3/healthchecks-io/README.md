# PHP healthchecks.io

![version](https://img.shields.io/github/v/tag/henrywhitaker3/PHP-healthchecks.io?color=light-green&label=version&style=flat-square) [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/henrywhitaker3/PHP-healthchecks.io/PHPUnit?label=build&logo=github&style=flat-square)](https://github.com/henrywhitaker3/PHP-healthchecks.io/actions) [![last_commit](https://img.shields.io/github/last-commit/henrywhitaker3/PHP-healthchecks.io?style=flat-square)](https://github.com/henrywhitaker3/PHP-healthchecks.io/commits) [![issues](https://img.shields.io/github/issues/henrywhitaker3/PHP-healthchecks.io?style=flat-square)](https://github.com/henrywhitaker3/PHP-healthchecks.io/issues) [![commit_freq](https://img.shields.io/github/commit-activity/m/henrywhitaker3/PHP-healthchecks.io?style=flat-square)](https://github.com/henrywhitaker3/PHP-healthchecks.io/commits) [![license](https://img.shields.io/github/license/henrywhitaker3/PHP-healthchecks.io?style=flat-square)](https://github.com/henrywhitaker3/PHP-healthchecks.io/blob/master/LICENSE)

This package provides simple methods to interact with the [healthchecks.io](https://healthchecks.io) API.

## Installation

```bash
composer require henrywhitaker3/healthchecks-io
```

## Usage

This package provides two classes you can use:

### Healthchecks

This class is just used to update a healthchecks.io endpoint.

To create a new instance, you need to pass the UUID for the endpoint you want to update:

```php
$hc = new Healthchecks('SAMPLE-UUID-HERE');
```

Once you have your created a new instance, there are 3 methods to use:

```php
$hc->success(); // sends a 'success' signal
$hc->fail(); // sends a 'fail' signal
$hc->start(); // sends a 'start' signal
```

### HealthchecksManager

This class can be used to interact with the healthchecks.io management API. Creating a new instance is simple:

```php
$hm = new HealthchecksManager('SAMPLE-API-KEY');
```

You now have several methods you can use:

```php
$hm->listChecks(); // lists all checks
$hm->getCheck('UUID'); // get info for a specific check
$hm->pauseCheck('UUID'); // pauses a check
$hm->resumeCheck('UUID'); // pings a check to resume it
$hm->deleteCheck('UUID'); // deletes the check
$hm->getCheckPings('UUID'); // gets a list of pings for the check
$hm->getCheckStatusChanges('UUID'); // returns a list of "flips" this check has experienced
$hm->createCheck($args); // create a new check
$hm->updateCheck('UUID', $args); // update an existing check
$hm->integrations(); // get a list of integrations (i.e. slack, discord etc.)
```

## Contributing

Contributions are welcome, but please write tests for whatever you add. Create a new file in the `tests/tests` directory with the following structure:

```php
<?php

namespace Henrywhitaker3\Healthchecks\Tests;

class NewTest extends Test
{
    //
}
```

As these tests will need valid credentials and a UUID, you will need to setup a `.env` file to run tests successfully. Simply copy the `.env.example` to `.env` and fill in the relevant details.