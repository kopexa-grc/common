# Changelog

## [1.11.0](https://github.com/kopexa-grc/common/compare/v1.10.0...v1.11.0) (2025-05-12)


### Features

* **blob:** implement Azure Blob storage integration with bucket oper… ([#36](https://github.com/kopexa-grc/common/issues/36)) ([1122452](https://github.com/kopexa-grc/common/commit/11224525d4d979d2f9c7442a9d5aec116f2cb749))
* **messaging:** implement NATS client and embedded server with authe… ([#37](https://github.com/kopexa-grc/common/issues/37)) ([9c621d5](https://github.com/kopexa-grc/common/commit/9c621d598c0ce07ae96c3bf869a626f2da91098b))

## [1.10.0](https://github.com/kopexa-grc/common/compare/v1.9.0...v1.10.0) (2025-05-12)


### Features

* **parser:** add query parameter parsing functions with unit tests ([#34](https://github.com/kopexa-grc/common/issues/34)) ([4fea3eb](https://github.com/kopexa-grc/common/commit/4fea3eb66c6b7dc057eae718f3e91e74ac55bd57))

## [1.9.0](https://github.com/kopexa-grc/common/compare/v1.8.0...v1.9.0) (2025-05-12)


### Features

* **khttp:** add CORS and metrics middleware with comprehensive tests and documentation ([#32](https://github.com/kopexa-grc/common/issues/32)) ([f7330b0](https://github.com/kopexa-grc/common/commit/f7330b03191f1de44abf2831ce2aa11409d196fe))

## [1.8.0](https://github.com/kopexa-grc/common/compare/v1.7.0...v1.8.0) (2025-05-11)


### Features

* **wellknown:** add Prometheus namespace constant and documentation ([#30](https://github.com/kopexa-grc/common/issues/30)) ([192b15d](https://github.com/kopexa-grc/common/commit/192b15d75d46f83c3f07b06d5e6a33ca792f538f))

## [1.7.0](https://github.com/kopexa-grc/common/compare/v1.6.0...v1.7.0) (2025-05-11)


### Features

* **graceful:** implement graceful shutdown mechanism with tests and documentation ([#29](https://github.com/kopexa-grc/common/issues/29)) ([63d72ac](https://github.com/kopexa-grc/common/commit/63d72ac79d15f012e6702de0247e9013a9003f28))
* **krn:** implement Kopexa Resource Name (KRN) system with JSON/YAML support and database integration ([#27](https://github.com/kopexa-grc/common/issues/27)) ([6591a93](https://github.com/kopexa-grc/common/commit/6591a9354c850a39b5ee497b350383073d208698))

## [1.6.0](https://github.com/kopexa-grc/common/compare/v1.5.0...v1.6.0) (2025-05-11)


### Features

* **logger:** implement buffered and colorized logging with environment detection ([#25](https://github.com/kopexa-grc/common/issues/25)) ([e29ca04](https://github.com/kopexa-grc/common/commit/e29ca0496f985942b4121bd98a466950a441ab05))

## [1.5.0](https://github.com/kopexa-grc/common/compare/v1.4.0...v1.5.0) (2025-05-11)


### Features

* soft delete docs ([#23](https://github.com/kopexa-grc/common/issues/23)) ([94611ea](https://github.com/kopexa-grc/common/commit/94611eaffd0c8fa98ab6a294ce46922f956f911c))

## [1.4.0](https://github.com/kopexa-grc/common/compare/v1.3.0...v1.4.0) (2025-05-11)


### Features

* add space-aware ID support to IDMixin ([#21](https://github.com/kopexa-grc/common/issues/21)) ([10aa53e](https://github.com/kopexa-grc/common/commit/10aa53e0900d31006b0b2941646d60863d8c4856))

## [1.3.0](https://github.com/kopexa-grc/common/compare/v1.2.1...v1.3.0) (2025-05-11)


### Features

* add error handling ([#5](https://github.com/kopexa-grc/common/issues/5)) ([2168ffb](https://github.com/kopexa-grc/common/commit/2168ffb0c5c3765e3caa92af1dd3c6ea2399129b))
* add ID and Audit mixins for enhanced entity management ([#15](https://github.com/kopexa-grc/common/issues/15)) ([fb44be7](https://github.com/kopexa-grc/common/commit/fb44be71e978b69d661b9f8b00907675cf0cd5f3))
* add release-please and dependabot configuration ([#2](https://github.com/kopexa-grc/common/issues/2)) ([0985f49](https://github.com/kopexa-grc/common/commit/0985f498217e0885d6a30aa2c6f79074e4fc133e))
* add to package for pointer management with 100% test coverage ([#4](https://github.com/kopexa-grc/common/issues/4)) ([bd41daf](https://github.com/kopexa-grc/common/commit/bd41daf9f6ae8c4e9f8d4ff72a7b4fd0a56f58bd))
* add-sessions-package ([#12](https://github.com/kopexa-grc/common/issues/12)) ([b2d4cdc](https://github.com/kopexa-grc/common/commit/b2d4cdc9e3d1062fd8c18a074badba8a55a517be))
* ctxutils ([#11](https://github.com/kopexa-grc/common/issues/11)) ([c07ed18](https://github.com/kopexa-grc/common/commit/c07ed186f1aacdc05dd6168f834eea71a0360e88))
* optimize makefile ([#6](https://github.com/kopexa-grc/common/issues/6)) ([896aed9](https://github.com/kopexa-grc/common/commit/896aed99de1195272aeafe040f4d838e21ee9163))
* **passwd:** use constants for Argon2 config, static errors, docs an… ([#9](https://github.com/kopexa-grc/common/issues/9)) ([a4e4dc7](https://github.com/kopexa-grc/common/commit/a4e4dc74dce227a26629410f7ed1cae46161c54e))
* totp implementation ([#13](https://github.com/kopexa-grc/common/issues/13)) ([768a62a](https://github.com/kopexa-grc/common/commit/768a62a86f452a990de98a74acb3fe43a19d15fc))


### Bug Fixes

* id mixin to use string ([#17](https://github.com/kopexa-grc/common/issues/17)) ([448c143](https://github.com/kopexa-grc/common/commit/448c143bb05be741261c06e010dacd4cc6deb8df))

## [1.2.1](https://github.com/kopexa-grc/common/compare/common-v1.2.0...common-v1.2.1) (2025-05-11)


### Bug Fixes

* id mixin to use string ([#17](https://github.com/kopexa-grc/common/issues/17)) ([448c143](https://github.com/kopexa-grc/common/commit/448c143bb05be741261c06e010dacd4cc6deb8df))

## [1.2.0](https://github.com/kopexa-grc/common/compare/common-v1.1.0...common-v1.2.0) (2025-05-10)


### Features

* add ID and Audit mixins for enhanced entity management ([#15](https://github.com/kopexa-grc/common/issues/15)) ([fb44be7](https://github.com/kopexa-grc/common/commit/fb44be71e978b69d661b9f8b00907675cf0cd5f3))

## [1.1.0](https://github.com/kopexa-grc/common/compare/common-v1.0.0...common-v1.1.0) (2025-05-10)


### Features

* totp implementation ([#13](https://github.com/kopexa-grc/common/issues/13)) ([768a62a](https://github.com/kopexa-grc/common/commit/768a62a86f452a990de98a74acb3fe43a19d15fc))

## 1.0.0 (2025-05-10)


### Features

* add error handling ([#5](https://github.com/kopexa-grc/common/issues/5)) ([2168ffb](https://github.com/kopexa-grc/common/commit/2168ffb0c5c3765e3caa92af1dd3c6ea2399129b))
* add release-please and dependabot configuration ([#2](https://github.com/kopexa-grc/common/issues/2)) ([0985f49](https://github.com/kopexa-grc/common/commit/0985f498217e0885d6a30aa2c6f79074e4fc133e))
* add to package for pointer management with 100% test coverage ([#4](https://github.com/kopexa-grc/common/issues/4)) ([bd41daf](https://github.com/kopexa-grc/common/commit/bd41daf9f6ae8c4e9f8d4ff72a7b4fd0a56f58bd))
* add-sessions-package ([#12](https://github.com/kopexa-grc/common/issues/12)) ([b2d4cdc](https://github.com/kopexa-grc/common/commit/b2d4cdc9e3d1062fd8c18a074badba8a55a517be))
* ctxutils ([#11](https://github.com/kopexa-grc/common/issues/11)) ([c07ed18](https://github.com/kopexa-grc/common/commit/c07ed186f1aacdc05dd6168f834eea71a0360e88))
* optimize makefile ([#6](https://github.com/kopexa-grc/common/issues/6)) ([896aed9](https://github.com/kopexa-grc/common/commit/896aed99de1195272aeafe040f4d838e21ee9163))
* **passwd:** use constants for Argon2 config, static errors, docs an… ([#9](https://github.com/kopexa-grc/common/issues/9)) ([a4e4dc7](https://github.com/kopexa-grc/common/commit/a4e4dc74dce227a26629410f7ed1cae46161c54e))
