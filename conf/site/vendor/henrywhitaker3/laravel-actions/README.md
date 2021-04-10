# Laravel Actions

[![Latest Version on Packagist][icon-version]][link-packagist]
![Tests status][icon-tests]
[![StyleCI][icon-styleci]][link-styleci]

A simple Laravel actions package.

## Installation

```bash
composer require henrywhitaker3/laravelactions
```

## Usage

To create a new action, run:

```bash
php artisan make:action <name>
```

This will create a new action in `App/Actions`. You can then use the action in one of two ways:

```php
$action = new SomeAction();
$action->run();
```

Or you can use the `run` helper function:

```php
run(SomeAction::class);
```

You can pass arguments for the action's run method as the subsequent arguments for the helper method:

```php
run(SomeAction::class, $arg);
run(SomeAction::class, $arg1, $arg2);
```

## License

MIT license. Please see the [license file](LICENSE.md) for more information.

[icon-tests]: https://img.shields.io/github/workflow/status/henrywhitaker3/laravel-actions/PHP%20Composer/master?label=Tests&logoColor=%234c1&style=flat-square
[icon-styleci]: https://github.styleci.io/repos/335909164/shield?branch=master
[icon-version]: https://img.shields.io/packagist/v/henrywhitaker3/laravel-actions?style=flat-square&label=Version

[link-styleci]: https://github.styleci.io/repos/335909164
[link-packagist]: https://packagist.org/packages/henrywhitaker3/laravel-actions