# hrt - HTTP(S) Request Tool

[![Go Report Card](https://goreportcard.com/badge/github.com/majermarci/hrt)](https://goreportcard.com/report/github.com/majermarci/hrt)
[![Go Reference](https://pkg.go.dev/badge/github.com/majermarci/hrt.svg)](https://pkg.go.dev/github.com/majermarci/hrt)
![Latest release](https://img.shields.io/github/v/release/majermarci/hrt)

The app is still in early development (and it is a learning project for me), but the idea is to have a simple CLI alternative to Insomnia / httpie and similar tools.
`hrt` works with a simple and easy-to-back-up `yaml` configuration, which can also be customized per project and included in repos.

## Features

- Simple yaml configuration for organizing your request collection(s)
  - Specify headers, body and method for each request
  - Bearer Token and Basic Auth support
  - Support for local and global config files
  - Option to create a default config file
- Various option flags for running requests
  - Usage of specific TLS certificate, adding new CA chains to existing one or skipping certificate verification
  - Global timeout option
  - Verbose outputs showing request and TLS details
  - Option to run every request in currently used collection file

### Usage

For usage information, run `hrt -h` or read about options at the [documentation](usage.md) page.

## Installing with Go

If you have Go installed, you can install the app with the following command:

```bash
go install github.com/majermarci/hrt/cmd/hrt@latest
```

Tip: Make sure that your `$PATH` contains the Go bin directory (`$HOME/go/bin` by default).

### Downloading the binary

You can download the latest binary from the releases page.
For Linux and macOS, you can use the following command:

```bash
curl -L https://github.com/majermarci/hrt/releases/latest/download/hrt-linux-amd64) -o hrt

chmod +x hrt

sudo install -m 755 hrt /usr/local/bin
```

### Build and install locally

Requirements:

- Go 1.21.5+
- GNU Make

```bash
git clone https://github.com/majermarci/hrt.git
cd hrt/

make build

sudo install -m 755 bin/hrt /usr/local/bin
```

## To-Do / Plans

- ~~Make response outputs nicer and more readable~~
- ~~Hide response body if none is given back~~
- ~~Add timeout options for each request~~
- ~~Add option to list / hide response headers~~
- ~~Create show details option for specific request~~
- ~~Add basic and bearer token auth options~~
- ~~Add option to call with specific certificate~~
- ~~Make an example config output, and offer to create it~~
- ~~Default config search in `$HOME/.config/hrt/config.yaml`, but prioritize local config~~
- ~~Add option to list all available requests~~
- Add OAuth2 auth option
- Better support for HEAD and OPTIONS methods
- ~~Add version variable to the build process along with commit ID~~
- Autocomplete the request names from default config
- Interactive TUI with dynamic output and selection (?)
- Add tests...
