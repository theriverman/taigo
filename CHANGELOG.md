# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
 
## [Unreleased] - 2022-07-??
 
Taiga has changed it's authentication system to a more robust JWT. This requires the user to refresh their token every 24 hours (default setting.) If you're using Taigo in a long-running environment, such as, a webserver where your taigo.Client instance is preserved for days/weeks/months, then you need a way to keep your stored Token refreshed.

Taigo gets this task done automatically by polling a ticker in a goroutine and refreshing the stored tokens every 12 hours.

If you'd like to implement your own token refreshing mechanism, you have two options:
- implement your own routine based on `defaultTokenRefreshRoutine(c *Client, ticker *time.Ticker)`
- disable the `RefreshTokenRoutine` by calling `DisableAutomaticTokenRefresh()` and do the Token update your way.
  don't forget to update the contents of `Client.Headers`.
 
### Added
- `AuthService.RefreshAuthToken()` implemented
- New fields added to `Client`:
  - AutoRefreshDisabled
  - AutoRefreshTickerDuration
  - TokenRefreshTicker
- New methods added to `Client`:
  - `DisableAutomaticTokenRefresh()`

### Changed

- MAJOR Changed the signature of `*Client.AuthByToken(tokenType, token, refreshToken string) error`.
  It now requires `refreshToken` too.
- GitHub Workflows: Stepped go version to 1.18
- GitHub Workflows: Tests enabled on `feature_*` and `issue_*` branches
 
### Fixed

- [TAIGO-8](https://github.com/theriverman/taigo/issues/8)
  MAJOR Add missing Auth/Refresh auth token

## [1.4.0] - 2021-02-28
 
### Added

- Support for getting, creating and editing Taiga Issue objects
- Added custom attribute get/update example for user stories
- Added support for working easily in a specific project's scope via *ProjectService

### Changed

- TgObjectCustomAttributeValues are now exported and can be extended on-the-fly

### Fixed

- Struct members representing various agile points have been changed from regular int to float64

## [1.3.0] - 2020-10-04
 
### Added

- Support for custom attribute handling
- Support for working easily in a specific project's scope via `*ProjectService`

### Changed

### Fixed

## [1.2.2] - 2020-10-04
 
### Added

- A simplified init was implemented

### Changed

- update `contribute/main.go`

### Fixed

## [1.1.0] - 2020-09-23
  
### Added

### Changed
  
- Models have been improved for better accuracy
 
### Fixed

## [1.0.0] - 2020-09-18
 
### Added

- First version of the product

### Changed

### Fixed
