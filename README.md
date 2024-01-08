# hrt - HTTP(S) Request Tool

[![Go Report Card](https://goreportcard.com/badge/github.com/majermarci/hrt)](https://goreportcard.com/report/github.com/majermarci/hrt)
[![Go Reference](https://pkg.go.dev/badge/github.com/majermarci/hrt.svg)](https://pkg.go.dev/github.com/majermarci/hrt)
![License](https://img.shields.io/github/license/majermarci/hrt?label=License)
[![Build Status](https://github.com/majermarci/hrt/actions/workflows/build.yaml/badge.svg)](https://github.com/majermarci/hrt/actions/workflows/build.yaml)
![Latest Pre-Release)](https://img.shields.io/github/v/release/majermarci/hrt?include_prereleases&label=pre-release&logo=github)
<!-- ![Latest Release)](https://img.shields.io/github/v/release/majermarci/hrt?logo=github) -->

The application is currently in its initial development stages and serves as a learning project. The ultimate goal is to provide a streamlined command-line interface alternative to tools such as Insomnia and httpie.
The application works with a simple and easy-to-back-up `yaml` configuration, which can also be customized per project and included in repos.

## Features

- Simple [yaml configuration](cmd/hrt/example_config.yaml) for organizing your request collection(s)
  - Specify headers, body and method for each request
  - Bearer Token and Basic Auth support
  - Support for local and global config files
  - Option to create a default config file
- Various option flags for running requests
  - Usage of specific TLS certificate, adding new CA chains to existing one or skipping certificate verification
  - Global timeout option
  - Verbose outputs showing request and TLS details
  - Option to run every request from the active collection file right after each other

## Usage

For more information on how to use `hrt`, you can run `hrt -h` in your terminal or refer to the [documentation](usage.md) page.

## Installing

If [Go](https://github.com/golang/go) is [installed](https://go.dev/doc/install) on your system, you can install the app using the following command:

```bash
go install github.com/majermarci/hrt/cmd/hrt@latest
```

Tip: Make sure that your `$PATH` contains the Go bin directory (`$HOME/go/bin` by default).

### Downloading the binary

You can download the latest binary from the [releases page](https://github.com/majermarci/hrt/releases).
For Linux and macOS, you can use the following command:

```bash
curl -L https://github.com/majermarci/hrt/releases/latest/download/hrt-linux-amd64 -o hrt

chmod +x hrt

sudo install -m 755 hrt /usr/local/bin
```

### Build and install locally

Requirements:

- Go 1.21+
- GNU Make

```bash
git clone https://github.com/majermarci/hrt.git
cd hrt/

make build

sudo install -m 755 bin/hrt /usr/local/bin
```
