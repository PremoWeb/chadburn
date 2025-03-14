## Chadburn v1.0.9

This release adds support for placing job-run labels directly on target containers and updates the project dependencies to their latest versions.

### Features
- Added support for job-run labels on target containers, allowing users to schedule containers to start periodically without having to specify the container names in the Chadburn container's labels

### Changes
- Updated Go version from 1.20 to 1.23.0 (with toolchain 1.23.7)
- Updated all direct dependencies to their latest versions, including:
  - github.com/fsouza/go-dockerclient: v1.6.5 -> v1.12.1
  - github.com/jessevdk/go-flags: v1.4.0 -> v1.6.1
  - github.com/mitchellh/mapstructure: v1.3.3 -> v1.5.0
  - github.com/docker/docker: v1.4.2-0.20191101170500-ac7306503d23 -> v28.0.1+incompatible
  - github.com/prometheus/client_golang: v1.13.0 -> v1.21.1
- Fixed tests to work with the updated dependencies

### Compatibility
This release maintains backward compatibility with previous versions.
