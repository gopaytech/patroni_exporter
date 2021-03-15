[![CI](https://github.com/gopaytech/patroni_exporter/workflows/Main%20Deployment/badge.svg)][ci]
[![Go Report Card](https://goreportcard.com/badge/github.com/gopaytech/patroni_exporter)][goreportcard]
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)][license]

[ci]: https://github.com/gopaytech/patroni_exporter/actions?query=workflow%3A%22Master+Deployment%22+branch%3Amaster
[goreportcard]: https://goreportcard.com/report/github.com/gopaytech/patroni_exporter
[license]: https://opensource.org/licenses/Apache-2.0

# Patroni Exporter for Prometheus
Simple server that scrapes Patroni stats and exports them via HTTP for Prometheus consumption.

## Getting Started

To run it:

```bash
./patroni_exporter [flags]
```

Help on flags:

```bash
./patroni_exporter --help
```

For more information check the [source code documentation][gdocs].

[gdocs]: http://godoc.org/github.com/gopaytech/patroni_exporter

## Usage

--patroni.host="http://localhost"
Specify Patroni API URL using the `--patroni.host` flag.
Specify Patroni API port using the `--patroni.port` flag.
```bash
./patroni_exporter --patroni.host="http://localhost" --patroni.port=8008
```

### Building

```bash
make build
```

### Testing

```bash
make test
```

## License

Apache License 2.0, see [LICENSE](https://github.com/gopaytech/patroni_exporter/blob/master/LICENSE).
