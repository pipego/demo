# demo

[![Build Status](https://github.com/pipego/demo/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/pipego/demo/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/pipego/demo/branch/main/graph/badge.svg?token=y5anikgcTz)](https://codecov.io/gh/pipego/demo)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipego/demo)](https://goreportcard.com/report/github.com/pipego/demo)
[![License](https://img.shields.io/github/license/pipego/demo.svg)](https://github.com/pipego/demo/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/pipego/demo.svg)](https://github.com/pipego/demo/tags)



## Introduction

*demo* is the demo of [pipego](https://github.com/pipego) written in Go.



## Prerequisites

- Go >= 1.18.0



## Run

```bash
version=latest make build
./bin/demo --config-file="$PWD"/test/config/config.yml --runner-file="$PWD"/test/data/runner.json --scheduler-file="$PWD"/test/data/scheduler.json
```



## Docker

```bash
version=latest make docker
docker run -v "$PWD"/test:/tmp ghcr.io/pipego/demo:latest --config-file=/tmp/config/config.yml --runner-file=/tmp/data/runner.json --scheduler-file=/tmp/data/scheduler.json
```



## Usage

```
usage: demo --config-file=CONFIG-FILE --runner-file=RUNNER-FILE --scheduler-file=SCHEDULER-FILE [<flags>]

pipego demo

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --version                  Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
  --runner-file=RUNNER-FILE  Runner file (.json)
  --scheduler-file=SCHEDULER-FILE
                             Scheduler file (.json)
```



## Settings

*demo* parameters can be set in the directory [config](https://github.com/pipego/demo/blob/main/config).

An example of configuration in [config.yml](https://github.com/pipego/demo/blob/main/config/config.yml):

```yaml
apiVersion: v1
kind: demo
metadata:
  name: demo
spec:
  runner:
    host: 127.0.0.1
    port: 29090
    timeout: 300
  scheduler:
    host: 127.0.0.1
    port: 28082
    timeout: 300
```



## License

Project License can be found [here](LICENSE).



## Reference

- [deploy](https://github.com/pipego/deploy)
- [runner](https://github.com/pipego/runner)
- [scheduler](https://github.com/pipego/scheduler)
