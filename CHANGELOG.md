# Changelog

## [0.4.0](https://github.com/muleyuck/linippet/compare/v0.3.1...v0.4.0) (2026-04-11)


### Features

* align modal text to left and adjust padding ([41082c6](https://github.com/muleyuck/linippet/commit/41082c6b3b8819b23d0550db824521a943f7e920))
* improve modal text display with snippet preview and layout adjustments ([b01d719](https://github.com/muleyuck/linippet/commit/b01d7194a6e9b4b81227f57bcbd3bb91ea733d60))
* select all text in input fields when they receive focus ([aaa4d15](https://github.com/muleyuck/linippet/commit/aaa4d15985a97add0641910ae82a167681ee859a))
* select all text in input fields when they receive focus ([65ff9e7](https://github.com/muleyuck/linippet/commit/65ff9e7bb596b2ed06fefe09f0e3891a36e0bf2f))
* show snippet preview with resolved args in modal text ([f8b9c15](https://github.com/muleyuck/linippet/commit/f8b9c150b55e619f156572f39ef618a91b5769ad))

## [0.3.1](https://github.com/muleyuck/linippet/compare/v0.3.0...v0.3.1) (2026-03-04)


### Features

* add default value support for snippet arguments (`${{arg_name:default}}`) ([#25](https://github.com/muleyuck/linippet/pull/25)) ([db67226](https://github.com/muleyuck/linippet/commit/db67226))


### Bug Fixes

* close modal and return to list on cancel or Escape ([8c88d1e](https://github.com/muleyuck/linippet/commit/8c88d1e))
* disable the textView scrollbar so that it doesn't accept tabIndex ([5301316](https://github.com/muleyuck/linippet/commit/5301316))


## [0.3.0](https://github.com/muleyuck/linippet/compare/v0.2.4...v0.3.0) (2026-02-14)


### Features

* add cancellation support for snippet execution ([#24](https://github.com/muleyuck/linippet/pull/24)) ([6fd321e](https://github.com/muleyuck/linippet/commit/6fd321e))
* add install shell script and exclude Windows from releases ([#22](https://github.com/muleyuck/linippet/pull/22)) ([0b3327e](https://github.com/muleyuck/linippet/commit/0b3327e))


### Bug Fixes

* disable accept-line when selected snippet has no contents ([b9a6324](https://github.com/muleyuck/linippet/commit/b9a6324))
* ensure correct validation of line ending in snippet ([#18](https://github.com/muleyuck/linippet/pull/18)) ([c36ac79](https://github.com/muleyuck/linippet/commit/c36ac79))


### Performance Improvements

* refine fuzzy search scoring algorithm ([#24](https://github.com/muleyuck/linippet/pull/24)) ([6fd321e](https://github.com/muleyuck/linippet/commit/6fd321e))


## [0.2.4](https://github.com/muleyuck/linippet/compare/v0.2.3...v0.2.4) (2026-01-12)


### Bug Fixes

* ensure correct validation of line ending in snippet ([ff07364](https://github.com/muleyuck/linippet/commit/ff07364))


## [0.2.3](https://github.com/muleyuck/linippet/compare/v0.2.2...v0.2.3) (2025-12-21)


### Features

* indicate filtered item in list ([#7](https://github.com/muleyuck/linippet/pull/7)) ([64de59b](https://github.com/muleyuck/linippet/commit/64de59b))


### Bug Fixes

* change comparison operator from `>=` to `-ge` in initializer.bash ([#15](https://github.com/muleyuck/linippet/pull/15)) ([02057fd](https://github.com/muleyuck/linippet/commit/02057fd))


## [0.2.2](https://github.com/muleyuck/linippet/compare/v0.2.1...v0.2.2) (2025-05-04)


### Bug Fixes

* fix multiple args substitution ([1ed9f9a](https://github.com/muleyuck/linippet/commit/1ed9f9a))
* fix data file path when `LINIPPET_DATA` env var is set ([b777e34](https://github.com/muleyuck/linippet/commit/b777e34))


### Performance Improvements

* improve fuzzy search scoring algorithm ([59f821a](https://github.com/muleyuck/linippet/commit/59f821a))


## [0.2.1](https://github.com/muleyuck/linippet/compare/v0.2.0...v0.2.1) (2025-04-27)


### Features

* highlight characters matched by fuzzy search ([2a94a83](https://github.com/muleyuck/linippet/commit/2a94a83))


## [0.2.0](https://github.com/muleyuck/linippet/compare/v0.1.7...v0.2.0) (2025-03-30)


### Features

* fuzzy search for snippet list ([3419fcd](https://github.com/muleyuck/linippet/commit/3419fcd))


## [0.1.0](https://github.com/muleyuck/linippet/releases/tag/v0.1.0) (2025-03-30)


### Features

* interactive TUI to browse and execute snippets ([ab39662](https://github.com/muleyuck/linippet/commit/ab39662))
* CRUD operations for snippets (`create`, `edit`, `remove`) ([95f0af4](https://github.com/muleyuck/linippet/commit/95f0af4))
* read/write snippets to JSON file ([9a1f46d](https://github.com/muleyuck/linippet/commit/9a1f46d))
* dynamic argument placeholders (`${{arg_name}}`) in snippets ([9a1f46d](https://github.com/muleyuck/linippet/commit/9a1f46d))
* shell initialization via `linippet init` for bash/zsh integration ([22c1f80](https://github.com/muleyuck/linippet/commit/22c1f80))


### Bug Fixes

* fix set ID on changing select index ([606dd09](https://github.com/muleyuck/linippet/commit/606dd09))
