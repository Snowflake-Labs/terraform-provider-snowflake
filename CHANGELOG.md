# Changelog

## [0.31.0](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.30.0...v0.31.0) (2022-04-11)


### Features

* Add manage support cases account grants ([#961](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* snowflake_user_ownership_grant resource ([#969](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))

## [0.30.0](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.29.0...v0.30.0) (2022-03-29)


### Features

* support host option to pass down to driver ([#939](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))

## [0.29.0](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.8...v0.29.0) (2022-03-23)


### Features

* Allow in-place renaming of Snowflake tables ([#904](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))


### BugFixes

* Legacy role grantID to work with new grant functionality ([#941](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))

### [0.28.8](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.7...v0.28.8) (2022-03-18)


### BugFixes

* change the function_grant documentation example privilege to usage ([#901](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* remove share feature from stage because it isn't supported ([#918](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))

### [0.28.7](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.6...v0.28.7) (2022-03-15)


### BugFixes

* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* sequence import ([#775](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))

### [0.28.6](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.5...v0.28.6) (2022-03-14)


### BugFixes

* Add release step in goreleaser ([#919](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))

### [0.28.5](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.4...v0.28.5) (2022-03-12)


### BugFixes

* Add manifest json ([#914](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))

### [0.28.4](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.3...v0.28.4) (2022-03-11)


### BugFixes

* Add gpg signing to goreleaser ([#911](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))

### [0.28.3](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.2...v0.28.3) (2022-03-10)


### BugFixes

* issue with ie-proxy ([#903](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))

### [0.28.2](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.1...v0.28.2) (2022-03-09)


### BugFixes

* Ran make deps to fix dependencies ([#899](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))

### [0.28.1](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.28.0...v0.28.1) (2022-03-09)


### BugFixes

* Release by updating go dependencies ([#894](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))

## [0.28.0](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.27.0...v0.28.0) (2022-03-05)


### Features

* Implemented External OAuth Security Integration Resource ([#879](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))


### BugFixes

* escape string escape_unenclosed_field ([#877](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))

## [0.27.0](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.26.3...v0.27.0) (2022-03-02)


### Features

* Data source for list databases ([#861](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* Support DIRECTORY option on stage create ([#872](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))


### Misc

* Upgarde all dependencies to latest ([#878](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

### [0.26.3](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.26.2...v0.26.3) (2022-02-08)


### BugFixes

* Remove keybase since moving to github actions ([#852](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))

### [0.26.2](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.26.1...v0.26.2) (2022-02-07)


### BugFixes

* Update goreleaser env Dirty to false ([#850](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))

### [0.26.1](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.26.0...v0.26.1) (2022-02-07)


### BugFixes

* Release tag ([#848](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))

## [0.26.0](https://github.com/chanzuckerberg/terraform-provider-snowflake/compare/v0.25.36...v0.26.0) (2022-02-03)


### Features

* Add replication support ([#832](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Release GH workflow ([#840](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* TitleLinter customized ([#842](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))


### Misc

* Move titlelinter workflow ([#843](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))


### BugFixes

* Allow multiple resources of the same object grant ([#824](https://github.com/chanzuckerberg/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/chanzuckerberg/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
