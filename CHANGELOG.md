## [1.9.1](https://github.com/PremoWeb/chadburn/compare/v1.9.0...v1.9.1) (2025-03-16)


### Bug Fixes

* handle Docker socket connection issues with backoff and prevent job de-registration (Issue [#117](https://github.com/PremoWeb/chadburn/issues/117)) ([34ee8b6](https://github.com/PremoWeb/chadburn/commit/34ee8b613230d77c532d3d993ad4bb07b33f6983))
* handle EOF errors in Docker event watching with exponential backoff ([466ce1a](https://github.com/PremoWeb/chadburn/commit/466ce1a1a485ea74aa9163d1750d2fa2c004c310))
* lowercase error message to comply with Go style guide (ST1005) ([a76f3bc](https://github.com/PremoWeb/chadburn/commit/a76f3bc0bce714feaa20d68811dfd7440630d3ec))

# [1.9.0](https://github.com/PremoWeb/chadburn/compare/v1.8.2...v1.9.0) (2025-03-16)


### Bug Fixes

* asset loading issues by adding .nojekyll and setting paths.relative to false ([c25444d](https://github.com/PremoWeb/chadburn/commit/c25444df73ab78b34c0bd8c24786144f4bce144d))
* **docs:** remove _routes.json file ([3b73955](https://github.com/PremoWeb/chadburn/commit/3b7395578b693ceecfbe2f350ca4207795abb46a))
* **docs:** remove routes configuration from Cloudflare adapter ([55575b1](https://github.com/PremoWeb/chadburn/commit/55575b16079cd2740ea367bae5b50a406562574e))
* **docs:** update Cloudflare Pages deployment directory ([8f21143](https://github.com/PremoWeb/chadburn/commit/8f21143051183ccf8d69e80e0286f686da4e25ae))
* **docs:** update Cloudflare Pages deployment directory to use .svelte-kit/cloudflare ([05cdf64](https://github.com/PremoWeb/chadburn/commit/05cdf64f59baf2668bcfa557c3f81f1edf23a81e))
* **docs:** update Cloudflare Pages deployment to use Wrangler v4 ([65d7e0c](https://github.com/PremoWeb/chadburn/commit/65d7e0c01a8ed6af4caf11222a8d466bc06dcf96))
* update contributors script to use ES modules syntax ([c6be872](https://github.com/PremoWeb/chadburn/commit/c6be8724fe9512d9f04da31dc3d1edb6047c30c8))


### Features

* add metrics setup documentation page ([f041813](https://github.com/PremoWeb/chadburn/commit/f041813cffd6b3467b8b3b36f667be4a399c48f1))
* add version display from VERSION file ([c606a67](https://github.com/PremoWeb/chadburn/commit/c606a67b4594f1f0feb7619056b86b4ee9eba09b))

## [1.8.2](https://github.com/PremoWeb/chadburn/compare/v1.8.1...v1.8.2) (2025-03-15)


### Bug Fixes

* update GitHub Pages deployment workflow to use JamesIves action ([335f7e8](https://github.com/PremoWeb/chadburn/commit/335f7e84ef83e0e58e78064f1c72a6022752cd1a))

## [1.8.1](https://github.com/PremoWeb/chadburn/compare/v1.8.0...v1.8.1) (2025-03-15)


### Bug Fixes

* update GitHub Pages deployment workflow ([8e62de3](https://github.com/PremoWeb/chadburn/commit/8e62de3d40584d73b5fde16544531453808d293b))

# [1.8.0](https://github.com/PremoWeb/chadburn/compare/v1.7.4...v1.8.0) (2025-03-15)


### Features

* configure GitHub Pages with custom domain chadburn.dev ([9dd358c](https://github.com/PremoWeb/chadburn/commit/9dd358c56568dce3f72a25c7ca099219bde8ddd4))

## [1.7.4](https://github.com/PremoWeb/chadburn/compare/v1.7.3...v1.7.4) (2025-03-15)


### Bug Fixes

* remove redundant header from documentation pages ([ecf0a44](https://github.com/PremoWeb/chadburn/commit/ecf0a44cd76b8230a6760c297d528e3660ba5dd9))

## [1.7.3](https://github.com/PremoWeb/chadburn/compare/v1.7.2...v1.7.3) (2025-03-15)


### Bug Fixes

* improve client-side navigation in Svelte docs ([d1a2784](https://github.com/PremoWeb/chadburn/commit/d1a278469781d484230748f892938a35b4cb5502))

## [1.7.2](https://github.com/PremoWeb/chadburn/compare/v1.7.1...v1.7.2) (2025-03-15)


### Bug Fixes

* update tests after migration to official Docker client ([26103ce](https://github.com/PremoWeb/chadburn/commit/26103ce692a1ce888859babceeeb9a2305642fdf))

## [1.7.1](https://github.com/PremoWeb/chadburn/compare/v1.7.0...v1.7.1) (2025-03-14)


### Bug Fixes

* update tests to work with the new Docker client implementation ([5ee9329](https://github.com/PremoWeb/chadburn/commit/5ee93297fb0f9bfb44123c1664c5775d11606e85))

# [1.7.0](https://github.com/PremoWeb/chadburn/compare/v1.6.0...v1.7.0) (2025-03-14)


### Bug Fixes

* prevent context leak in OfficialDockerHandler ([f1c83ff](https://github.com/PremoWeb/chadburn/commit/f1c83ff7b3ed3bb76eb0e566a62dbe0cf2cfff51))


### Features

* remove all references to the legacy polyfill dependency ([9792563](https://github.com/PremoWeb/chadburn/commit/97925638200cea0f9249dd26cb159e1ff664a4b9))
* replace fsouza/go-dockerclient with official Docker client ([d607852](https://github.com/PremoWeb/chadburn/commit/d6078526a11913f00a5aa31d156783e05aef93dd))

# [1.6.0](https://github.com/PremoWeb/chadburn/compare/v1.5.1...v1.6.0) (2025-03-14)


### Bug Fixes

* avoid copying mutex in LifecycleJob.Run method ([aa52ac1](https://github.com/PremoWeb/chadburn/commit/aa52ac11890c932d8921813a3968e962566464b9))


### Features

* migrate from fsouza/go-dockerclient to official Docker client library ([commit-hash](https://github.com/PremoWeb/chadburn/commit/commit-hash))
* add support for container lifecycle jobs (issue [#68](https://github.com/PremoWeb/chadburn/issues/68)) ([defc60f](https://github.com/PremoWeb/chadburn/commit/defc60f9d087b15d368008f3ac321e71679c30cb))

## [1.5.1](https://github.com/PremoWeb/chadburn/compare/v1.5.0...v1.5.1) (2025-03-14)


### Bug Fixes

* fix tests after removing Pull field from RunJob struct ([90ee9b0](https://github.com/PremoWeb/chadburn/commit/90ee9b021b8d3938f686830334b292e2ff361ef9))

# [1.5.0](https://github.com/PremoWeb/chadburn/compare/v1.4.0...v1.5.0) (2025-03-14)


### Features

* add support for variables in job commands ([8b392ea](https://github.com/PremoWeb/chadburn/commit/8b392ea965272afca328b844b1b58b826e157abd))
* improve documentation for job-run with Docker Compose ([commit-hash](https://github.com/PremoWeb/chadburn/commit/commit-hash)), closes [#70](https://github.com/PremoWeb/chadburn/issues/70)

# [1.4.0](https://github.com/PremoWeb/chadburn/compare/v1.3.8...v1.4.0) (2025-03-14)


### Features

* add workdir parameter to job-exec ([7979d52](https://github.com/PremoWeb/chadburn/commit/7979d5202f17e399cbaf8379635b9278cc282180)), closes [#100](https://github.com/PremoWeb/chadburn/issues/100)
* add support for variables in job commands ([8b392ea](https://github.com/PremoWeb/chadburn/commit/8b392ea)), closes [#66](https://github.com/PremoWeb/chadburn/issues/66)

## [1.3.8](https://github.com/PremoWeb/chadburn/compare/v1.3.7...v1.3.8) (2025-03-14)


### Bug Fixes

* update version management in semantic-release workflow ([9ff2d1f](https://github.com/PremoWeb/chadburn/commit/9ff2d1f45d79303df420cd0967817a22a9704fc7)), closes [#100](https://github.com/PremoWeb/chadburn/issues/100)

## [1.3.7](https://github.com/PremoWeb/chadburn/compare/v1.3.6...v1.3.7) (2025-03-14)


### Bug Fixes

* update GitHub Actions workflow to use Go 1.23 ([7316d1d](https://github.com/PremoWeb/chadburn/commit/7316d1d15b01c6bee20d58112767b5f33d257131))

## [1.3.6](https://github.com/PremoWeb/chadburn/compare/v1.3.5...v1.3.6) (2025-03-14)


### Bug Fixes

* update publish_release workflow to properly handle tag-based releases ([4ad3071](https://github.com/PremoWeb/chadburn/commit/4ad3071e852b4e140fc8bc81841d56012b083c01))

## [1.3.5](https://github.com/PremoWeb/chadburn/compare/v1.3.4...v1.3.5) (2025-03-14)


### Bug Fixes

* update GitHub release action to use wangyoucao577/go-release-action with correct parameters ([2d47d69](https://github.com/PremoWeb/chadburn/commit/2d47d6926cfa2c9d52d9b6764a608c31e16a1b5d))

## [1.3.4](https://github.com/PremoWeb/chadburn/compare/v1.3.3...v1.3.4) (2025-03-14)


### Bug Fixes

* add golangci-lint config to disable problematic linters ([f01cdcd](https://github.com/PremoWeb/chadburn/commit/f01cdcdb65d543c714a357c9bf96abe72e941f97))

## [1.3.3](https://github.com/PremoWeb/chadburn/compare/v1.3.2...v1.3.3) (2025-03-14)


### Bug Fixes

* update Dockerfile to use Go 1.19 instead of Go 1.23.2 ([82e4f28](https://github.com/PremoWeb/chadburn/commit/82e4f28b90a6c56078ae0e7a0894305b69714b43))
* update Go version to 1.23 ([09ba910](https://github.com/PremoWeb/chadburn/commit/09ba9104a83d10c6bb8689ae7aee4abc46bfc379))

## [1.3.2](https://github.com/PremoWeb/chadburn/compare/v1.3.1...v1.3.2) (2025-03-14)


### Bug Fixes

* update go.mod to use Go 1.19 and remove toolchain directive ([edb5477](https://github.com/PremoWeb/chadburn/commit/edb5477f9c019ae17f80e78e6f3dc80b51e0433d))

## [1.3.1](https://github.com/PremoWeb/chadburn/compare/v1.3.0...v1.3.1) (2025-03-14)


### Bug Fixes

* update go.mod to use Go 1.19 and remove toolchain directive ([bf47752](https://github.com/PremoWeb/chadburn/commit/bf477527a355e15245dc9aaaa1a3327cafc6b38a))

# [1.3.0](https://github.com/PremoWeb/chadburn/compare/v1.2.1...v1.3.0) (2025-03-14)


### Features

* implement Git hooks for commit message validation ([2ed31b8](https://github.com/PremoWeb/chadburn/commit/2ed31b8a15c5e8085ae2e6f6c069bd89086c25d6))

## [1.2.1](https://github.com/PremoWeb/chadburn/compare/v1.2.0...v1.2.1) (2025-03-14)


### Bug Fixes

* add note about semantic versioning to README ([ae96d78](https://github.com/PremoWeb/chadburn/commit/ae96d7810939d1bb5fb26f1a6921ddd213cd5fac))

# [1.2.0](https://github.com/PremoWeb/chadburn/compare/v1.1.0...v1.2.0) (2025-03-14)


### Bug Fixes

* update semantic release workflow to properly handle variable substitution ([610e445](https://github.com/PremoWeb/chadburn/commit/610e445cb28747b38d5c8a2895ed2ece932b3e34))


### Features

* implement automatic version bump system and update README badges ([7ff96c3](https://github.com/PremoWeb/chadburn/commit/7ff96c36dec1b73d0deb472dda1375a5281ca272))

## Unreleased

### Added
- Added support for container lifecycle jobs (`job-lifecycle`) that run once on container start or stop events.
