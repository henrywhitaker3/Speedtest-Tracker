# ramsey/collection

[![Source Code][badge-source]][source]
[![Latest Version][badge-release]][packagist]
[![Software License][badge-license]][license]
[![PHP Version][badge-php]][php]
[![Build Status][badge-build]][build]
[![Coverage Status][badge-coverage]][coverage]
[![Total Downloads][badge-downloads]][downloads]

ramsey/collection is a PHP 7.2+ library for representing and manipulating collections.

Much inspiration for this library came from the [Java Collections Framework][java].

This project adheres to a [code of conduct](CODE_OF_CONDUCT.md).
By participating in this project and its community, you are expected to
uphold this code.


## Installation

Install this package as a dependency using [Composer](https://getcomposer.org).

``` bash
composer require ramsey/collection
```

## Usage

The [latest class API documentation][apidocs] is available online.

Examples of how to use this framework can be found in the
[Wiki pages](https://github.com/ramsey/collection/wiki/Examples).

## Contributing

Contributions are welcome! Before contributing to this project, familiarize
yourself with [CONTRIBUTING.md](CONTRIBUTING.md).

To develop this project, you will need [PHP](https://www.php.net) 7.2 or greater
and [Composer](https://getcomposer.org).

After cloning this repository locally, execute the following commands:

``` bash
cd /path/to/repository
composer install
```

Now, you are ready to develop!

### Tooling

This project uses [CaptainHook](https://github.com/CaptainHookPhp/captainhook)
to validate all staged changes prior to commit.

#### Composer Commands

To see all the commands available in the project `br` namespace for
Composer, type:

``` bash
composer list br
```

##### Composer Command Autocompletion

If you'd like to have Composer command auto-completion, you may use
[bamarni/symfony-console-autocomplete](https://github.com/bamarni/symfony-console-autocomplete).
Install it globally with Composer:

``` bash
composer global require bamarni/symfony-console-autocomplete
```

Then, in your shell configuration file — usually `~/.bash_profile` or `~/.zshrc`,
but it could be different depending on your settings — ensure that your global
Composer `bin` directory is in your `PATH`, and evaluate the
`symfony-autocomplete` command. This will look like this:

``` bash
export PATH="$(composer config home)/vendor/bin:$PATH"
eval "$(symfony-autocomplete)"
```

Now, you can use the `tab` key to auto-complete Composer commands:

``` bash
composer br:[TAB][TAB]
```

#### Coding Standards

This project follows a superset of [PSR-12](https://www.php-fig.org/psr/psr-12/)
coding standards, enforced by [PHP_CodeSniffer](https://github.com/squizlabs/PHP_CodeSniffer).
The project PHP_CodeSniffer configuration may be found in `phpcs.xml.dist`.

CaptainHook will run PHP_CodeSniffer before committing. It will attempt to fix
any errors it can, and it will reject the commit if there are any un-fixable
issues. Many issues can be fixed automatically and will be done so pre-commit.

You may lint the entire codebase using PHP_CodeSniffer with the following
commands:

``` bash
# Lint
composer br:lint

# Lint and autofix
composer br:lint:fix
```

#### Static Analysis

This project uses a combination of [PHPStan](https://github.com/phpstan/phpstan)
and [Psalm](https://github.com/vimeo/psalm) to provide static analysis of PHP
code. Configurations for these are in `phpstan.neon.dist` and `psalm.xml`,
respectively.

CaptainHook will run PHPStan and Psalm before committing. The pre-commit hook
does not attempt to fix any static analysis errors. Instead, the commit will
fail, and you must fix the errors manually.

You may run static analysis manually across the whole codebase with the
following command:

``` bash
# Static analysis
composer br:analyze
```

### Project Structure

This project uses [pds/skeleton](https://github.com/php-pds/skeleton) as its
base folder structure and layout.

| Name              | Description                                    |
| ------------------| ---------------------------------------------- |
| **bin/**          | Commands and scripts for this project          |
| **build/**        | Cache, logs, reports, etc. for project builds  |
| **docs/**         | Project-specific documentation                 |
| **resources/**    | Additional resources for this project          |
| **src/**          | Project library and application source code    |
| **tests/**        | Tests for this project                         |

## Copyright and License

The ramsey/collection library is copyright © [Ben Ramsey](https://benramsey.com)
and licensed for use under the terms of the
MIT License (MIT). Please see [LICENSE](LICENSE) for more information.


[java]: http://docs.oracle.com/javase/8/docs/technotes/guides/collections/index.html
[apidocs]: https://docs.benramsey.com/ramsey-collection/latest/

[badge-source]: http://img.shields.io/badge/source-ramsey/collection-blue.svg?style=flat-square
[badge-release]: https://img.shields.io/packagist/v/ramsey/collection.svg?style=flat-square&label=release
[badge-license]: https://img.shields.io/packagist/l/ramsey/collection.svg?style=flat-square
[badge-php]: https://img.shields.io/packagist/php-v/ramsey/collection.svg?style=flat-square
[badge-build]: https://img.shields.io/travis/ramsey/collection/master.svg?style=flat-square
[badge-coverage]: https://img.shields.io/coveralls/github/ramsey/collection/master.svg?style=flat-square
[badge-downloads]: https://img.shields.io/packagist/dt/ramsey/collection.svg?style=flat-square&colorB=mediumvioletred

[source]: https://github.com/ramsey/collection
[packagist]: https://packagist.org/packages/ramsey/collection
[license]: https://github.com/ramsey/collection/blob/master/LICENSE
[php]: https://php.net
[build]: https://travis-ci.org/ramsey/collection
[coverage]: https://coveralls.io/r/ramsey/collection?branch=master
[downloads]: https://packagist.org/packages/ramsey/collection
