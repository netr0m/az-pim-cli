# Changelog

## [1.3.0](https://github.com/netr0m/az-pim-cli/compare/v1.2.0...v1.3.0) (2024-10-21)


### Features

* Support for Entra roles ([#61](https://github.com/netr0m/az-pim-cli/issues/61)) ([dd9ed19](https://github.com/netr0m/az-pim-cli/commit/dd9ed193c7bee3a85ad3cc62ada4bc2630378393))

## [1.2.0](https://github.com/netr0m/az-pim-cli/compare/v1.1.0...v1.2.0) (2024-10-21)


### ⚠ BREAKING CHANGES

* use proper terms for 'azure resources' type ([#59](https://github.com/netr0m/az-pim-cli/issues/59))

### Features

* use proper terms for 'azure resources' type ([#59](https://github.com/netr0m/az-pim-cli/issues/59)) ([6411902](https://github.com/netr0m/az-pim-cli/commit/641190289f99d2599d7dd789c5c3ea10845746ae))

## [1.1.0](https://github.com/netr0m/az-pim-cli/compare/v1.0.1...v1.1.0) (2024-09-13)


### Features

* add support for specifying 'ticket number' and 'ticket system' ([#56](https://github.com/netr0m/az-pim-cli/issues/56)) ([a62c52f](https://github.com/netr0m/az-pim-cli/commit/a62c52ff158a018d46598fa6c631ebc020c52d53))

## 1.0.1 (2024-07-01)


### Features

* activate roles ([7cdb3be](https://github.com/netr0m/az-pim-cli/commit/7cdb3be77fe393028096d066192a6c1631b3ac3d))
* add 'version' command ([#30](https://github.com/netr0m/az-pim-cli/issues/30)) ([e24a15f](https://github.com/netr0m/az-pim-cli/commit/e24a15f6fb1aa020e6e7191080c3b56363eac355))
* add reason to activate command ([#4](https://github.com/netr0m/az-pim-cli/issues/4)) ([8b43135](https://github.com/netr0m/az-pim-cli/commit/8b4313595e4b534c304619c973d42e2c8e8b1d35))
* check for various request status types ([#14](https://github.com/netr0m/az-pim-cli/issues/14)) ([57e4472](https://github.com/netr0m/az-pim-cli/commit/57e447247280dc092cc2b9ee817a53b599b47ae9))
* dry-run for 'activate' ([#22](https://github.com/netr0m/az-pim-cli/issues/22)) ([05c4095](https://github.com/netr0m/az-pim-cli/commit/05c40956017909a14f3015f2de10c4a5e43303e2))
* list eligible roles ([eb3e15a](https://github.com/netr0m/az-pim-cli/commit/eb3e15ae475d065613c1cb816dc6082e9d008c76))
* support for PIM Entra groups ([#16](https://github.com/netr0m/az-pim-cli/issues/16)) ([6fddc87](https://github.com/netr0m/az-pim-cli/commit/6fddc870a990bc6065b8dd053544fc141421428f))
* support new Azure Entra ID PIM API ([#6](https://github.com/netr0m/az-pim-cli/issues/6)) ([700323c](https://github.com/netr0m/az-pim-cli/commit/700323cc0c90674f8d1b8fd9db6db96933e15bbc))
* use az-cli for auth ([95e7553](https://github.com/netr0m/az-pim-cli/commit/95e7553cd7142b0ba35f7054f4762b23764804d3))


### Bug Fixes

* **activate:** Role selection on `activate` selects incorrect role ([#8](https://github.com/netr0m/az-pim-cli/issues/8)) ([6cb1079](https://github.com/netr0m/az-pim-cli/commit/6cb1079b62cabf219232c9e829198d70b4b122e8))
* fix casing role on activate ([#3](https://github.com/netr0m/az-pim-cli/issues/3)) ([9d92cff](https://github.com/netr0m/az-pim-cli/commit/9d92cff54a4515eb44e6226c623fe8f59cf9817c))
* use exact matching for the role selection ([#12](https://github.com/netr0m/az-pim-cli/issues/12)) ([0bf37e6](https://github.com/netr0m/az-pim-cli/commit/0bf37e6db2e648179442326c0b101328e4fd7e82))


### Documentation

* **github:** add project guidelines ([#31](https://github.com/netr0m/az-pim-cli/issues/31)) ([5fc195b](https://github.com/netr0m/az-pim-cli/commit/5fc195bda5e78fd66b0fc996b3259d380b40f102))
* initial docs ([c315b5c](https://github.com/netr0m/az-pim-cli/commit/c315b5c44dab5102e8a7678c09e3c81d35f87a09))


### Continuous Integration

* add release-please workflow ([#36](https://github.com/netr0m/az-pim-cli/issues/36)) ([67f357d](https://github.com/netr0m/az-pim-cli/commit/67f357d1dfb1a2bc981ad257085757e59d934b90))
* add workflow triggered by merge to main ([#37](https://github.com/netr0m/az-pim-cli/issues/37)) ([4b24cf9](https://github.com/netr0m/az-pim-cli/commit/4b24cf90b8a58a5a71c36347149418b233fa038b))
