# Changelog

All notable changes to `telegram` will be documented in this file

## 0.5.1 - 2020-12-07

- PHP 8 Support.

## 0.5.0 - 2020-09-08

- Add previous `ClientException` when constructing `CouldNotSendNotification` exception. PR [#86](https://github.com/laravel-notification-channels/telegram/pull/86).
- Add Laravel 8 Support. PR [#88](https://github.com/laravel-notification-channels/telegram/pull/88).
- Add Bot token per notification support. Closed [#84](https://github.com/laravel-notification-channels/telegram/issues/84).
- Add view file support for notification content. Closed [#82](https://github.com/laravel-notification-channels/telegram/issues/82).
- Add support to set HTTP Client.

## 0.4.1 - 2020-07-07

- Add Guzzle 7 Support. PR [#80](https://github.com/laravel-notification-channels/telegram/pull/80).

## 0.4.0 - 2020-06-02

- Add support to set custom api `base_uri` for web bridge.
- Revise README with instructions for Proxy or Bridge support.
- Revise on-demand notification instructions - Fixes [#72](https://github.com/laravel-notification-channels/telegram/issues/72).
- Fix typo in test.
- Remove redundant test.
- Remove exception when chat id isn't provided - PR [#75](https://github.com/laravel-notification-channels/telegram/pull/75).

## 0.3.0 - 2020-03-26

- Add ability to set param in `disableNotification` method.

## 0.2.0 - 2020-02-19

- Laravel 7 Support.
- Support response handling from Telegram.

## 0.1.1 - 2019-11-07

- Support PHP 7.1 and up.

## 0.1.0 - 2019-10-11

- New Helper Methods to work with file attachments.
- Code cleanup.
- Documentation updated with more examples and previews.
- Micro optimization and improvements.
- Typehint and return type declaration.
- Fixed tests.

## 0.0.6 - 2019-09-28

- Laravel 6 Support.
- Add Photo, Document, Audio, Location and other file notification type support.
- Token getter and setter.

## 0.0.5 - 2018-09-08

- Laravel 5.7 Support.
- Add ability to change button columns.

## 0.0.4 - 2018-02-08

- Laravel 5.6 Support.

## 0.0.3 - 2017-09-01

- Laravel 5.5 Support with Auto-Discovery.

## 0.0.2 - 2017-03-24

- Laravel 5.4 Support.

## 0.0.1 - 2016-08-14

- Initial Release.
