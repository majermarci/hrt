# hrt

A basic HTTP(S) Request Tool for CLI.

The app is still a very early prototype, but the idea is to have a simple alternative to Insomnia / httpie and so on...
This app works with a simple and easy-to-back-up yaml configuration, which can be customized per project directory also.

## Installing

If you have Go installed, you can install the app with the following command:

```bash
go install github.com/majermarci/hrt/cmd/hrt@latest
```

### Build locally

```bash
git clone https://github.com/majermarci/hrt.git
cd hrt/

go build ./cmd/hrt

sudo install -m 755 hrt /usr/local/bin
```

## To-Do / Plans

- ~~Make response outputs nicer and more readable~~
- ~~Hide response body if none is given back~~
- ~~Add timeout options for each request~~
- ~~Add option to list / hide response headers~~
- ~~Create show details option for specific request~~
- ~~Add basic and bearer token auth options~~
- Add OAuth2 auth option
- ~~Add option to call with specific certificate~~
- ~~Make an example config output, and offer to create it~~
- Better support for HEAD and OPTIONS methods
- ~~Default config search in `$HOME/.config/hrt/config.yaml`, but prioritize local config~~
- Create listing for all available requests with no option given
- Autocomplete the request names from default config
- Interactive TUI with dynamic output and selection (?)
