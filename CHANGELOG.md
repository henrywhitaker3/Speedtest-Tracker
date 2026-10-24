## [2.3.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.2.4..v2.3.0) - 2026-10-24

### 🚀 Features

- Move to workflow based speedtest runner - ([dbdc013](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/dbdc013d790ef8fcb2b50c2eb1105ddb73023f92))
- Display timeline of the run workflow in the results details page - ([38f7ed8](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/38f7ed8940a320ee613fc96fb0eb57ac7f50182b))

### 🐛 Bug Fixes

- *(ui)* Remember results table page - ([59bd3f8](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/59bd3f8bd7e4c20f84c39351166b13ab81f5ac02))
- *(ui)* Poll for fresh result data when it is not complete yet - ([1f4e231](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1f4e231b6f976ebaae9e0f457ec6e3d5cc79b0bf))
- *(ui)* Nicer view of steps in the timeline - ([4e7d156](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4e7d15644e262ea70fcee441b7d9cc95f1cbb266))
- *(workflow)* Don't block running when healthchecks fail - ([7ec9ec6](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7ec9ec68805e7fb6bff245d7ab1e39a77483f556))
- *(workflow)* Error logging in failed run - ([bc898c4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/bc898c422a341be763dc7db04750531e863e3760))
- Close healthchecks response bodies - ([e74bbfe](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e74bbfee39eeac1e45da36bf43f2cf0665e86edf))
- Cleanup dangling result when workflow fails to start - ([82a7b1b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/82a7b1bf2db20b402759b9263f4540939452dcee))

### 🧪 Testing

- Add workflow testing - ([772176f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/772176fd8c1a34ec1aef5f6e2810aa244437ed2a))
- Update tests to ensure workflow continues on healthcheck failure - ([c29d7e9](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/c29d7e9e07c93610ae1eaf0b2105643cee4c27eb))

### ⚙️ Miscellaneous Tasks

- *(test)* Update step name - ([2f9190f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/2f9190f68f344a167a581a8b98cd4cbc5d18dde6))
- *(test)* Don't persists git creds - ([47e1e1e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/47e1e1ef8b5e10d6aab6cb679fdd7d9e9b44c4e7))
- Add lint job - ([fefe7f1](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fefe7f11ed162ebd615710ef11585233d241c361))
- Add ui build job too - ([a300ba8](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a300ba8d9012faf252a465c917e627d5aae96705))
- Add ui build job too - ([d9c796b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d9c796b1e1d30164a024cca0bfe1742ef58d67bc))
- Add lint job - ([dc07830](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/dc0783042ce408c47a8864680a96fbed74d870a1))
## [2.2.4](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.2.3..v2.2.4) - 2026-10-23

### 🐛 Bug Fixes

- *(api)* Log requests at debug level - ([8d31fe9](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/8d31fe9ea6c7914a0b10a885d50259e2fb9f9613))
- *(api)* Add handler to delete individual test results - ([5fc7a0e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5fc7a0e9e45d19b61eef5d4730aa13ef16c7a1f1))
- *(auth)* Admin settings broken when auth disabled - ([cc6fdaa](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/cc6fdaa2b1b0384ba9a793a31c1f345d66d45d5c))
- *(ui)* Add link to changelog - ([5dd0be4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5dd0be44dfc664c9bfc180bfde07b0305b77eaa8))

### 🚜 Refactor

- Move repeated service name to constant - ([3b11de3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/3b11de39aac0dfa0a2cf5d3e8a1a40bc9c37b821))

### ⚡ Performance

- *(auth)* Cache oidc config and sigining keys - ([7e658c5](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7e658c52a8b990c00aa4091c7ee9bc40c6b2ab81))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([d688d64](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d688d6447443359acbcd8ff754b52d42ab0e7309))
## [2.2.3](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.2.2..v2.2.3) - 2026-10-23

### 🐛 Bug Fixes

- *(ui)* Nicer date formatting - ([9c2ec5f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/9c2ec5f7625fcd3fd7810f4db7f9d1d7befa58cf))
- Add endpoint to purge results with errors - ([9830262](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/983026284f386dfd545016c882218ddf99b034a8))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([689b22c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/689b22cdd68db5a53e290392223cfd79509659da))
- Update git cliff version - ([f57d129](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/f57d129e6e43c655153f8145252e1812c37934af))
## [2.2.2](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.2.1..v2.2.2) - 2026-10-23

### 🐛 Bug Fixes

- *(auth)* Keep sessions for 7 days - ([fa0e6c6](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fa0e6c688986b29146b72a1b4ff161f8190bae1e))
- *(retention)* Optionally remove results older than x duration - ([b753a0f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b753a0f3bffcc5f191740ef45d3f84f10b7e68ab))
- *(ui)* Add dropdown for stats time range - ([796c641](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/796c641be1f37721f7146659546e58f6c7758758))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([ebd846d](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ebd846d2103a0ea2d4006efd027bdf9386b2cde0))
## [2.2.1](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.2.0..v2.2.1) - 2026-10-23

### 🐛 Bug Fixes

- *(ui)* Text box showing created token light mode - ([4575f45](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4575f453ef958067a27ac381c016fe4b57107c5f))
- *(ui)* Hide points from the graph - ([1ea6b57](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1ea6b57d68f18f5724e64f2afcf195251c242543))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([d31fada](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d31fadacbf5aff4fe8e100ec3b415f485adf1ca6))
## [2.2.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.1.4..v2.2.0) - 2026-10-23

### 🚀 Features

- *(auth)* Add http handlers for oidc auth - ([6406773](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6406773198cf332028aaad372efbbc22309f7c3f))
- *(auth)* Added oidc auth session actor - ([4e82768](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4e827689578716062f18d5d22423bbba39552008))
- *(auth)* Add backend handling for api tokens - ([cde17a4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/cde17a47514b76757f497af6ffe69a004e3803fa))
- *(auth)* Added auth token middleware support - ([92a3ab4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/92a3ab49e6388c2c8a54e91d55330d4daa3613e6))
- *(ui)* Gate speedtest run behind auth - ([29637fd](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/29637fd4bacec95439adc8b8a9e11a4fa4e689d4))
- *(ui)* Add crud for api tokens - ([e52c023](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e52c0238f377605f41a33959ab44c31e4bf9724d))
- Add login handler - ([b1737c2](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b1737c290536f16b3c1c75b7722fdd6f4c204385))

### 🐛 Bug Fixes

- *(auth)* Add dependencies for oidc - ([68eff68](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/68eff687e79228546d180e53368bca5ebfd70fa7))
- *(auth)* Add authz for admin-only endpoints - ([4f8b56d](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4f8b56d98643bcc36ccf5423e087675fd3458bae))
- *(auth)* Update group data when a token is refreshed - ([7bb9823](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7bb9823fc51164531fad5ba4a8325212f84758dd))
- Add metrics for ping/up/down - ([5ebd92b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5ebd92bb0be53f4c704320b6fe380f1e2c2e13dd))
- Add config endpoint - ([0a7da2d](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0a7da2d46faeb07c9364e9634794d73018eaeee9))
- Add token/:id paths to auth - ([71e1c8f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/71e1c8f4eb02009b5cb38ca9373e5454aab1e6e4))

### 🧪 Testing

- Add checks for auth in run endpoint - ([d492846](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d4928462f174b1ccb22da7a2838edac2ac1a80b0))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update readme with auth info - ([b2e0b82](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b2e0b828ede8964c52a7100bd280b7ccabb83c3e))
- *(docs)* Update changelog - ([6fba9d0](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6fba9d0eb71734cacca44f59fbc4b8623aba245c))
## [2.1.4](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.1.3..v2.1.4) - 2026-10-23

### 🐛 Bug Fixes

- Pass server id through to actor - ([be80943](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/be809439639fd2b79b681de9a1d530bddf76a496))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([522fe46](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/522fe461d3db49620e114214cc5dc598c3927f43))
## [2.1.3](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.1.2..v2.1.3) - 2026-10-23

### 🐛 Bug Fixes

- Disable navbar overflow - ([3a125fb](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/3a125fbe7482a0071942478781d79b9f3f65412c))
- Make latest cards responsive - ([0a51d51](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0a51d518511b475506fed16613f6fd8c6120d741))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([c2e15bf](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/c2e15bf92ade10b1d68872da1ba149922f9dada6))
- Update README.md - ([e14f975](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e14f975f0b1efbdaeede13edb9013fa2e915a644))
- Update readme.md - ([8ba21a1](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/8ba21a10237f7bccc47c8c405973f122f7ae047c))
## [2.1.2](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.1.1..v2.1.2) - 2026-10-23

### 🐛 Bug Fixes

- Ready probes never firing - ([60faafc](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/60faafce8630d768a0b1084a9b58254c0cc656d3))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([f0f0bc0](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/f0f0bc0605e764253caace5f77710d63c2b7d0e3))
- Remove dangling package-lock.json - ([8e99fbc](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/8e99fbc4ca99870a79a8db6c5e186035debe5315))
## [2.1.1](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.1.0..v2.1.1) - 2026-10-23

### 🐛 Bug Fixes

- Issue with fallback 404s for direct to spa routes - ([176cc07](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/176cc07f40aa1b227bd882c6dd129866e9f91044))
- Order timeseries results - ([b6b0ff2](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b6b0ff216ae3d31a49ec457aa8eea8490e71484d))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([5dedd95](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5dedd9590757aad247eb6a5d919e68f0ff3bb26f))
## [2.1.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.0.3..v2.1.0) - 2026-10-23

### 🚀 Features

- Backend support for rate limiting speedtest runs - ([fa33818](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fa33818f94be49095e3603226197695973073aec))
- Add page to show single result details - ([7e5801b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7e5801b5bfa86a06c11ebb769e2ca4a68b75c57c))

### 🐛 Bug Fixes

- *(deps)* Update module github.com/labstack/echo/v5 to v5.4.2 - ([8350704](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/8350704c7cb296f2140ca4b9981c33ccb9ddeace))
- *(deps)* Update module gorm.io/driver/postgres to v1.6.3 - ([6b07aa1](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6b07aa1ab6290f053aa32bfd728fb93ae284c052))
- Don't attempt a flood of requests when api is down - ([1c02622](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1c0262281724bc651adb76271e7cb64728371725))
- Incorrect gorm table - ([386aa74](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/386aa74dd4dc8a6ad172bf2efd40ae8663b389b3))
- Add ui support for 429s - ([88bda7c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/88bda7c255e378e56567e7465bceaa9c91a7fb65))
- Handle error-d values in results table - ([aad9ed2](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/aad9ed2c5085bec20082f6d2b89cb86cc5aee17e))

### ⚡ Performance

- Add index for timeseries query - ([8672434](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/86724348a85c6743a9e34c2f8a8b114dd0d89331))

### 🧪 Testing

- Add checks for limiter - ([a48206f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a48206fd389f8f7a22e00f3592fe812310def788))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([2c6b3d5](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/2c6b3d5c8497040a28ef45fef70cc6a460ff3e6e))
- Add release scripts - ([270f8bc](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/270f8bc4fc99a12a12f6c26d1fd64863452cbaab))

## ❤️ New Contributors

* @renovate-operator-k3s[bot] made their first contribution## [2.0.3](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.0.2..v2.0.3) - 2026-10-22

### 🐛 Bug Fixes

- Add build version - ([2a78b46](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/2a78b46fd32d0c9a21156e7fe777a5688568a384))
- Add build version - ([a94abdd](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a94abdd0e89d8f32e56deaf79df51a6fe8a7633a))
## [2.0.2](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.0.1..v2.0.2) - 2026-10-22

### 🐛 Bug Fixes

- Allow overridin francis hostname - ([5769e06](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5769e06dec3ecac8411c894bc6b6af7d5377ed65))
## [2.0.1](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.0.0..v2.0.1) - 2026-10-22

### 🐛 Bug Fixes

- Api url path in dockefile - ([cf1c6c7](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/cf1c6c7a3f070f6954c5fe816fcbda3e3b8fb85c))
## [2.0.0] - 2026-10-22

### 🚀 Features

- Actor runs tests - ([ffdcaa6](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ffdcaa62a6dbf29f6205ca763a39e4976374e390))
- Add cronjob to run tests on a schedule - ([1230c00](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1230c0051504d026fa1e80aaa7fcc2207d9dc502))
- Add handler to list records - ([6e1646e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6e1646e3bc772c77c980ce39e73af28967469f25))
- Add local docker-compose - ([ac14817](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ac1481727fa8ce9f7a10ff4cef79826494e36b58))
- Add ui with button to trigger test - ([1a091a1](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1a091a1a9930dc5796c923b65f600f4ab5b069de))
- Add landing page with latest result - ([1fc9107](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1fc910762810a8211de7ad9075b7e5fc749cf0fb))
- Added latest results table - ([41c8f92](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/41c8f925d7c5e9dc65eca56bfff7574baa31edc5))
- Dynamic loading of data across components - ([3544233](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/3544233683ce64a018cbb777470ab62d773f3646))
- Add timeseries graphs - ([5717b5a](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5717b5a8dab4476435dfb97b02f57345c9c678a3))
- Add ui into docker image - ([cbf0f7a](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/cbf0f7aa28522482c066ee44099247afe6fb9e6a))

### 🐛 Bug Fixes

- Move api endpoint path - ([33a778b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/33a778b4b41a832bba0aea4eb6d44f5509553e1c))
- Add endpoint to get speedtest result - ([a38738f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a38738f93e33e97e902c359990f1b0dc98cf9c98))
- Add dockerfile - ([57e718e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/57e718ed55fdb4b563c75cb13becf018ba0feb76))
- Use bytes in db - ([e0c303e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e0c303e226da14bd0cf8bd0cdc4ededeec08803d))
- Move to using raw bit per sec values instead of formatted - ([f34181b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/f34181b6581ffcb13dc410a1806fddf9b4f164c1))
- Add metrics - ([ff816dd](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ff816ddd7d5f5cd971e5242aa6b963d78e928484))
- Add mock graphs for stats - ([7a6ab6c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7a6ab6c6f0b791ed89de8f77b2f4cec83804cace))
- Disable run button when test is inflight - ([06ae93a](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/06ae93a95c5a66056278d2e3667312923710f550))
- Import v2 package - ([84f2049](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/84f204958af172abbff4ca7a5a27d3b03c7707cf))

### 🚜 Refactor

- Move results to own package - ([6fd2fe3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6fd2fe3ecf93ab8c7ce35afc74b496eec2a80354))

### 🧪 Testing

- Add results repo tests - ([2db1105](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/2db1105118e8d5e9f8d2abbf4be0af02313fb807))
- Add tests for listing - ([b82b51c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b82b51cfdc64af36cf6a78ea1dba35fdc0caf3ac))
- Split them out a bit - ([908f534](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/908f53458bc019692446c465a9de3cba542f1d01))
- Remove hist test - ([4284371](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/42843712fc32a162691ad8ec46ef1ba514233292))

### ⚙️ Miscellaneous Tasks

- Initial empty branch - ([7614b56](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7614b5638b41374460f6bec19698d9f0e8988bcf))
- Add test workflow - ([def71ad](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/def71ad743d0880d8dcf4ed12d20ba16e63c18d6))
- Add gitignore - ([6eadae3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6eadae3973ceeedf29e66b24dd5fbbf7df115147))
- Remove leftover frontend things - ([26ce4cd](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/26ce4cd7c45712c674aa9245a15a53af9a0601ba))
- Add readme - ([dab9cb7](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/dab9cb7fde37770e5731d2226095a283b3b8cc5e))
- Add license - ([e319d5c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e319d5cadeb0ed664c647b9acd89809908abeb27))
- Add build job - ([c405bed](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/c405bed51f27acb7f39b2719514a6f5fc4beb4ce))
- Add build job - ([05c2c8d](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/05c2c8d16ca285fad7e988ca2b8807314474abb2))
- Add build job - ([1dc634e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1dc634e75e7c642b116bd4a446453666685b961e))
- Add build job - ([90724db](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/90724dbeb4900c59accc223d464be012dc27fa54))

## ❤️ New Contributors

* @henrywhitaker3 made their first contribution
