# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2020-05-11

### Added
- CORS headers are better cachable now, with correct Vary headers (#70, #74)

### Changed
- CORS headers are added to non-Origin requests when possible (#73)
- Requests are no longer blocked by the server, only by the browser (#70)
