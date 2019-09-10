# Architecture

**source{d} Community Editon** provides a frictionless experience for trying source{d} for Code Analysis.

## Technical Architecture

The `sourced` binary, a single CLI binary [written in Go](https://github.com/dpordomingo/sourced-ce/tree/c91e1e000bdf08fe29aae1fd05aa5ea61c785786/cmd/sourced/main.go), is the user's main interaction mechanism with **source{d} CE**. It is also the only piece \(other than Docker\) that the user will need to explicitly download on their machine to get started.

The `sourced` binary manages the different installed environments and their configurations, acting as a wrapper of Docker Compose.

The whole architecture is based on Docker containers, orchestrated by Docker Compose and managed by `sourced`.

## Docker Set Up

In order to make this work in the easiest way, some design decisions were made:

### Isolated Environments.

_Read more in_ [_Working With Multiple Data Sets_](https://github.com/dpordomingo/sourced-ce/tree/c91e1e000bdf08fe29aae1fd05aa5ea61c785786/usage/multiple-datasets.md)

Each dataset runs in an isolated environment, and only one environment can run at the same time. Each environment is defined by one `docker-compose.yml` and one `.env`, stored in `~/.sourced`.

### Docker Naming

All the Docker containers from the same environment share its prefix: `srcd-<HASH>_` followed by the name of the service running inside, e.g `srcd-c3jjlwq_gitbase_1` and `srcd-c3jjlwq_bblfsh_1` will contain gitbase and babelfish for the same environment.

### Docker Networking

In order to provide communication between the multiple containers started, all of them are attached to the same single bridge network. The network name also has the same prefix than the containers inside the same environment, e.g. `srcd-c3jjlwq_default`.

Some environment services can be accessed from the outside, using their exposed port and connection values:

* `bblfsh`:
  * port: `9432`
* `gitbase`:
  * port: `3306`
  * database: `gitbase`
  * user: `root`
* `metadatadb`:
  * port: `5433`
  * database: `metadata`
  * user: `metadata`
  * password: `metadata`
* `sourced-ui`:
  * port: `8088`

### Persistence

To prevent losing data when restarting services, or upgrading containers, the data is stored in volumes. These volumes also share the same prefix with the containers in the same environment, e.g. `srcd-c3jjlwq_gitbase_repositories`.

These are the most relevant volumes:

* `gitbase_repositories`, stores the repositories to be analyzed
* `gitbase_indexes`, stores the gitbases indexes
* `metadata`, stores the metadata from GitHub pull requests, issues, users...
* `postgres`, stores the dashboards and charts used by the web interface

