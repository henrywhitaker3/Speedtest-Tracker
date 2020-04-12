# Contributing

Welcome to your server-side syntax highlighting solution! If you're interested in contributing to the project, here are a few things to keep in mind!

## Language Definitions + Styles

We do not accept PRs with new or updated language definitions or stylesheets. These should be contributed to the `highlight.js` project instead. The process of maintaining styles and transforming language definitions is completely automated and left up to project maintainers to run whenever a new version of `highlight.js` is released.

## Updates from `highlight.js`

If you'd like to make a PR containing behavior changes made in the `highlight.js` project, please make sure you link to the appropriate `highlight.js` commits and change logs. Because this project is a direct port, we need do our best to keep the behavior as identical as possible.

PRs reflecting changes made in `highlight.js` should **only** be made after `highlight.js` has had a release tagged with those changes. This project will always push out updates _after_ `highlight.js` releases.

We make no guarantees that the latest master `highlight.php` will be compatible with the latest master version of `highlight.js`.

## Project structure

The project contains the following folders:

1. [Highlight](#highlight)
2. [styles](#styles)
3. [demo](#demo)
4. [test](#test)
5. [tools](#tools)

## Highlight

This folder contains the main source and the following files (classes):

### Highlight.php (Highlight)

This is the one that does the highlighting for you and the one you'll probably look for.

### Language.php (Language)

Auxiliary class used in the Highlight class. Instances of these classes represent the rather complex structures of regular expressions needed to scan through programming code in order to highlight them.

You don't need this class.

### JsonRef.php (JsonRef)

Auxiliary class to decode JSON files containing path-based references. The language definition data from [highlight.js](http://www.highlightjs.org) is too complex to be described in ordinary JSON files. Therefore it was chosen to use [dojox.json.ref](https://dojotoolkit.org/reference-guide/1.9/dojox/json/ref.html) to export them. This class is able (to a very limited extend) to decode JSON data that was created with this [dojo](https://dojotoolkit.org) toolkit.

This class has a very distinct purpose and might be useful in other projects as well (and might be a good starting point for a new project ;) ).

### Autoloader.php (Autoloader)

A simple autoloader class that you possible might want or more likely not want to use. It is used for the tools and tests.

### the languages folder

This folder contains all language definitions: one JSON file for each language. These files are not hand coded but extracted from the original [highlight.js](http://www.highlightjs.org) project.

## Styles

These are the the CSS files needed to actually color the code. Not much to say about: these are just copied from the [highlight.js](https://github.com/isagalaev/highlight.js/tree/master/src/styles) project.

## Demo

This folder contains two demo pages that can be accessed through your browser.

### demo.php

A test page showing all supported languages and styles.

### compare.php

Much like [demo.php](#demo-php) but this page focuses on the comparison between _highlight.php_ and _highlight.js_. Both should yield the same results.

## Test

This folder contains the unit test for _highlight.php_. To run them, run _phpunit_ from this directory:

```bash
phpnunit .
```

Note that the following tests for _highlight.js_ were not rewritten for _highlight.php_:

### special explicitLanguage

Controlling language selection by setting an attribute to the containing `<div>` is irrelevant for _highlight.php_

### special customMarkup

In _highlight.js_, code may contain additional HTML markup like in the following PHP fragment: `$sum = <b>$a</b> + $b;`. Currently this is not supported by _highlight.php_ which can only highlight (unescaped) code. Also highlighting `<br>` (HTML break element) is not supported. _highlight.php_ does however support tab replacement (which defaults to 4 spaces).

### special noHighlight

There is no need to turn off highlighting through a class name on the code container.

### special buildClassName

_highlight.php_ does not modify class names of code containers.

## Tools

A collection of scripts that are used to extract data from the original [highlight.js](http://www.highlightjs.org) project.

The process of bringing in languages from highlight.js has been put into a single script. This script requires that you have cURL, PHP, and node.js installed on your machine. This is only necessary if you'd like to update languages or bring in new languages, it's **not** required for using this library in your application.

```bash
cd tools
bash process.sh
```
