# Changelog

## [0.41.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.40.0...v0.41.0) (2022-08-10)


### Features

* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))


### BugFixes

* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))

## [0.40.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.39.0...v0.40.0) (2022-07-14)


### Features

* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))


### BugFixes

* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))

## [0.39.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.38.0...v0.39.0) (2022-07-13)


### Features

* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))


### BugFixes

* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))

## [0.37.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.37.0...v0.37.1) (2022-07-01)


### BugFixes

* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))

## [0.37.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.36.0...v0.37.0) (2022-06-28)


### Features

* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))


### BugFixes

* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))

## [0.36.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.35.0...v0.36.0) (2022-06-16)


### Features

* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))


### BugFixes

* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))

## [0.35.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.34.0...v0.35.0) (2022-06-07)


### Features

* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))


### BugFixes

* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))


### Misc

* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.34.0...v0.34.0) (2022-05-25)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* Resource to manage a user's public keys ([#540](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/540)) ([590c22e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/590c22ec40ed28c7d280192ed66c4d93623e32fd))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))


### BugFixes

* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* build: Add trimpath Go flag to build ([#316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/316)) ([420a84b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/420a84b61cf342e8210f440ccfaca5dcaa589ede))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* OSCP -> OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* table: Properly read and import table state ([#314](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/314)) ([df1f6bc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/df1f6bcabfca27058c860a7db815d263457afd6c))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))


### Misc

* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

### [0.33.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.33.0...v0.33.1) (2022-05-03)


### BugFixes

* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))

## [0.33.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.32.0...v0.33.0) (2022-04-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))

## [0.32.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.31.0...v0.32.0) (2022-04-14)


### Features

* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))

## [0.31.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.30.0...v0.31.0) (2022-04-11)


### Features

* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))

## [0.30.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.29.0...v0.30.0) (2022-03-29)


### Features

* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))

## [0.29.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.8...v0.29.0) (2022-03-23)


### Features

* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))


### BugFixes

* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))

### [0.28.8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.7...v0.28.8) (2022-03-18)


### BugFixes

* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))

### [0.28.7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.6...v0.28.7) (2022-03-15)


### BugFixes

* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))

### [0.28.6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.5...v0.28.6) (2022-03-14)


### BugFixes

* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))

### [0.28.5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.4...v0.28.5) (2022-03-12)


### BugFixes

* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))

### [0.28.4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.3...v0.28.4) (2022-03-11)


### BugFixes

* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))

### [0.28.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.2...v0.28.3) (2022-03-10)


### BugFixes

* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))

### [0.28.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.1...v0.28.2) (2022-03-09)


### BugFixes

* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))

### [0.28.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.0...v0.28.1) (2022-03-09)


### BugFixes

* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))

## [0.28.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.27.0...v0.28.0) (2022-03-05)


### Features

* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))


### BugFixes

* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))

## [0.27.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.3...v0.27.0) (2022-03-02)


### Features

* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))


### Misc

* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

### [0.26.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.2...v0.26.3) (2022-02-08)


### BugFixes

* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))

### [0.26.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.1...v0.26.2) (2022-02-07)


### BugFixes

* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))

### [0.26.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.0...v0.26.1) (2022-02-07)


### BugFixes

* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))

## [0.26.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.25.36...v0.26.0) (2022-02-03)


### Features

* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))


### Misc

* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))


### BugFixes

* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
