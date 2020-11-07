# ramsey/collection Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

## [1.1.1] - 2020-09-10

### Fixed

* Fixed broken `AbstractCollection::map()` implementation.

## [1.1.0] - 2020-08-10

### Fixed

* Fixed `AbstractCollection::diff()`, `AbstractCollection::intersect()` and
  `AbstractCollection::merge()` when used with Generic collections.
* Fixed `AbstractCollection::diff()` and `AbstractCollection::intersect()`
  returning inconsistent results when used on collections containing objects.
* Removed warning about deprecated dependency when running `composer install`

## [1.0.1] - 2020-01-04

### Fixed

* Fixed `AbstractCollection::offsetSet()` so that it uses the provided `$offset`
  when setting `$value` in the array.

## [1.0.0] - 2018-12-31

### Added

* Added support for *queue* data structures to represent collections of ordered
  entities. Together with *double-ended queues* (a.k.a. *deques*),
  first-in-first-out (FIFO), last-in-first-out (LIFO), and other queue and stack
  behaviors may be implemented. This functionality includes interfaces
  `QueueInterface` and `DoubleEndedQueueInterface` and classes `Queue` and
  `DoubleEndedQueue`.
* Added support for *set* data structures, representing collections that cannot
  contain any duplicated elements; includes classes `AbstractSet` and `Set`.
* Added support for *typed map* data structures to represent maps of elements
  where both keys and values have specified data types; includes
  `TypedMapInterface` and the classes `AbstractTypedMap` and `TypedMap`.
* Added new manipulation and analyze methods for collections: `column()`,
  `first()`, `last()`, `sort()`, `filter()`, `where()`, `map()`, `diff()`,
  `intersect()`, and `merge()`. See [CollectionInterface](https://github.com/ramsey/collection/blob/master/src/CollectionInterface.php)
  for more information.
* Added the following new exceptions specific to the ramsey/collection library:
  `CollectionMismatchException`, `InvalidArgumentException`,
  `InvalidSortOrderException`, `NoSuchElementException`, `OutOfBoundsException`,
  `UnsupportedOperationException`, and `ValueExtractionException`.

### Changed

* Minimum PHP version supported is 7.2.
* Strict types are enforced throughout.

### Removed

* Removed support for HHVM.

### Security

* Fixed possible exploit using `AbstractArray::unserialize()`
  (see [#47](https://github.com/ramsey/collection/issues/47)).

## [0.3.0] - 2016-05-23

### Added

* Added `MapInterface::keys()` method to return the keys from a `MapInterface`
  object. This was added to the `AbstractMap` class.

### Removed

* Removed `getType()` and constructor methods from `AbstractCollection`. Children
  of `AbstractCollection` must now implement `getType()`, which should return a
  string value that defines the data type of items for the collection.

### Fixed

* Improve error messages in exceptions when `Collection` and `NamedParameterMap`
  items fail type checks.

## [0.2.1] - 2016-02-22

### Fixed

* Allow non-strict checking of values in typed collections.

## [0.2.0] - 2016-02-05

### Added

* Support typed collections.

## [0.1.0] - 2015-10-27

### Added

* Support generic arrays and maps.

[Unreleased]: https://github.com/ramsey/collection/compare/1.1.0...HEAD
[1.1.0]: https://github.com/ramsey/collection/compare/1.0.1...1.1.0
[1.0.1]: https://github.com/ramsey/collection/compare/1.0.0...1.0.1
[1.0.0]: https://github.com/ramsey/collection/compare/0.3.0...1.0.0
[0.3.0]: https://github.com/ramsey/collection/compare/0.2.1...0.3.0
[0.2.1]: https://github.com/ramsey/collection/compare/0.2.0...0.2.1
[0.2.0]: https://github.com/ramsey/collection/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/ramsey/collection/commits/0.1.0
