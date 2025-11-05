# Changelog

## [1.58.0](https://github.com/kopexa-grc/common/compare/v1.57.0...v1.58.0) (2025-11-05)


### Features

* implement WithIncludeArchived function to manage archived records in context ([37aa182](https://github.com/kopexa-grc/common/commit/37aa182b55114a68ff5d9aa923f250b7ace23cb2))


### Bug Fixes

* rename ArchiveSkipKey to SkipKey for consistency in context handling ([a09bdc1](https://github.com/kopexa-grc/common/commit/a09bdc1fd5b01e6a5b58ce54b8036a63fc262fdf))

## [1.57.0](https://github.com/kopexa-grc/common/compare/v1.56.1...v1.57.0) (2025-11-04)


### Features

* add ContextualTuples to AccessCheck and ListRequest for enhanced context handling ([1356cf2](https://github.com/kopexa-grc/common/commit/1356cf246e37d45db30c2639a8767f1255197ad4))

## [1.56.1](https://github.com/kopexa-grc/common/compare/v1.56.0...v1.56.1) (2025-11-03)


### Bug Fixes

* CONSISTENCYPREFERENCE_HIGHER_CONSISTENCY to fga ListObjects ([3117d5d](https://github.com/kopexa-grc/common/commit/3117d5d067e211e1ef515432d9024bfbd98667e0))
* fga list tests ([f83f79d](https://github.com/kopexa-grc/common/commit/f83f79d0589cb6afa48b5a37e54cfc71646d3d3a))

## [1.56.0](https://github.com/kopexa-grc/common/compare/v1.55.0...v1.56.0) (2025-10-01)


### Features

* reset token ([a939050](https://github.com/kopexa-grc/common/commit/a93905027d0bab07478a17bf712fb30b3fc25e34))

## [1.55.0](https://github.com/kopexa-grc/common/compare/v1.54.0...v1.55.0) (2025-09-30)


### Features

* ptr package ([fbee33f](https://github.com/kopexa-grc/common/commit/fbee33f3b2e6c8674098bcf63a10fca064e7cfd0))

## [1.54.0](https://github.com/kopexa-grc/common/compare/v1.53.1...v1.54.0) (2025-09-30)


### Features

* expose expand api on fga ([ea410fa](https://github.com/kopexa-grc/common/commit/ea410fa5d3f8afb45976aab76d2ae48323201c61))
* expose expand api on fga ([5b2bb83](https://github.com/kopexa-grc/common/commit/5b2bb83453b4fcbf22638585aa0442972b7e3e47))
* **fga:** implement ListUsersWithAccess function and add tests ([59db99d](https://github.com/kopexa-grc/common/commit/59db99dce4a94847d39fe4b872502f4da2cedce0))

## [1.53.1](https://github.com/kopexa-grc/common/compare/v1.53.0...v1.53.1) (2025-09-23)


### Bug Fixes

* **logging:** change log level from Info to Debug in listObjects function ([d24349f](https://github.com/kopexa-grc/common/commit/d24349f397357392b4eecc1198724fefa897e1cd))

## [1.53.0](https://github.com/kopexa-grc/common/compare/v1.52.0...v1.53.0) (2025-09-10)


### Features

* **fga:** add logging for listed objects in listObjects function ([d41b275](https://github.com/kopexa-grc/common/commit/d41b2756104a9a7cc0f35985803f6f21558a3e53))

## [1.52.0](https://github.com/kopexa-grc/common/compare/v1.51.0...v1.52.0) (2025-08-05)


### Features

* **to:** add BoolValue, StringValue, and TimeValue utility functions ([8ba6b1c](https://github.com/kopexa-grc/common/commit/8ba6b1c8d6cb5abb9fe355f528fbf834ff7517bc))

## [1.51.0](https://github.com/kopexa-grc/common/compare/v1.50.0...v1.51.0) (2025-07-23)


### Features

* **tokens:** add verification token creation and validation with email checks ([0d8035f](https://github.com/kopexa-grc/common/commit/0d8035f1742275c396ace9b39b810d7d6323c098))

## [1.50.0](https://github.com/kopexa-grc/common/compare/v1.49.0...v1.50.0) (2025-07-11)


### Features

* **ent:** add JSON scalar support and related utilities ([318c8f2](https://github.com/kopexa-grc/common/commit/318c8f2eb28cab746b758d7e0fce8df90611db56))


### Bug Fixes

* prompting ([78d1292](https://github.com/kopexa-grc/common/commit/78d12923e7ed02aca0702540d185fa07166e8bb3))
* remove truncate ([dfab4a0](https://github.com/kopexa-grc/common/commit/dfab4a051f09af25a5e5d5f53ea40b2e39b5804f))

## [1.49.0](https://github.com/kopexa-grc/common/compare/v1.48.0...v1.49.0) (2025-07-11)


### Features

* **mixins:** add TagMixin with comprehensive documentation and tests ([765904b](https://github.com/kopexa-grc/common/commit/765904b3ed9082b15bc3912799c5404da726dc90))

## [1.48.0](https://github.com/kopexa-grc/common/compare/v1.47.0...v1.48.0) (2025-07-07)


### Features

* **validation:** add comprehensive URL and domain validation tests ([e2e7e10](https://github.com/kopexa-grc/common/commit/e2e7e1085df9bf2b28be3c7b47ba2d01d289648a))

## [1.47.0](https://github.com/kopexa-grc/common/compare/v1.46.0...v1.47.0) (2025-07-03)


### Features

* **summarizer:** add NewFromLLM function for LLM integration ([18d8573](https://github.com/kopexa-grc/common/commit/18d857394adff3917cd07b027fc4fbf47404adfc))

## [1.46.0](https://github.com/kopexa-grc/common/compare/v1.45.1...v1.46.0) (2025-06-26)


### Features

* **gql:** add GraphQL pagination and constant handling ([15c55bc](https://github.com/kopexa-grc/common/commit/15c55bc930a22f6b17bbd0e6b2d38c3e80e33acd))

## [1.45.1](https://github.com/kopexa-grc/common/compare/v1.45.0...v1.45.1) (2025-06-26)


### Bug Fixes

* blob docs ([105ea40](https://github.com/kopexa-grc/common/commit/105ea40064c012f4a0e274269fe4f5de24573457))

## [1.45.0](https://github.com/kopexa-grc/common/compare/v1.44.1...v1.45.0) (2025-06-26)


### Features

* add LLM summarization package with multi-provider support ([848b7c3](https://github.com/kopexa-grc/common/commit/848b7c35c4359aaf50da922f7227ece7615cf3cd))

## [1.44.1](https://github.com/kopexa-grc/common/compare/v1.44.0...v1.44.1) (2025-06-25)


### Bug Fixes

* solves a duplication error ([e81397d](https://github.com/kopexa-grc/common/commit/e81397dbc0356ef0d8f5c2bbdf620c335ca862f5))

## [1.44.0](https://github.com/kopexa-grc/common/compare/v1.43.1...v1.44.0) (2025-06-24)


### Features

* implement a no-op validation method for Address struct ([#111](https://github.com/kopexa-grc/common/issues/111)) ([24c1133](https://github.com/kopexa-grc/common/commit/24c1133064b5cfa8860c99ef10beea984b2e29fb))

## [1.43.1](https://github.com/kopexa-grc/common/compare/v1.43.0...v1.43.1) (2025-06-24)


### Bug Fixes

* remove validation logic from Address struct and related tests ([#109](https://github.com/kopexa-grc/common/issues/109)) ([6ac7687](https://github.com/kopexa-grc/common/commit/6ac7687411f83fa812000e6aecf0180cc59b543a))

## [1.43.0](https://github.com/kopexa-grc/common/compare/v1.42.1...v1.43.0) (2025-06-23)


### Features

* add duration parsing and formatting functionality with tests ([#107](https://github.com/kopexa-grc/common/issues/107)) ([5e20220](https://github.com/kopexa-grc/common/commit/5e20220f041b00485735f9ce181fb2ed62b32420))

## [1.42.1](https://github.com/kopexa-grc/common/compare/v1.42.0...v1.42.1) (2025-06-21)


### Bug Fixes

* nil pointer for metrics if not provided ([dfe2c38](https://github.com/kopexa-grc/common/commit/dfe2c3819d3c24076c05687f163ff9f282f8d0e1))

## [1.42.0](https://github.com/kopexa-grc/common/compare/v1.41.1...v1.42.0) (2025-06-21)


### Features

* Implement new Reader and Writer interfaces with range reading and writing capabilities ([ed43e11](https://github.com/kopexa-grc/common/commit/ed43e1140426793d40cd0e24a84fac11f12061df))


### Bug Fixes

* Update blob package for improved error handling and code clarity ([5a9fe3a](https://github.com/kopexa-grc/common/commit/5a9fe3a8cbdbc43f0f9737ff096657ec62db16e9))

## [1.41.1](https://github.com/kopexa-grc/common/compare/v1.41.0...v1.41.1) (2025-06-19)


### Bug Fixes

* pointer usage bug ([ffaca51](https://github.com/kopexa-grc/common/commit/ffaca51d5bea2694d45a8664ef3226397ccc432e))

## [1.41.0](https://github.com/kopexa-grc/common/compare/v1.40.0...v1.41.0) (2025-06-16)


### Features

* Add wildcard support and public access tuple creation ([e05ac49](https://github.com/kopexa-grc/common/commit/e05ac498aa684cff2c3ef62c55a94a8281173221))

## [1.40.0](https://github.com/kopexa-grc/common/compare/v1.39.0...v1.40.0) (2025-06-12)


### Features

* risk audit entity ([390d960](https://github.com/kopexa-grc/common/commit/390d96031d5f8d0d01caf6438936fb4fb3c68367))

## [1.39.0](https://github.com/kopexa-grc/common/compare/v1.38.0...v1.39.0) (2025-06-10)


### Features

* **localization:** Implement localization utilities and tests ([6f098ea](https://github.com/kopexa-grc/common/commit/6f098ea05fe4e098fbc3d25407253bbb7091c3d0))

## [1.38.0](https://github.com/kopexa-grc/common/compare/v1.37.0...v1.38.0) (2025-06-08)


### Features

* **types:** Add compliance types and GraphQL marshaling/unmarshaling tests ([5ad6dd4](https://github.com/kopexa-grc/common/commit/5ad6dd4b6141b589c3cf32b32497e0f6c5c8b11b))

## [1.37.0](https://github.com/kopexa-grc/common/compare/v1.36.0...v1.37.0) (2025-06-07)


### Features

* **audit:** Add GraphQL annotations to audit fields for improved mutation handling ([54a0157](https://github.com/kopexa-grc/common/commit/54a01576413d94fe8461437faaf73c008856ffd2))

## [1.36.0](https://github.com/kopexa-grc/common/compare/v1.35.0...v1.36.0) (2025-05-30)


### Features

* added display id annotations ([#96](https://github.com/kopexa-grc/common/issues/96)) ([0d004c2](https://github.com/kopexa-grc/common/commit/0d004c267345a6dd20142a5b7df6343f60424d7d))

## [1.35.0](https://github.com/kopexa-grc/common/compare/v1.34.0...v1.35.0) (2025-05-24)


### Features

* fga create store ([#94](https://github.com/kopexa-grc/common/issues/94)) ([453f78a](https://github.com/kopexa-grc/common/commit/453f78a2515924ed1b6425a78acf253429c5d6e1))

## [1.34.0](https://github.com/kopexa-grc/common/compare/v1.33.0...v1.34.0) (2025-05-23)


### Features

* **auth:** Add organization subscription context and price struct ([#92](https://github.com/kopexa-grc/common/issues/92)) ([b244ab0](https://github.com/kopexa-grc/common/commit/b244ab0d5b28b5fc62c9c983abb8cc4fd044cc6a))

## [1.33.0](https://github.com/kopexa-grc/common/compare/v1.32.0...v1.33.0) (2025-05-22)


### Features

* **ci:** Add GitHub Actions workflow for linting pull requests ([#90](https://github.com/kopexa-grc/common/issues/90)) ([d47a411](https://github.com/kopexa-grc/common/commit/d47a411807f7e34a66e8a1192171b43ad153fb1c))

## [1.32.0](https://github.com/kopexa-grc/common/compare/v1.31.0...v1.32.0) (2025-05-22)


### Features

* **fga:** Enhance error handling and add ListAccess functionality ([#87](https://github.com/kopexa-grc/common/issues/87)) ([f198146](https://github.com/kopexa-grc/common/commit/f198146f0ef43ddcbcde81e22158c6cc77cfab7a))

## [1.31.0](https://github.com/kopexa-grc/common/compare/v1.30.0...v1.31.0) (2025-05-22)


### Features

* **fga:** Add batch access check functionality and ULID correlation ID generation ([#85](https://github.com/kopexa-grc/common/issues/85)) ([28b9b32](https://github.com/kopexa-grc/common/commit/28b9b3292ff2e9b1cdaa9a2a6734e594257b03e5))

## [1.30.0](https://github.com/kopexa-grc/common/compare/v1.29.0...v1.30.0) (2025-05-22)


### Features

* **auth:** Add AuthenticationMethodsReferences and constants for RFC8176 compliance ([#83](https://github.com/kopexa-grc/common/issues/83)) ([451bc79](https://github.com/kopexa-grc/common/commit/451bc793dfd04ea486704064fd41e779c636dd53))

## [1.29.0](https://github.com/kopexa-grc/common/compare/v1.28.1...v1.29.0) (2025-05-22)


### Features

* **auth:** Implement AuthLevel type with methods for string conversion and GraphQL marshaling ([#81](https://github.com/kopexa-grc/common/issues/81)) ([a822d8f](https://github.com/kopexa-grc/common/commit/a822d8fdd9b6739ca427db822f0871a7d94aea2a))

## [1.28.1](https://github.com/kopexa-grc/common/compare/v1.28.0...v1.28.1) (2025-05-21)


### Bug Fixes

* **types:** Change MarshalGQL method receivers from pointer to value for Reference, ResponseData, and ResponseMeta types ([#79](https://github.com/kopexa-grc/common/issues/79)) ([d23832a](https://github.com/kopexa-grc/common/commit/d23832a44ebb8f58b4c26d5e897c470ba8c02c40))

## [1.28.0](https://github.com/kopexa-grc/common/compare/v1.27.0...v1.28.0) (2025-05-21)


### Features

* **auth:** Add ActorType methods and error handling ([#75](https://github.com/kopexa-grc/common/issues/75)) ([242589e](https://github.com/kopexa-grc/common/commit/242589ee51e35d121b4721b11d98e6708d2d9198))


### Bug Fixes

* **types:** Remove unnecessary validation in UnmarshalGQL method of Author type ([#73](https://github.com/kopexa-grc/common/issues/73)) ([0824c20](https://github.com/kopexa-grc/common/commit/0824c205e3be1dcbbe6393ecc17d0fe97b8563c9))
* **types:** Update MarshalGQL method to use pointer receiver for ResponseMeta ([#78](https://github.com/kopexa-grc/common/issues/78)) ([8e015c2](https://github.com/kopexa-grc/common/commit/8e015c2c2c60e7b5ae7abd36b2e10f1521f65ad2))

## [1.27.0](https://github.com/kopexa-grc/common/compare/v1.26.1...v1.27.0) (2025-05-21)


### Features

* **auth:** Add ActorType methods and error handling ([#75](https://github.com/kopexa-grc/common/issues/75)) ([242589e](https://github.com/kopexa-grc/common/commit/242589ee51e35d121b4721b11d98e6708d2d9198))

## [1.26.1](https://github.com/kopexa-grc/common/compare/v1.26.0...v1.26.1) (2025-05-21)


### Bug Fixes

* **types:** Remove unnecessary validation in UnmarshalGQL method of Author type ([#73](https://github.com/kopexa-grc/common/issues/73)) ([0824c20](https://github.com/kopexa-grc/common/commit/0824c205e3be1dcbbe6393ecc17d0fe97b8563c9))

## [1.26.0](https://github.com/kopexa-grc/common/compare/v1.25.0...v1.26.0) (2025-05-19)


### Features

* author type ([#70](https://github.com/kopexa-grc/common/issues/70)) ([b9d9f38](https://github.com/kopexa-grc/common/commit/b9d9f38cf536270459e7ca45574aa35dedbdaf3c))

## [1.25.0](https://github.com/kopexa-grc/common/compare/v1.24.0...v1.25.0) (2025-05-19)


### Features

* **types:** Verbesserte Reference Implementierung ([#68](https://github.com/kopexa-grc/common/issues/68)) ([12ab449](https://github.com/kopexa-grc/common/commit/12ab4492149056a198e8e668642e8dd32c69faa8))

## [1.24.0](https://github.com/kopexa-grc/common/compare/v1.23.1...v1.24.0) (2025-05-19)


### Features

* **types:** Verbesserte Metadata Implementierung ([#66](https://github.com/kopexa-grc/common/issues/66)) ([bded201](https://github.com/kopexa-grc/common/commit/bded2019630759d4902d37395c5a5fa744ef8a4d))

## [1.23.1](https://github.com/kopexa-grc/common/compare/v1.23.0...v1.23.1) (2025-05-19)


### Bug Fixes

* **types:** Konsistente Receiver-Namen in ResponseData ([#64](https://github.com/kopexa-grc/common/issues/64)) ([4a02f45](https://github.com/kopexa-grc/common/commit/4a02f45bb0a111ff82b31715e211ee63773d6cea))

## [1.23.0](https://github.com/kopexa-grc/common/compare/v1.22.0...v1.23.0) (2025-05-19)


### Features

* added risk type ([#62](https://github.com/kopexa-grc/common/issues/62)) ([421569b](https://github.com/kopexa-grc/common/commit/421569b07bbf0c6977cd19fc502b7413414d046d))

## [1.22.0](https://github.com/kopexa-grc/common/compare/v1.21.0...v1.22.0) (2025-05-18)


### Features

* **types:** Verbesserte Dokumentation und Code-Qualität ([#60](https://github.com/kopexa-grc/common/issues/60)) ([0e9bd76](https://github.com/kopexa-grc/common/commit/0e9bd766c6c7577ed45ba116d2a9fdc16e4e2121))

## [1.21.0](https://github.com/kopexa-grc/common/compare/v1.20.0...v1.21.0) (2025-05-17)


### Features

* org and space context ([#58](https://github.com/kopexa-grc/common/issues/58)) ([52b66aa](https://github.com/kopexa-grc/common/commit/52b66aa4017dd35c53035e75368ae6fc3fca79b7))

## [1.20.0](https://github.com/kopexa-grc/common/compare/v1.19.0...v1.20.0) (2025-05-16)


### Features

* iam tokens ([#56](https://github.com/kopexa-grc/common/issues/56)) ([b7944b8](https://github.com/kopexa-grc/common/commit/b7944b87d13905a22e2ede40ab90366b11bfce84))

## [1.19.0](https://github.com/kopexa-grc/common/compare/v1.18.0...v1.19.0) (2025-05-16)


### Features

* fgalist ([#53](https://github.com/kopexa-grc/common/issues/53)) ([bed7a00](https://github.com/kopexa-grc/common/commit/bed7a00d852f4112d1aaf2f4af545a3b5d30e310))

## [1.18.0](https://github.com/kopexa-grc/common/compare/v1.17.0...v1.18.0) (2025-05-15)


### Features

* **fga:** add Context field to AccessCheck for conditional relationships ([#51](https://github.com/kopexa-grc/common/issues/51)) ([4664a5f](https://github.com/kopexa-grc/common/commit/4664a5f9b03718cfdc992ab0ba43639f46bdac08))

## [1.17.0](https://github.com/kopexa-grc/common/compare/v1.16.0...v1.17.0) (2025-05-15)


### Features

* **fga:** update Entity.String method to handle wildcard entities ([#49](https://github.com/kopexa-grc/common/issues/49)) ([ee521e1](https://github.com/kopexa-grc/common/commit/ee521e11aa868c9b5448858ae40ee1d75748c6ea))

## [1.16.0](https://github.com/kopexa-grc/common/compare/v1.15.0...v1.16.0) (2025-05-14)


### Features

* **fga:** add API token authentication support ([#47](https://github.com/kopexa-grc/common/issues/47)) ([8ca632f](https://github.com/kopexa-grc/common/commit/8ca632fa87f38c95465c6c998d457ffbceeefa2d))

## [1.15.0](https://github.com/kopexa-grc/common/compare/v1.14.0...v1.15.0) (2025-05-13)


### Features

* added fga client to common library ([#45](https://github.com/kopexa-grc/common/issues/45)) ([21ba8b1](https://github.com/kopexa-grc/common/commit/21ba8b13507252a65b7a166a9e7186e9d700fe88))

## [1.14.0](https://github.com/kopexa-grc/common/compare/v1.13.0...v1.14.0) (2025-05-13)


### Features

* **mixins:** add TranslationMixin for handling locale and parent field indexing ([#43](https://github.com/kopexa-grc/common/issues/43)) ([0664217](https://github.com/kopexa-grc/common/commit/0664217580f336c8d14fed445f2fd8be26af5b28))

## [1.13.0](https://github.com/kopexa-grc/common/compare/v1.12.0...v1.13.0) (2025-05-12)


### Features

* **tests:** add PostgreSQL test container utility ([#41](https://github.com/kopexa-grc/common/issues/41)) ([a85d99d](https://github.com/kopexa-grc/common/commit/a85d99d543c0e7fb5165658988bfc60c3b80967a))

## [1.12.0](https://github.com/kopexa-grc/common/compare/v1.11.0...v1.12.0) (2025-05-12)


### Features

* change default gravatar kind ([#39](https://github.com/kopexa-grc/common/issues/39)) ([beadddb](https://github.com/kopexa-grc/common/commit/beadddbbc2158e92326678df20a9702613bb31a5))

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
