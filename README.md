# cli

[![Build Status](https://github.com/pipego/cli/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/pipego/cli/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/pipego/cli/branch/main/graph/badge.svg?token=NODHGUZJ9X)](https://codecov.io/gh/pipego/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipego/cli)](https://goreportcard.com/report/github.com/pipego/cli)
[![License](https://img.shields.io/github/license/pipego/cli.svg)](https://github.com/pipego/cli/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/pipego/cli.svg)](https://github.com/pipego/cli/tags)



## Introduction

*cli* is the CLI of [pipego](https://github.com/pipego) written in Go.



## Prerequisites

- Go >= 1.18.0



## Run

```bash
version=latest make build
./bin/cli --config-file="$PWD"/test/config/config.yml --runner-file="$PWD"/test/data/runner.json --scheduler-file="$PWD"/test/data/scheduler.json
```



## Docker

```bash
version=latest make docker
docker run -v "$PWD"/test:/tmp ghcr.io/pipego/cli:latest --config-file=/tmp/config/config.yml --runner-file=/tmp/data/runner.json --scheduler-file=/tmp/data/scheduler.json
```



## Usage

```
usage: cli --config-file=CONFIG-FILE --runner-file=RUNNER-FILE --scheduler-file=SCHEDULER-FILE [<flags>]

pipego cli

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --version                  Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
  --runner-file=RUNNER-FILE  Runner file (.json)
  --scheduler-file=SCHEDULER-FILE
                             Scheduler file (.json)
```



## Settings

*cli* parameters can be set in the directory [config](https://github.com/pipego/cli/blob/main/config).

An example of configuration in [config.yml](https://github.com/pipego/cli/blob/main/config/config.yml):

```yaml
apiVersion: v1
kind: cli
metadata:
  name: cli
spec:
  runner:
    host: 127.0.0.1
    port: 29090
  scheduler:
    host: 127.0.0.1
    port: 28082
```



## License

Project License can be found [here](LICENSE).



## Reference

- [deploy](https://github.com/pipego/deploy)
- [runner](https://github.com/pipego/runner)
- [scheduler](https://github.com/pipego/scheduler)
