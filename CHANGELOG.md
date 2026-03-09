# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [2.0.0] - 2026-03-08

### Added
- Added full service wiring for status/type/custom-attribute endpoints:
  - points, priorities, severities, issue types
  - epic/issue/task/user story statuses
  - epic/issue/task/user story custom-attribute definition services
- Added missing operations for partially implemented resources:
  - tasks: edit/delete/get-by-ref improvements
  - issues: by-ref/delete/attachments
  - wiki: list/create/get/get-by-slug/edit/delete/render/attachments
  - webhooks: test endpoint fix, log helpers, resend helpers
- Added coverage for previously stubbed endpoint groups:
  - applications/application-tokens, searches, user-storage
  - project-templates, memberships/invitations, wiki-links
  - history, notify-policies, contact, feedback
  - export/import, timelines, locales, importers
  - contrib-plugins, objects-summary
- Added `MIGRATION.md` for v1 -> v2 upgrade guidance.
- Added unit tests for key v2 correctness areas (`v2_proposals_test.go`).

### Changed
- **BREAKING:** module path moved to `github.com/theriverman/taigo/v2`.
- **BREAKING:** normalised multiple method signatures and added `Update(...)` aliases where relevant.
- Pagination header semantics now match Taiga behaviour:
  - disabling pagination sets `x-disable-pagination`
  - enabling pagination removes the header
- Query semantics updated:
  - optional booleans use pointer-bool fields
  - task tag filter uses comma-delimited encoding helper
  - project order_by now serialises correctly
- Non-2xx responses now return typed `APIError`.
- Integration tests are now opt-in via `TAIGO_RUN_INTEGRATION_TESTS=1`.

### Fixed
- Fixed webhook test route to `POST /webhooks/{id}/test`.
- Fixed custom attribute value payload keys for issue/task/user story value DTOs.
- Fixed watched/liked user endpoint decode shape to list responses.
- Fixed token refresh shutdown behaviour and nil-safe disable logic.
- Removed placeholder-only stub files by implementing concrete service surfaces.
 
## [1.5.0] - 2022-08-03
 
Taiga has changed its authentication system to a more sophisticated JWT implementation. This requires the user to refresh their token every 24 hours (default setting). If you're using Taigo in a system which tends to run for longer than 24 hours, such as, a webserver where your `taigo.Client` instance is preserved for days/weeks/months, then you need a way to keep your stored token fresh.

Taigo gets this task done automatically by polling a ticker in a goroutine and refreshing the stored tokens every 12 hours.

If you'd like to implement your own token refreshing mechanism, you have two options:
- implement your own routine based on `defaultTokenRefreshRoutine(c *Client, ticker *time.Ticker)`
- disable the `RefreshTokenRoutine` by calling `DisableAutomaticTokenRefresh()` and do the Token update your way (don't forget to update the contents of `Client.Headers`).
 
### Added
- `AuthService.RefreshAuthToken()` implemented
- New fields added to `Client`:
  - `AutoRefreshDisabled`
  - `AutoRefreshTickerDuration`
  - `TokenRefreshTicker`
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
