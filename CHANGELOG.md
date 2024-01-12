# Changelog

## [2.0.1](https://github.com/rickstaa/warnings-discord-bot/compare/v2.0.0...v2.0.1) (2024-01-12)


### Bug Fixes

* fix roles intent ([#26](https://github.com/rickstaa/warnings-discord-bot/issues/26)) ([de8a26b](https://github.com/rickstaa/warnings-discord-bot/commit/de8a26be97877cd792740793dc1c53d2c8a3bffa))


### Performance Improvements

* add goroutines ([#28](https://github.com/rickstaa/warnings-discord-bot/issues/28)) ([2318952](https://github.com/rickstaa/warnings-discord-bot/commit/2318952014f27f69567186a6d7a0768a7d500786))

## [2.0.0](https://github.com/rickstaa/warnings-discord-bot/compare/v1.1.2...v2.0.0) (2024-01-12)


### ⚠ BREAKING CHANGES

* add ability to specify regex ([#20](https://github.com/rickstaa/warnings-discord-bot/issues/20))

### Features

* add ability to only trigger on new members ([#23](https://github.com/rickstaa/warnings-discord-bot/issues/23)) ([29a9ef5](https://github.com/rickstaa/warnings-discord-bot/commit/29a9ef505e9b376a9c52b4b5af88b90a06dd2929))
* add ability to specify regex ([#20](https://github.com/rickstaa/warnings-discord-bot/issues/20)) ([ff28c63](https://github.com/rickstaa/warnings-discord-bot/commit/ff28c639ab6e722a7ab532b613a7e1f33104f3b4))
* add welcome warning feature ([#25](https://github.com/rickstaa/warnings-discord-bot/issues/25)) ([cbdc4ab](https://github.com/rickstaa/warnings-discord-bot/commit/cbdc4ab27e16baf00556b21aa6a082c302180e1b))
* change reply to embed reply ([#22](https://github.com/rickstaa/warnings-discord-bot/issues/22)) ([847fe17](https://github.com/rickstaa/warnings-discord-bot/commit/847fe178cfeaaf9dfbbcc236c5b740b4e9e314d9))


### Bug Fixes

* ensure internal links don't trigger warnings ([#24](https://github.com/rickstaa/warnings-discord-bot/issues/24)) ([0039713](https://github.com/rickstaa/warnings-discord-bot/commit/00397134f96eee1a61725f25fdf79828a2faa6e4))

## [1.1.2](https://github.com/rickstaa/warnings-discord-bot/compare/v1.1.1...v1.1.2) (2024-01-09)


### Bug Fixes

* ensure the role parameters are linked to their keyword rule ([#18](https://github.com/rickstaa/warnings-discord-bot/issues/18)) ([d4bb1a4](https://github.com/rickstaa/warnings-discord-bot/commit/d4bb1a4d9d842f45fd8590f04484547318c98df1))
* fix empty required roles not triggered bug ([#15](https://github.com/rickstaa/warnings-discord-bot/issues/15)) ([4f8175b](https://github.com/rickstaa/warnings-discord-bot/commit/4f8175b2b6a229d43a2a6d2566d51692ff2d897f))

## [1.1.1](https://github.com/rickstaa/warnings-discord-bot/compare/v1.1.0...v1.1.1) (2024-01-08)


### Bug Fixes

* fix no roles condition ([#11](https://github.com/rickstaa/warnings-discord-bot/issues/11)) ([683f962](https://github.com/rickstaa/warnings-discord-bot/commit/683f962fcf5f74a86bc52de2287c7c4e3aca797a))

## [1.1.0](https://github.com/rickstaa/warnings-discord-bot/compare/v1.0.0...v1.1.0) (2024-01-08)


### Features

* add 'excluded_roles' config parameter ([#9](https://github.com/rickstaa/warnings-discord-bot/issues/9)) ([0a5f661](https://github.com/rickstaa/warnings-discord-bot/commit/0a5f6615ac3b4a13cd47cb75499818df1a7aa354))

## 1.0.0 (2024-01-08)


### Features

* add 'required_roles' configuration value ([1a989a8](https://github.com/rickstaa/warnings-discord-bot/commit/1a989a8cbc93214893b3f291180489b48dc93957))
* add docker deployment ([#8](https://github.com/rickstaa/warnings-discord-bot/issues/8)) ([107184b](https://github.com/rickstaa/warnings-discord-bot/commit/107184b059004ddb52eaa1f3fd2fa81c3836e6f0))
* add initial bot code ([720bc56](https://github.com/rickstaa/warnings-discord-bot/commit/720bc56350a4458d0ab814a044df101ca162845a))
* change bot message to reply ([ddb57a3](https://github.com/rickstaa/warnings-discord-bot/commit/ddb57a380bbe44f3b3f9004f7ea9dbd0277a658e))
