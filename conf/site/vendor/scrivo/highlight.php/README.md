# highlight.php

[![Build Status](https://travis-ci.org/scrivo/highlight.php.svg?branch=master)](https://travis-ci.org/scrivo/highlight.php)
[![Latest Packagist release](https://img.shields.io/packagist/v/scrivo/highlight.php.svg)](https://packagist.org/packages/scrivo/highlight.php)
[![Monthly downloads on Packagist](https://img.shields.io/packagist/dm/scrivo/highlight.php.svg)](https://packagist.org/packages/scrivo/highlight.php/stats)

*highlight.php* is a server-side syntax highlighter written in PHP that currently supports 185 languages. It's a port of [highlight.js](http://www.highlightjs.org) by Ivan Sagalaev that makes full use of the language and style definitions of the original JavaScript project.

[TOC]: # "## Table of Contents"

## Table of Contents
- [Installation + Setup](#installation--setup)
    - [Composer Version Constraints](#composer-version-constraints)
- [Usage](#usage)
    - [Explicit Mode](#explicit-mode)
    - [Automatic Language Detection Mode](#automatic-language-detection-mode)
    - [Stylesheets](#stylesheets)
- [Highlighter Utilities](#highlighter-utilities)
- [Versioning](#versioning)
    - [Backward Compatibility Promise](#backward-compatibility-promise)
- [Some History](#some-history)
- [License](#license)


## Installation + Setup

The recommended approach is to install the project through [Composer](https://getcomposer.org/).

```bash
composer require scrivo/highlight.php
```

If you're not using Composer, ensure that the classes defined in the `Highlight` namespace can be found either by inclusion or by an autoloader. A trivial autoloader for this purpose is included in this project as `Highlight\Autoloader.php`

### Composer Version Constraints

When requiring this project in your `composer.json`, it is recommended you use the [caret version range](https://getcomposer.org/doc/articles/versions.md#caret-version-range-) and use only the major and minor values; i.e. `^9.14`.

It's come to our attention that a lot of tutorials and projects out there are locking themselves into highly specific versions of this project; e.g. `"scrivo/highlight.php": "v9.12.0.1"`. Please do **not** do this or encourage it. We promise a [reliable backward compatibility policy](#backward-compatibility-promise) so there's no reason to lock yourself to such a specific version. By doing this, you are preventing yourself or your users from receiving updates to language definitions and bug fixes.

## Usage

The `\Highlight\Highlighter` class contains the syntax highlighting functionality. You can choose between two highlighting modes:

1. explicit mode
2. automatic language detection mode

### Explicit Mode

In explicit mode, you must define which language you will be highlighting as.

```php
// Instantiate the Highlighter.
$hl = new \Highlight\Highlighter();
$code = file_get_contents('some_ruby_script.rb');

try {
    // Highlight some code.
    $highlighted = $hl->highlight('ruby', $code);

    echo "<pre><code class=\"hljs {$highlighted->language}\">";
    echo $highlighted->value;
    echo "</code></pre>";
}
catch (DomainException $e) {
    // This is thrown if the specified language does not exist

    echo "<pre><code>";
    echo $code;
    echo "</code></pre>";
}
```

### Automatic Language Detection Mode

Alternatively you can use the automatic detection mode, which highlights your code with the language the library thinks is best. It is highly recommended you explicitly choose the language or limit the number of languages to automatically detect to reduce the number of inaccuracies.

> **Warning:** Auto-detection occurs in a brute force fashion and the language with the most accurate result will be selected. This is extremely inefficient as you supply more languages and may not always be 100% accurate if similar languages are configured.

```php
$hl = new \Highlight\Highlighter();
$hl->setAutodetectLanguages(array('ruby', 'python', 'perl'));

$highlighted = $hl->highlightAuto(file_get_contents('some_ruby_script.rb'));

echo "<pre><code class=\"hljs {$highlighted->language}\">";
echo $highlighted->value;
echo "</code></pre>";
```

#### Default Languages

In version 9.x of this project, the following languages are the default auto-detected languages:

- XML
- JSON
- JavaScript
- CSS
- PHP
- HTTP

These default languages are considered "legacy behavior" and will be removed in version 10.x of this library to match highlight.js behavior; the new default behavior in 10.x will be to use *every* language.

### Stylesheets

The same stylesheets available in the **highlight.js** project are available in the `styles` directory of this project and may be included in your own CSS or made accessible to your web server.

## Highlighter Utilities

The core of the project is loyal port of **highlight.js** and is available under the main `Highlight` namespace. A series of convenience functions are provided under the `HighlightUtilities` namespace to introduce additional functionality without the need for another dependency.

Available functions:

- [`getAvailableStyleSheets(bool $filePaths = false): string[]`](HighlightUtilities/functions.php#L35-L48)
- [`getStyleSheet(string $name): false|string`](HighlightUtilities/functions.php#L94-L107)
- [`getStyleSheetFolder(): string`](HighlightUtilities/functions.php#L115-L123)
- [`getStyleSheetPath(string $name): string`](HighlightUtilities/functions.php#L131-L143)
- [`getThemeBackgroundColor(string $name): float[]`](HighlightUtilities/functions.php#L75-L88)
- [`splitCodeIntoArray(string $html): false|string[]`](HighlightUtilities/functions.php#L156-L169)

## Versioning

This project will follow the same version numbers as the highlight.js project with regards to languages, meaning that a language definition available in highlight.js 9.12.0 will be available in highlight.php 9.12.0. However, there are times where bugs may arise in this project or its translated definition files, so there'll be one more number appended to the version number. For example, version 9.12.0.1 will contain all of the same languages as highlight.js 9.12.0 but also contain fixes solely to this project. This is done so this project can have version bumps without conflicts should highlight.js release version 9.12.1.

### Backward Compatibility Promise

Despite the fact that the semantic versioning used in this project mirrors that of highlight.js, this project will adhere to [Symfony's Backward Compatibility Promise](https://symfony.com/doc/current/contributing/code/bc.html#using-symfony-code). You can rest assured that there will be no breaking changes during `9.x` and any deprecations will be marked with `@deprecated` and won't be removed until the next major release.

## Some History

Geert Bergman
Sep 30, 2013

JavaScript code highlighting is very convenient and in many cases just what you want to use. Especially for programming blogs I would not advice you to use otherwise. But there are occasions where you're better off with a more 'static' approach, for instance if you want to send highlighted code in an email or for API documents. For this I needed a code highlighting program preferably written in PHP.

I couldn't found any satisfactory PHP solution so I decided to port one from JavaScript. After some comparison of different highlighting programs based on license, technology, language support [highlight.js](http://www.highlightjs.org) came out most favorable in my opinion.

It was my decision not to make a PHP highlighter but to do a port of highlight.js, these are different things. The goal was to make it work exactly as [highlight.js](http://www.highlightjs.org) to make as much use as possible of the language definitions and CSS files of the original program.

Happy coding!

## License

[BSD](./LICENSE.md)
