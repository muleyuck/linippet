# Changelog

## [0.4.0](https://github.com/muleyuck/linippet/compare/v0.3.1...v0.4.0) (2026-03-29)


### Features

* add default value support for snippet arguments (${{arg_name:default}}) ([#25](https://github.com/muleyuck/linippet/issues/25)) ([db67226](https://github.com/muleyuck/linippet/commit/db6722653c43faf2241c592e06a70511e2b61ab2))
* crud snippet ([95f0af4](https://github.com/muleyuck/linippet/commit/95f0af43f060b1f9cc08df87d7bff2016944c5ac))
* read/write snippet to json ([9a1f46d](https://github.com/muleyuck/linippet/commit/9a1f46de2952b31aa66d050a8801077cb374c6d8))
* selectable snippet by tui ([ab39662](https://github.com/muleyuck/linippet/commit/ab39662c9cbebf9224f6c9cf644f6a44fb250640))


### Bug Fixes

* bugfix ([f1cff46](https://github.com/muleyuck/linippet/commit/f1cff46fb45eafaa98195db601bc58bf7a84a443))
* change comparision operator from "&gt;=" to "-ge" in initializer.bash ([#15](https://github.com/muleyuck/linippet/issues/15)) ([02057fd](https://github.com/muleyuck/linippet/commit/02057fd0f77739a4d94d2665d65b407740d4aeb7))
* close modal and return to list on cancel or Escape ([8c88d1e](https://github.com/muleyuck/linippet/commit/8c88d1e3d618d71c66e676d05bed3c263c1be44a))
* data file path when env_value is set ([b777e34](https://github.com/muleyuck/linippet/commit/b777e34515ef29834187e03a4e4f4f2591f4576d))
* Disable accept-line when selected snippet has no contents ([b9a6324](https://github.com/muleyuck/linippet/commit/b9a632431da89dbe86702356822c91276b7718d9))
* Disable the textView scrollbar so that it doesn’t accept tabIndex ([5301316](https://github.com/muleyuck/linippet/commit/5301316a555632ecdc3041d6c6acefed037ae60a))
* ensure correct validation of line ending in snippet ([#18](https://github.com/muleyuck/linippet/issues/18)) ([c36ac79](https://github.com/muleyuck/linippet/commit/c36ac79085a3765a6a0a810a8b79841ffce9981b))
* multiple args ([1ed9f9a](https://github.com/muleyuck/linippet/commit/1ed9f9a19e4d4a91eb7381a2937b485c3b6731eb))
* set ID on changing select index ([606dd09](https://github.com/muleyuck/linippet/commit/606dd09453e1624d1983278846fbd734296398b5))
* skipworktree app_version ([9ba7115](https://github.com/muleyuck/linippet/commit/9ba711596122c7837982fe51786ef3f3032917da))
* stop to tracking app_version ([e4fe1fc](https://github.com/muleyuck/linippet/commit/e4fe1fccc39de8c3767388c7563dc7a8166a281b))
* typo ([b4634cf](https://github.com/muleyuck/linippet/commit/b4634cfaac0330d23a5a0151274a5812869b47c3))
* version ([3419fcd](https://github.com/muleyuck/linippet/commit/3419fcd8bede54adc6ac81b50d2a350e93e756a0))

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
