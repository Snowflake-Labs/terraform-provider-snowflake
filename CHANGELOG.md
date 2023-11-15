# Changelog

## [0.76.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.75.0...v0.76.0) (2023-11-15)


### Features

* Add "CREATE DYNAMIC TABLE" to schema_grant ([#2144](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2144)) ([6f026f6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f026f64e6e24638df2b9d4110362836a9071011))
* Add views to the SDK ([#2171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2171)) ([ed079d3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed079d3d06dc3af083da04ca18314c8e7b07308e))
* Migrate application role to new sdk ([#2149](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2149)) ([7abb4db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7abb4dbb645f0eb0fcb1d39414b1ed0c322916c9))
* Migrate stage to the new sdk ([#2163](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2163)) ([1d08c46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d08c46c470cda2d6b7c7ed3d599d2c1a57b837b))
* Poc custom error type ([#2052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2052)) ([b86c4c3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b86c4c34d05f8b982fb6218a3a3a7500a23abf72))
* Use tasks from the SDK followup ([#2153](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2153)) ([82c3c13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82c3c13b6166168e470d7cb9b2982a8979275f17))


### Misc

* Add migration guide ([#2142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2142)) ([ee0f6af](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee0f6af54dbd269f8bfa2c3d73a396d98d10a6ab))
* Split existing alter operations ([#2156](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2156)) ([dbb7c91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dbb7c9136c586490a0856cc07ae879be491c8150))


### BugFixes

* Fix [#1947](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1947), [#2134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2134), [#2173](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2173), and [#2176](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2176) ([#2192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2192)) ([98d8ccc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/98d8ccc3c52fbd76d9d24c5fc7091ec0afa30e1a))
* provider muxing for plugin-framework ([#2130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2130)) ([f3c85c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3c85c0ebe3e54dc91836fe998c2d7c8e3373a52))

## [0.75.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.74.0...v0.75.0) (2023-10-26)


### Features

* add parse_header option to file format resource ([#2132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2132)) ([1e6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e6e54f828efa60edd258b316709fc4dfd370f93))
* Use streams from the new SDK in resource / datasource ([#2129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2129)) ([5c633be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5c633be461fd373d412b02b108e64b6cfc4eb856))
* Use task from SDK in resource and data source ([#2140](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2140)) ([de23f2b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de23f2ba939eb368d9734217e1bb2d4ebc75eef4))


### Misc

* Return multiple errors in existing validations ([#2122](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2122)) ([4d4bcdb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d4bcdbe841807da2fa08d534eaf846234934f7c))
* Set up a single warehouse for the SDK integration tests ([#2141](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2141)) ([16022ef](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16022ef4171e7dccf2932ae6e8d451b51c93291c))


### BugFixes

* cleanup acc tests ([#2135](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2135)) ([5db751d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5db751d1aa71952b1528e81cf2fdcd05d9d5d0fb))
* provider config ([#2136](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2136)) ([07b9b4f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/07b9b4fee800fe3f34890783cc463d4fc5904717))
* view statement update ([#2152](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2152)) ([6de32ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6de32ae6ec16ad76fb40afddfcaa7f650322cb67))

## [0.74.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.73.0...v0.74.0) (2023-10-18)


### Features

* dynamic table resource + data source ([#2104](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2104)) ([59c5de6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59c5de6becf4b16542da4db59e839f7ca5625e50))
* Migrate streams to the new sdk ([#2113](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2113)) ([521fde5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/521fde58b7836a21614d2ae31a7060352a68b465))
* Move integration tests to separate package ([#2111](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2111)) ([2755589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2755589bc8c7ecbda3115dc574339c40bfd5b096))
* Provider config refactor ([#2126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2126)) ([9858fea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9858feaa2fb382b57e3c89aae044f09ce2dcc5a3))
* tags sdk refactoring ([#2079](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2079)) ([7013f83](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7013f835864139ef40c2521e6af87589e909e201))
* Use network policy from sdk ([#2087](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2087)) ([50f2935](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/50f2935dc12fc2479c2919a5666d5573a8b737ff))
* use SDK in schema resource and datasource ([#2082](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2082)) ([f7d0d97](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f7d0d9716c4da9b00e5171a368a03a9cd52c31ae))


### Misc

* Adjust integration tests after moving to separate package ([#2115](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2115)) ([3f528a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f528a87f4c0b3bc95a0dfb35d93d22251b5112e))
* Integration tests setup db once ([#2118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2118)) ([f533368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5333684d5e46877ec4da6dc708ce2b09b5694af))
* Set up schema once for integration tests ([#2121](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2121)) ([5e3ebf4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e3ebf4e809a4be163a30642f69ea1b9146f21a9))
* Use helper methods in old unit tests ([#2119](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2119)) ([0c44571](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0c44571a92b842ac97d46a38e5bcfc1066367177))


### BugFixes

* use schema object identifier in external tables ([#2112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2112)) ([f5d4aeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5d4aebf810697aca96764261b891415f002ee92))

## [0.73.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.72.0...v0.73.0) (2023-10-11)


### Features

* add authenticator ([#2109](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2109)) ([4f3a551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4f3a5519484b0aab91ff5fa08f37a8cf512d1ec0))
* Add task clone to sdk ([#2105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2105)) ([acddb2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/acddb2cd6bfb1a7ffaf6dbb3c8349f7bc550c124))
* Add task to SDK ([#2099](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2099)) ([d52f334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d52f3347f091f0edff5e6daded1120542f1e9bd1))


### Misc

* **deps:** bump golang.org/x/crypto from 0.13.0 to 0.14.0 ([#2100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2100)) ([060750f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/060750f668b986d690619cd2d915a7a945350ec5))
* **deps:** bump golang.org/x/tools from 0.13.0 to 0.14.0 ([#2101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2101)) ([45a14d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/45a14d96d5fd6b051a73ced2fb03cab8659ded84))

## [0.72.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.71.0...v0.72.0) (2023-10-04)


### Features

* add dynamic tables to sdk ([#2074](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2074)) ([d1dfb05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d1dfb05fbb3bcc59cc2622b6b2d02ebadf1cf33f))
* Add grant ownership to SDK ([#2064](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2064)) ([f85ec8b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f85ec8baa2f9aaeead4f619dccfa3d38880a16d7))
* Add Manage Warehouses Account Grant ([#2017](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2017)) ([89c7148](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/89c7148c11378af9e42ea32bdad3e5a5c465d39c))
* add mfa auth ([#2077](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2077)) ([922358a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/922358a43c5383ee5840bf2971ecd27d96f86573))
* Add session policy to SDK ([#2088](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2088)) ([038241c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/038241c00a158b389df8034864a52a252fcc41bf))
* Go Generator + Network Policy migration to the new SDK ([#2061](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2061)) ([231b081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/231b081a34462392954c18b0f3452b1fed7d016f))
* Migrate external tables to new sdk ([#2006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2006)) ([5af17cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5af17cfbe4b6570c61b4ac77ab7564fd1f0529b7))
* Migrate roles to the new sdk ([#2007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2007)) ([5b996e6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b996e6f4b3ea28bae0cd13e75bd112d9af09c68))
* support drop for accounts and fix read race condition ([#2065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2065)) ([5412252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/541225258ff19de864c8d25d05a8c1b4a7941813))


### Misc

* Change return type of show ([#2045](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2045)) ([21f069a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/21f069a3232a0039792f30445e25b5e88f09ab25))
* Cleanup comments ([#2092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2092)) ([3a06a66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a06a66f27d9b03d037affdda5b7dbb3bd81d5d7))
* **deps:** bump github.com/hashicorp/terraform-plugin-docs from 0.14.1 to 0.16.0 ([#1931](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1931)) ([07c4f14](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/07c4f149bf08ea9fecd28aa12a47a9fae3e0523b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#2090](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2090)) ([24b6313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/24b6313a7f28731ad4c10dd80c32e7b9cf397c5d))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1842)) ([9b0825d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9b0825db520c244cea962830d2139fb50186f23c))
* **deps:** bump golang.org/x/tools from 0.7.0 to 0.13.0 ([#2089](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2089)) ([0ace968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ace9689c6b29382b1d93d9a3542eacc0d7c4270))
* Make options naming consistent ([#2083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2083)) ([df84bd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/df84bd084482ef3f85db2b1f4a3517afe8df4ddb))


### BugFixes

* added privilege for `RESOLVE ALL` ([#1861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1861)) ([18cf7b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18cf7b09c9be927b5e1a625acdadb90abd048ff8))
* database grant read ([#2063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2063)) ([d93ddd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d93ddd564cd5cc76b47d49029f99c503d9314f29))
* Fix database role docs ([#2081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2081)) ([339b1ff](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/339b1ff4fecca3dfe815231c9d3f26dc43a1ed7f))
* Ignore two flaky table tests and disable comment unset for password policy ([#2086](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2086)) ([ee90014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee90014cc2d0dc5088d356f6c9b0d094c34d2a46))

## [0.71.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.70.1...v0.71.0) (2023-09-21)


### Features

* Add grant privileges to database role to SDK ([#2023](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2023)) ([717289f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/717289f71fd4a08f44d4f20f6e16dc4dded77803))
* Convert alert datasource to new sdk ([#2020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2020)) ([2d0eaeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d0eaebe8ea1605251dd24f86a2f83e6e3be67c8))
* Initialize SDK generator ([#2033](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2033)) ([96b47e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96b47e5fdbb6175a7a24d0518b975279927c22de))
* Parameters sdk v2 ([#1914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1914)) ([1d15355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1535572e1b82f3cb2308e3c19e632eb34fe24b))


### BugFixes

* add dynamic tables to grants ([#2059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2059)) ([3767de9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3767de97ec685acebd0e8e105878eb9129a0752c))
* SNOW-59564 remove default data retention in days value ([#2029](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2029)) ([53c20b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53c20b4dc3b1141bcad07976f2dba93198e4d200))

## [0.70.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.70.0...v0.70.1) (2023-09-01)


### Misc

* Extract reviewdog job from integration workflow ([#2027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2027)) ([243fc28](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/243fc285811de787328a7f3ca0d3a84193976b85))


### BugFixes

* add back in 386 build ([#2038](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2038)) ([5522504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55225041d7f4722adbd7b9445990f57c04a4adea))
* external id for shares ([#2040](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2040)) ([f8357c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8357c06376feda2d9732f5046a408f9b82bc339))
* Fix make sweep ([#2025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2025)) ([beb2f98](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/beb2f9867ee501e630f75ace42dbef577563d593))

## [0.70.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.69.0...v0.70.0) (2023-08-21)


### Features

* Add database role to SDK ([#2009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2009)) ([f5efc09](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5efc09ea60bd2d66c65c9e07cb84321f95531f0))
* Add missing database role operations to SDK ([#2014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2014)) ([d2ea67d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2ea67d7fee00b15e1222fe37efe8e7a1cddecb5))
* Add Pipes to SDK ([#1968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1968)) ([69a543f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69a543fd729b64bdd8964dc34626dee83b3f96a7))
* Generate constructor and builder methods for DTOs ([#2001](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2001)) ([79d9c9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79d9c9ce77e85dc9bd0889da30155cef0a74c293))
* Schemas sdk v2 ([#1975](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1975)) ([289ad8a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/289ad8a8d61cf217d4548a65cf4843e416baa1da))
* Use database role from SDK ([#2012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2012)) ([294075a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/294075a4d145dc08071c11f067afaecb78fe8ef7))
* Use pipes from SDK ([#2003](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2003)) ([079d47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/079d47d22af22edba6c6c499409b264c4c5f5945))
* Users sdk v2 ([#1945](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1945)) ([d644b63](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d644b63e673223b67f7092563a0bfc8c826fc38b))


### Misc

* Makefile cleanup ([#1995](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1995)) ([5c6fdbe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5c6fdbe48aadab9d56bfb7cf5d5ea30ca5a2b356))
* Remove lint ignore comments - test ([#2011](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2011)) ([8f81a42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f81a42544a668b03e0bb5d4d3199e6879d2f937))
* Rename and reuse validatableOpts interface ([#2016](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2016)) ([4b42848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b428480a3f8fb9aa38da7a311f8a04f8d294882))


### BugFixes

* unable to alter stage if url and si have changes  ([#1982](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1982)) ([#1983](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1983)) ([3813aaa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3813aaac52a3ec5186e285066e4db72c2046531a))

## [0.69.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.68.2...v0.69.0) (2023-07-29)


### ⚠ BREAKING CHANGES

* mark snowflake_user.login_name as non sensitive

### Features

* mark snowflake_user.login_name as non sensitive ([71c0a9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/71c0a9c2897e07b2251680d619cbbcaeba5779b2))
* resource monitor sdk v2 ([#1892](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1892)) ([707d723](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/707d7233bd867d60e13c8785d1c0adef558a5853))


### BugFixes

* 1969  change default privilege to OPERATE ([#1970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1970)) ([e8721f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e8721f07dd27ce7a903c8c4d2b19fee77ce72cdc))
* add external table grant to snowflake_grant_privileges_to_role resource ([#1967](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1967)) ([8c84c4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c84c4a5f6eee7cf35b5ec2d77fac1c7fcee05df))
* future ownership task grant ([#1954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1954)) ([#1955](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1955)) ([81ac1a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/81ac1a396bba713b0a5ec808d72b880d59e17d47))

## [0.68.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.68.1...v0.68.2) (2023-07-17)


### BugFixes

* support multiple grant_privilege_to_role resources on the same role ([#1953](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1953)) ([dfdd166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dfdd166e20347deb6f464a5bbbe9c345287fcd6b))
* supress warehouse size changes ([#1889](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1889) ([daf2cc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daf2cc289bab892e5a0954f852a72a9f99e1b0c1))
  
## [0.68.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.68.0...v0.68.1) (2023-07-07)


### BugFixes

* reading file formats with no default schema ([#1936](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1936)) ([c5602f5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5602f5f01e09fc90a82aa2fbb8e116648f04c58))

## [0.68.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.67.0...v0.68.0) (2023-07-06)


### Features

* alerts sdk v2 ([#1871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1871)) ([fcf0380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fcf03800cf92e1acff7f89e59968914f1bcad4f4))
* new grant resource "snowflake_grant_privileges_to_role"  ([#1929](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1929)) ([e241c22](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e241c223b0c5d159efc803b9fcd4671c8325ae00))
* set query tag for terraform sessions  ([#1826](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1826)) ([6629583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6629583bc288e5c23cf02d8a53597adbadf55fce))


### BugFixes

* add missing option ([#1924](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1924)) ([2991a16](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2991a16fe45491d77a8212a3ce0aa9566882de4f))
* do not set "query_acceleration_max_scale_factor" if Query Acceleration is not enabled on Warehouse creation ([#1866](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1866)) ([7679e0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7679e0b0a480c6c6867f6e03051e7c3f97709edf))
* file format now supports reading options ([#1917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1917)) ([2596990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2596990808d97a7b98fa459a74eb7a77edf25cd5))
* prevent view drop when copy_grants attribute changes ([#1920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1920)) ([0cf22d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0cf22d04513e70b0c945cd402b41b964f4d14be0))

## [0.67.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.66.2...v0.67.0) (2023-06-22)


### Features

* account password policy attachment ([#1824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1824)) ([f408828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f408828fd023c2207ec41f702cec7dae524b1e93))


### BugFixes

* Generate security integration statements using double quoted names. ([#1897](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1897)) ([a21d44f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a21d44fc9109c589a7e4c78032e2f5a62b6ed514))
* missing 'kind' field ([#1872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1872)) ([59b7ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59b7ef5d791ecce366063a5cb8e2d68c68fb3e39))
* support empty comments in oauth integrations ([#1900](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1900)) ([53d46b3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53d46b3a97680a4bf0e402d48c7db15c0d3c6f03))

## [0.66.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.66.1...v0.66.2) (2023-06-13)


### BugFixes

* warehouse size validation ([#1873](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1873)) ([5bbe460](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5bbe4600e2768a66a83df67a16368888b7f1d76b))

## [0.66.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.66.0...v0.66.1) (2023-06-03)


### BugFixes

* err when reading profile ([#1853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1853)) ([29c4633](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/29c463383da5fbb8de0af5c957aa9f3a34572c6e))

## [0.66.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.65.0...v0.66.0) (2023-06-02)


### Features

* databases sdk v2 ([#1840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1840)) ([410b493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/410b493ff75188b70cee91db56c60538d6823aac))


### Misc

* rename "db" symbol to "sql" and make input options more uniform in naming convention ([#1837](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1837)) ([244b5b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/244b5b45e22f66875f14be9bcb9af08ad41d37ca))

## [0.65.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.64.0...v0.65.0) (2023-05-30)


### Features

* add failover groups to sdk v2, and add data source for failover groups ([#1825](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1825)) ([44e8c06](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44e8c06ba4c665c81f0b909dbd3df90c4925e179))
* add shares to sdk v2 ([#1813](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1813)) ([a814841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a814841011f08857d8d37691fa5ff01cd9412176))
* adding a subset of missing parameters ([#1828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1828)) ([6d1a572](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d1a572623438dd96139edbd698ce5974fa3df61))
* adding CREATE STREAMLIT as privalege ([#1771](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1771)) ([7b112cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b112cc94822807d3c25c78ec812ec2c2c66ea16))
* alerts datasource ([#1727](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1727)) ([a1d0d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1d0d2eeb9a8a03d4c347b9f153182e66c237a4f))
* allow all priv on grants ([#1786](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1786)) ([aa3f873](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aa3f873b4669c5942c10681151d6afd4ea03623f))
* support custom ownership transfer for grants ([#1743](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1743)) ([eaa6e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eaa6e01820cb04ffa3c647ecd893e479af8e35a1))


### Misc

* Account SDK ([#1822](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1822)) ([4c22b89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c22b89e82a608375fd78993134cf3a0b45b985b))
* Migrate Warehouse resource + datasource to new SDK ([#1792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1792)) ([a14b994](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a14b994d62f73c272cde72651c3c8c18ac4213cf))


### BugFixes

* 0.65 integration tests and minor linting fixes ([#1835](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1835)) ([8d7d663](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8d7d6636a67d272d5b87e02a178ae799b5b6d777))
* add grant to sdk v2, better validation ([#1812](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1812)) ([76da3b8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/76da3b84d9d65493ae666fad6ea6d926b81add8f))
* Add missing Google Cloud computed property api_gcp_service_account to the api integration resource ([#1776](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1776)) ([1dd0672](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dd06724aa1add9b504c820583a2a752555c9bdb))
* allow password_policy.max_age_days = 0 ([#1819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1819)) ([eeaadfe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eeaadfef5b8fc31f84304899c7b205c9a97d8852))
* fixed enable_multiple_grants behaviour for role_grants ([#1816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1816)) ([f508129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f50812910a679bdea2d3add65f92d7cebeaa345c))
* gorelease freebsd ([#1790](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1790)) ([c87f74f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c87f74f0cdffba16108ee822f845d0e961205252))
* reading warehouse with underscore in name ([#1793](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1793)) ([5e184fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e184fab4cbabcf27f290eb79c01a0cd6cc79282))
* warehouse for sdk v2 ([#1804](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1804)) ([99f7621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/99f7621cad1aafa83f96a4ed52d0fd04d6dfb2ae))

## [0.64.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.63.0...v0.64.0) (2023-05-09)


### Features

* Add external functions translators ([#1735](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1735)) ([1f67286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1f672862adb29a658f5e81e940f9afb994347f2f))
* addition of the on_all option for all remaining resources ([#1742](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1742)) ([e3a7710](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3a771000fa527fc8420b892450031ba39394fcf))
* database roles datasource ([#1726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1726)) ([f1cb7f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1cb7f3b368c0a7fa6a0a08784d53b0e92efda81))
* email notification integration ([#1720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1720)) ([5d21fe1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d21fe19f0da9ed4939e9bd6dc66237685c5a327))
* masking policy in v2 sdk ([#1777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1777)) ([6978c42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6978c42b304bd8c5429c88bcd7d6ed20ac3fd98a))


### Misc

* sdk v2 implemented for password policy ([#1752](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1752)) ([0cb1164](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0cb1164414607ce2e4ac7f6fb2da1c563b22da7b))


### BugFixes

* add missing reference for email_notification ([#1768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1768)) ([6feeb9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6feeb9c09f56df0fb222bd0664508defb7c5af6e))
* database from_share handles organization ([#1711](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1711)) ([9e71eb6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e71eb6b4b69ed7fe1d42febe30fb166ea00812a))
* get `on_view` from `table_name` for `stream` ([#1740](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1740)) ([21fc2b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/21fc2b9ad9ddf02a241fbf2d1c3297d98cfa26c9))

## [0.63.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.62.0...v0.63.0) (2023-04-25)


### Features

* masking policy application resource ([#1739](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1739)) ([ce80f57](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce80f57225d6019765a50da6eb74619e401db622))


### Misc

* fix parameters datasource acceptance test ([#1749](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1749)) ([338a19d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/338a19de373525087372cc319aa065f2cddbb724))

## [0.62.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.61.0...v0.62.0) (2023-04-19)


### Features

* Add gcp_pubsub_topic_name parameter for gcp notification integration ([#1687](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1687)) ([a30d0cb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a30d0cb756a2281a4d880af9e32651c04409028e))
* added on_all grants for view, stage, schema and materialized view resource ([#1686](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1686)) ([f27a9e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f27a9e454102709e1134d796f1594b4932670ae8))
* alert resource ([#1653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1653)) ([4a84eb9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a84eb91af6201330a2ec60de306db5ed3b5ab38))
* data source shares ([#1651](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1651)) ([5dd447f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dd447fc79c153f025e420052a23b90751cdaa9f))
* database role resource ([#1654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1654)) ([9f4a516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f4a5165d6b7e1fa63ad06450ffc55689be66356))
* external_oauth_integration.scope_mapping_attribute ([#1722](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1722)) ([66d88bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/66d88bd29b6aa78dbc2acd11cf729fa9380b011d))
* Password policy resource ([#1702](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1702)) ([7ee293b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ee293b52262170a1d53811f9beb6fc3772ae913))
* provider config support insecureMode ([#1703](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1703)) ([e269925](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e269925e0f54ca9bbd742e5c60eaebad5f9b1f04))


### BugFixes

* Adding double quotes to the column names ([#1729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1729)) ([791dd0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/791dd0b3a613d0ada95d4af7153a7b0a2bbee219))
* be more cautious when setting db id to empty ([#1725](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1725)) ([e78e0c8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e78e0c816270779a78d839f78debc519fd53f5bf))
* change incorrect select privilege to usage ([#1681](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1681)) ([3d8f5b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d8f5b5ad417d2ebcf171024fc3217440267586b))
* don't ignore MUST_CHANGE_PASSWORD on account creation ([#1699](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1699)) ([7fcacfa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fcacfa59995766fb8c474b2984b9c9d3092d9ad))
* fmt / linting fix ([#1694](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1694)) ([d7d910e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d910edc95069f6f01bc350c3319eb36f75c127))
* grant imports ([#1713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1713)) ([19a4156](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19a415654f198a68e823b9ef62326992feba64df))
* resource monitor import test ([#1690](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1690)) ([a4e58bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a4e58bfa0d5d4f90668784980fea7e2cc443f25d))
* resource_monitor suspend triggers conflict ([#1682](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1682)) ([14d1200](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14d120008d6d987dd2c0773ad9815eaac17a2c3c))

## [0.61.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.61.0) (2023-04-03)


### Features

* Add COPY GRANTS arg to views resources. ([#1668](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1668)) ([7225d93](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7225d93ef3e50c6810c0dd57cfd7079e882d443f))
* add on_account to session and object params ([#1685](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1685)) ([1329430](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/13294304c7626c9d428682986669d2e97ab2c23b))
* mark PII as sensitive in account and user resources. ([#1678](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1678)) ([caa461f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/caa461fea56df4f33f73e586b537802c01f1eb5d))
* support is_org_admin column in SHOW ORGANIZATION ACCOUNTS ([#1673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1673)) ([263c521](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/263c5215dedf70195c0c2c7d8e1505e4b9c0828c))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* wildcards in database name ([#1666](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1666)) ([54bf74c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/54bf74ca3d0119d31668d18bd1599ed029e386c8))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
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
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.0 ([#1665](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1665)) ([ea23020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea23020801ea4d43f055f2b443400385d96a135b))
* **main:** release 0.60.0 ([#1667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1667)) ([9d3e40f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d3e40fab05023aff16795266ec8a30761560c26))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
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
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.0 ([#1665](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1665)) ([ea23020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea23020801ea4d43f055f2b443400385d96a135b))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* wildcards in database name ([#1666](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1666)) ([54bf74c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/54bf74ca3d0119d31668d18bd1599ed029e386c8))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
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
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.61.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
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
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

## [0.61.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.1...v0.61.0) (2023-03-27)


### Features

* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))


### BugFixes

* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))

## [0.60.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.60.1) (2023-03-23)


### BugFixes

* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))

## [0.60.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.59.0...v0.60.0) (2023-03-23)


### Features

* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))


### BugFixes

* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))


### Misc

* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))

## [0.59.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.58.2...v0.59.0) (2023-03-21)


### Features

* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))


### BugFixes

* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))

## [0.58.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.58.1...v0.58.2) (2023-03-16)


### BugFixes

* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))

## [0.58.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.58.0...v0.58.1) (2023-03-16)


### BugFixes

* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))

## [0.58.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.57.0...v0.58.0) (2023-03-03)


### Features

* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))


### BugFixes

* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))

## [0.57.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.5...v0.57.0) (2023-02-28)


### Features

* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))


### BugFixes

* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))

## [0.56.5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.4...v0.56.5) (2023-02-21)


### Misc

* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))

## [0.56.4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.3...v0.56.4) (2023-02-17)


### BugFixes

* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))

## [0.56.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.2...v0.56.3) (2023-02-02)


### BugFixes

* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))

## [0.56.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.1...v0.56.2) (2023-02-01)


### BugFixes

* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))

## [0.56.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.0...v0.56.1) (2023-01-31)


### BugFixes

* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))

## [0.56.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.55.1...v0.56.0) (2023-01-27)


### Features

* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))


### BugFixes

* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))

## [0.55.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.55.0...v0.55.1) (2023-01-11)


### BugFixes

* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))

## [0.55.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.54.0...v0.55.0) (2023-01-10)


### Features

* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))


### Misc

* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))


### BugFixes

* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.54.0...v0.34.0) (2022-12-23)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* Resource to manage a user's public keys ([#540](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/540)) ([590c22e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/590c22ec40ed28c7d280192ed66c4d93623e32fd))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
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
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
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
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.54.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.53.1...v0.54.0) (2022-12-23)


### Features

* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))


### Misc

* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))


### BugFixes

* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))

## [0.53.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.53.0...v0.53.1) (2022-12-08)


### Misc

* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))


### BugFixes

* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))

## [0.53.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.52.0...v0.53.0) (2022-12-07)


### Features

* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))


### BugFixes

* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.52.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.51.0...v0.52.0) (2022-11-17)


### Features

* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))


### Misc

* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))


### BugFixes

* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))

## [0.51.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.50.0...v0.51.0) (2022-11-07)


### Features

* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))


### Misc

* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))


### BugFixes

* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))

## [0.50.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.0...v0.50.0) (2022-11-04)


### Features

* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))

## [0.49.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.2...v0.49.3) (2022-11-01)


### BugFixes

* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))

## [0.49.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.1...v0.49.2) (2022-11-01)


### BugFixes

* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))

## [0.49.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.0...v0.49.1) (2022-10-31)


### BugFixes

* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))

## [0.49.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.48.0...v0.49.0) (2022-10-31)


### Features

* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))


### Misc

* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))


### BugFixes

* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))

## [0.48.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.47.0...v0.48.0) (2022-10-24)


### Features

* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))


### BugFixes

* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))

## [0.47.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.46.0...v0.47.0) (2022-10-11)


### Features

* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))


### BugFixes

* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))

## [0.46.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.45.0...v0.46.0) (2022-09-29)


### Features

* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))


### BugFixes

* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))

## [0.45.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.44.0...v0.45.0) (2022-09-22)


### Features

* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))

## [0.44.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.43.1...v0.44.0) (2022-09-20)


### Features

* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))


### BugFixes

* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))

## [0.43.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.43.0...v0.43.1) (2022-09-20)


### BugFixes

* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))

## [0.43.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.42.1...v0.43.0) (2022-08-31)


### Features

* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))


### BugFixes

* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))

## [0.42.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.42.0...v0.42.1) (2022-08-24)


### Misc

* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))

## [0.42.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.41.0...v0.42.0) (2022-08-24)


### Features

* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))

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
