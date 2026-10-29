## [2.7.2](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.7.1..v2.7.2) - 2026-10-29

### 🐛 Bug Fixes

- *(config)* Allow comma-separated server ids - ([4905138](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4905138b0be35b5c19781c86552275879f479450))
- *(ookla)* Pick random id from configured servers - ([0142748](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0142748461dd8a88d3726386523bf645ae5b74c8))
- *(ui)* Add handling for multiple configured server ids - ([2434155](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/24341558838e5ff89cbc5047026f6d408c97f3a9))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([88caf0c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/88caf0c3ad935a3f42b4ade74efc7d3fea4e3456))
## [2.7.1](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.7.0..v2.7.1) - 2026-10-29

### 🐛 Bug Fixes

- *(deps)* Update module github.com/slack-go/slack to v0.32.0 - ([d5bf856](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d5bf85621f557ba03e65eafa77f0ef35b5a4378d))
- *(notifications)* Track notification history to stop repeated notifications when a single channel fails - ([d6bf9b0](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d6bf9b07f26096bd0cafa2173089cf83ba18b90e))
- *(notifications)* Filter out channels when they are not configured for the phase - ([ce07a67](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ce07a67ed4c918198c2257f7da1352f43ef76309))
- *(notifications)* Add logging - ([a649282](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a649282dd630e0858b05876789709ce9ad429647))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([d199513](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d199513d19e126e583d3a571c2affb1045ae6539))
- Remove duplicate tests - ([b4afdd3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b4afdd379f75d475ab93bb5caa8f2cc5ce89ebe6))
- Remove leftover comments - ([656ad2c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/656ad2c7686a895052c219635c4c2a31e001d6c3))
## [2.7.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.6.2..v2.7.0) - 2026-10-27

### 🚀 Features

- *(notifications)* Add slack notification channel - ([44ff29c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/44ff29ca30e7d5f95638fea4c31db26fdf15d086))
- *(notifications)* Add discord notifications - ([7ff632f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7ff632fb029c979e50fd8a6e16482c3a98d11e33))

### 🐛 Bug Fixes

- *(api)* Redact more server info for unauthenticated users - ([a2ae03c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a2ae03c03bf8fc31c8b4be9800491eb176723123))
- *(notifications)* Use constants for standardised message titles - ([60f99d2](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/60f99d2911f2d433bd812cf8a19df1aa5dcf2eba))
- *(notifications)* Use constrained context for discord and slack notifications - ([f9dff72](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/f9dff72d5217272a3d3b29a6c06ffdac9f7d115e))

### 🚜 Refactor

- *(test)* Use testify instead of manual assertions - ([642f58d](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/642f58d8e4644bc1b5a286cbf319847ab95efdae))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([d88fc1e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d88fc1ec3d251adacd665b7d862f4d36dd7ff19d))
- *(docs)* Add info on new notification channels - ([3b6aa3f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/3b6aa3f6c5744fff86714218e2423fc02c0a054a))
- *(docs)* Update openapi schema with auth info - ([30dd18f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/30dd18f59acffa42bdc34906b04b27963c445d9b))
## [2.6.2](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.6.1..v2.6.2) - 2026-10-27

### 🐛 Bug Fixes

- *(api)* Redact more server info for unauthenticated users - ([c1451b3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/c1451b332364964dd363e7084e73913b6ea22d1b))
- *(ui)* Nicer title and description for pick server step - ([7659674](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/76596749d0adaa013ce80edda1cf69ed06b998bb))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([18aa2c6](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/18aa2c685fc58fdadd780a1d7b0b3d9f10abcea8))
## [2.6.1](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.6.0..v2.6.1) - 2026-10-27

### 🐛 Bug Fixes

- *(deps)* Update module github.com/italypaleale/francis to v0.2.0 - ([fdecc67](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fdecc67a88981940d4c3f8d47e9d086112b9a294))
- *(ookla)* Pick random server from nearest list when one is not configured - ([2329ba4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/2329ba43a682e8a1f68a9dcfbb49587fdb60de2d))
- *(runtime)* Add optional runtime management api activation - ([f517be4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/f517be4e94a55a523a2932dc1ad2f87e24dc74ac))
- *(ui)* Replace vue default favicon - ([4ea508f](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4ea508f8e94e61a0eba1ce47e5a8bbb39bb74ae8))
- *(ui)* Result delete button not working in table view - ([83c96c0](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/83c96c0d1f96626ccf921b7390a013a6a6c3205a))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([60740c5](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/60740c5e1fa8aefc34e3965ceb52f1e261c3bf5f))
- Better handling for interrupts when running locally - ([313051e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/313051e93fbb34e5de9afc94a81ca61fe85850de))
## [2.6.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.5.0..v2.6.0) - 2026-10-26

### 🚀 Features

- *(notifications)* Add telegram notifications - ([4ec99af](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4ec99af7de14ff0397c4b655812c03421894780c))

### 🐛 Bug Fixes

- *(api)* Sanaitise raw output instead of nulling it - ([22ac945](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/22ac945e84d04c9f4e4ec8605f52e67e286752c8))
- *(compat)* Added a ltest result handler compatible with gethomepage v2 widget - ([5e8d443](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5e8d443d3d53b4a76c6c874959b55727d5af57c1))
- *(telegram)* Update message formatting - ([6743e00](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6743e00488796086abc57b49b2af4e1ca43fa1b5))
- Skip telegram GetMe call on bootstrap - ([62d536b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/62d536b7d8325ed44c8852b656192f86376fc0da))
- Redact telegram details in notification errors - ([d9a7aef](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d9a7aefa182b615d15b5cf3cb8e3ca34c789f1b4))
- Add handling for negative bits to human values - ([92e4de4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/92e4de4dc2e35dffb35076aa252dd90dc5122211))

### 🚜 Refactor

- *(api)* Use new helpers for middleware auth status - ([4f641be](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4f641bebe5769c17d71268d7178e9b647defa140))
- *(notifications)* Split webhooks and errors into own files - ([eec7396](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/eec7396edaa386e1cdbbabd9068a3c5e8e573142))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([5b0eac9](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/5b0eac955eae0843f6ff9a3a31f039e2f95c954c))
- *(docs)* Add info on telegram notifications - ([11644dd](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/11644dd0b69f5f043eb08dce73bcdaaf68652fe0))
## [2.5.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.4.5..v2.5.0) - 2026-10-26

### 🚀 Features

- *(result)* Add api handler to restore legacy json backups - ([1baad7e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1baad7e95f022540fe42c67c80beacd0c0a80d20))
- *(results)* Add function to batch restores of results with de-duplication - ([8c79de8](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/8c79de8e40e06d98caa18790d9daf003591b02c9))

### 🐛 Bug Fixes

- *(results)* Skip restoring results with a zero created_at field - ([aac576b](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/aac576bf8d05394e92a08b2f49e5f7c588c4e16b))
- *(results)* Chunk existing timestamp lookup - ([95c4bd6](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/95c4bd66fb0821baab213b5ec25ca84c4166898a))
- *(ui)* Add a restore action to the settings page - ([b939131](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b9391319f33bf64ffb894046fcb0538c1543d7e2))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog [skip ci] - ([3f895b3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/3f895b300fce5eb3a8aa9d8c93f9211adcb7c6d4))
- Update readme with example docker compose - ([1ecfbe5](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/1ecfbe5867e82abb89d2ccee557141cbac5cc16a))
## [2.4.5](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.4.4..v2.4.5) - 2026-10-25

### 🐛 Bug Fixes

- *(api)* Add operation ids to spec handlers - ([7aa4553](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/7aa45531bd1973485963f2631b174a7e27d99924))
- *(api)* Add openapi spec definitions for enums - ([d9d5cdb](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d9d5cdbfb26d4e8425805c86548d716ce4320b70))
- *(api)* Remove local replaced module - ([e360236](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e36023619ebf855667eeb09f2bc22627cc678f04))
- *(deps)* Update module modernc.org/sqlite to v1.61.0 - ([0c5bff1](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0c5bff16586ee2afcb90eddabb2cb598ee61c184))
- *(ui)* Light mode raw result syntax highlighting - ([4f0a00c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4f0a00cec8b24c61fe27f9a40292b82592377128))
- *(ui)* Run button navigates to result when test is in progress - ([12fb4e4](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/12fb4e45556a96e257a4d3247029289b93b4003a))
- *(ui)* Use new generated client types - ([22d8423](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/22d8423fb34d41d04ecee2bea66e39c2b454cfc6))

### ⚙️ Miscellaneous Tasks

- *(ci)* Tag ordering to still trigger - ([0449538](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/04495389ebadc40b7962f7aa45c2339ba8cd9157))
- *(docs)* Update changelog [skip ci] - ([e3494a7](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/e3494a77db0b35db3881ddbeb1d43a4819e5c5d9))
## [2.4.4](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.4.3..v2.4.4) - 2026-10-25

### 🐛 Bug Fixes

- *(db)* Use a non-cgo slite driver - ([9a13da2](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/9a13da224ac395621df324f80bd90c28b73673b1))
- *(ui)* Clean up typecheck errors - ([a78413a](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a78413a0eb2ff989a13947276a875713be0360ad))
- *(ui)* Add syntax highlighting to the raw result - ([2f90cea](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/2f90cea6c5010ca017810ccc441515a9629e05cb))

### 🧪 Testing

- *(ui)* Add a bunch of tests against a mock api - ([0cb2b02](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0cb2b02c1a927b76d266158f14a7e01cc68c33ad))
- *(ui)* Add tests for when auth is enabled - ([0653239](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0653239bb43f1891318820a2cf1fa07e556c6266))
- Added playwright commands - ([d5f370e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d5f370e8c2bcd069afb4b3a83f41f018f5add04c))

### ⚙️ Miscellaneous Tasks

- *(ci)* Skip ci on changelog push - ([d9a0657](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d9a065737279935227eed348453b18ba81b205fa))
- *(docs)* Update changelog [skip ci] - ([aa99beb](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/aa99beb604a1f06c22ed1a18bc851e3bde81ee8a))
## [2.4.3](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.4.2..v2.4.3) - 2026-10-25

### 🐛 Bug Fixes

- *(o11y)* Add otel metrics to prometheus registry - ([eae95e3](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/eae95e3afc1ecf94015a7e32213560acf75f6432))

### 📚 Documentation

- Update README.md - ([ea051dc](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ea051dc383b13d2b32f59f1ec903e3be4f15f476))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([bf0505c](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/bf0505c8eda43f5a50004c0a7be124384b954fe5))
## [2.4.2](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.4.1..v2.4.2) - 2026-10-24

### 🐛 Bug Fixes

- *(results)* Json - struct tag stopping raw output being shared in workflow outputs - ([fdb0962](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fdb09625aabe7b1faeadf38cc0586d50fcfa96a9))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([6ed2e19](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/6ed2e1951858246bb4b62b7c10c16237abde0a56))
## [2.4.1](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.4.0..v2.4.1) - 2026-10-24

### 🐛 Bug Fixes

- *(auth)* Tighten middleare failure modes - ([807fc01](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/807fc01238e88b65e46ca5ea1aef1095d76a17d8))
- *(db)* Open postgres connection with slog logger - ([fc723a9](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fc723a93c7e4191017a0707c83291a1b0a05f251))
- *(results)* Store raw result output and display it - ([b953c4e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b953c4ef433591086318774be3ffbeaa2aeaf1b8))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update changelog - ([63ded00](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/63ded006c757bd04cbb854fa19c08ba2c6bb3820))
- Run tests on pushes to v2 to cache go mod and builds - ([c8f937a](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/c8f937ac61629198e8ba0c11a81aa2d1c53c90ca))
- Upload test coverage reports - ([972cc84](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/972cc8471b9b67d28a229b631ea522b348b8331b))
## [2.4.0](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.3.0..v2.4.0) - 2026-10-24

### 🚀 Features

- *(config)* Added config for webhook notifications - ([ac4f4bd](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/ac4f4bdfdfcfc361686c9606a08820d0cf2b6a84))
- *(notifications)* Add workflow support for multiple notifications and phases - ([0f85cad](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/0f85cadd0a6c1f7ea33e98f770d33aaa491c90fa))

### 🐛 Bug Fixes

- *(config)* Notification url parsing - ([65a6285](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/65a6285fd48cc6329c2ce4d0a7a13f9f8509d093))
- *(notifications)* Assert http(s) scheme to the configured url - ([4d4c35e](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/4d4c35e4d6b656317bf9b8293a401f477cc9f4b0))
- *(ui)* Update timeline for new speedtest workflow - ([fd2affa](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fd2affad597700b9f58c13fbcb583d29d4a598ea))
- *(ui)* Add card to settings page to display configured notification channels - ([d4313ed](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/d4313ed6f20ae1657c2ca387034a61e103e5e742))
- *(ui)* Add handling for null response from notifications endpoint - ([b3cc0a9](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/b3cc0a9f2241cb986f11e626d004ce129620985e))
- *(workflow)* Mark failed healthchecks as optional - ([3f31388](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/3f31388d44f8e957f136e08e5f537d19462bc01f))

### ⚙️ Miscellaneous Tasks

- *(docs)* Update readme and add some docs - ([fa13268](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/fa13268dc2cb6d2637a61551f2d7f38f4f2c1391))
- *(docs)* Update changelog - ([c294186](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/c2941866f1348d981d00a678085455913f3c3ec3))
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

- *(docs)* Update changelog - ([a127a44](https://github.com/henrywhitaker3/Speedtest-Tracker/commit/a127a44c4492b0c54fe84cc93f2d28f8989e7398))
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
## [2.0.3](https://github.com/henrywhitaker3/Speedtest-Tracker/compare/v2.0.2..v2.0.3) - 2026-10-22

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
