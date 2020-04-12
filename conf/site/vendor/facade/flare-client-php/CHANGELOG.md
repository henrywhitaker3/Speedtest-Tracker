# Changelog

All notable changes to `flare-client-php` will be documented in this file

## 1.3.2 - 2020-03-02

- allow L7

## 1.3.1 - 2019-12-15

- allow var-dumper v5.0

## 1.3.0 - 2019-11-27

- Allow custom grouping types

## 1.2.1 - 2019-11-19

- Let `registerFlareHandlers` return $this

## 1.2.0 - 2019-11-19

- Add `registerFlareHandlers` method to register error and exception handlers in non-Laravel applications
- Fix get requests with query parameters (#4)

## 1.1.2 - 2019-11-08

- Ignore invalid mime type detection issues

## 1.1.1 - 2019-10-07

- Wrap filesize detection in try-catch block

## 1.1.0 - 2019-09-27

- Add ability to log messages

## 1.0.4 - 2019-09-11

- Fixes an issue when sending exceptions inside a queue worker

## 1.0.3 - 2019-09-05

- Ensure valid session data

## 1.0.2 - 2019-09-05

- Fix error when uploading multiple files using an array name

## 1.0.1 - 2019-09-02

- Fix issue with uploaded files in request context

## 1.0.0 - 2019-08-30

- initial release
