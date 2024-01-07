# Documentation

The documentation for HRT is split into the following sections:

- [Configuration](#configuration)
- [Flags](#command-line-flags)

## Configuration

By default the app will always use the local `hrt.yaml` as configuration if it exists.
If it doesn't exist, it will use the global configuration file at `~/.config/hrt/config.yaml`.
If that doesn't exist either, it will offer to create a local one with [default example](cmd/hrt/example_config.yaml) values.

There is also an option to create the global config file with [default example](cmd/hrt/example_config.yaml) values using the `-g` flag.

### Structure of configuration

The configuration file uses the YAML format and supports the following fields where `endpoint` can be any string:

- `endpoint`: The endpoint name to send the request to. Repeating field that can be used multiple times.
- `endpoint.url`: The URL of the request. Format: `scheme://host:port/path?query#fragment`
- `endpoint.method`: The HTTP method for the request. Default is `GET`. Supported methods are `GET`, `POST`, `PUT`, `PATCH` and `DELETE`. `HEAD` and `OPTIONS` work too, and will display response headers by default.
- `endpoint.headers`: The headers for the request. It can contain as many key-value pairs as you want.
- `endpoint.body`: The body of the request.
- `endpoint.basic_auth.username`: The username for basic authentication.
- `endpoint.basic_auth.password`: The password for basic authentication.
- `endpoint.bearer_token`: The bearer token for the request.

### Example configuration

All options are represented below, but the only mandatory field is `url`.
Please note that the two authentication options cannot be used together like in this example.

```yaml
<request_name>:
  url: <url>
  method: <http_method>
  headers:
    <header_name>: <header_value>
  body: <body_content>
  basic_auth:
    username: <username>
    password: <password>
  bearer_token: <bearer_token>
```

## Command Line Flags

HRT supports the following command line flags:

- `-r`: The name of the request you want to run form currently used configuration file.
- `-l`: List all available requests in the currently used configuration file.
- `-c`: Use a custom configuration file instead of the default / global one.
- `-a`: Run every request in the currently used configuration file.
- `-g`: Create global configuration file with default example values.
- `-t`: Timeout for the request in seconds. Default is 10 seconds.
- `-table`: Show response in table format. Warning: This is an experimental feature. It will be ugly with large responses.
- `-k`: Skip TLS certificate verification. (Insecure)
- `-cacert`: Path to CA certificate file. Content is added to the system certificate pool.
- `-cert`: Path to TLS certificate file.
- `-key`: Path to TLS key file.
- `-h`: Show help message.
- `-version`: Show version information.
- `-v`: Show verbose output, such as the used configuration file, request name, method, endpoint and response status code.
- `-vv`: Show verbose TLS info along with both request and response headers. Also shows details for normal verbose option.
