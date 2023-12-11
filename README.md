# hrt

A basic HTTP(S) Request Tool for CLI.

The app is still a very early prototype, but the idea is to have a simple alternative to Insomnia / httpie and so on...
This app works with a simple and easy-to-back-up yaml configuration, which can be customized per project directory also.

## To-Do / Plans

- Make response outputs nicer and more readable
- Hide response body if none is given back
- Add option to list / hide response headers
- Create listing for all available requests with no option given
- Create show details option for specific request
- Add Basic, OAuth2 and other token auth options
- Add timeout options for each request
- Add option to call with certificate
- Make an example config output, and offer to create it
- Autocomplete the request names from default config
- Interactive TUI with dynamic output and selection
