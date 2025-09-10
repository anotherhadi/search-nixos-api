<div align="center">
    <img src="https://raw.githubusercontent.com/anotherhadi/search-nixos/main/static/logo.png" width="120px" />
</div>

# Search NixOS API

<p>
    <a href="https://github.com/anotherhadi/search-nixos-api/releases"><img src="https://img.shields.io/github/release/anotherhadi/search-nixos-api.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/anotherhadi/search-nixos-api?tab=doc"><img src="https://godoc.org/github.com/anotherhadi/search-nixos-api?status.svg" alt="GoDoc"></a>
    <a href="https://goreportcard.com/report/github.com/anotherhadi/search-nixos-api"><img src="https://goreportcard.com/badge/github.com/anotherhadi/search-nixos-api" alt="GoReportCard"></a>
</p>

The Search NixOS API is a service designed to facilitate the search and retrieval of options across various Nix-related projects, including Nixpkgs, NixOS, Home Manager, nix-darwin, and the Nix User Repository (NUR). This API provides developers and users with a unified interface to query and explore configuration options available within these ecosystems.

This API is used in the [Search NixOS](https://github.com/anotherhadi/search-nixos) project, which is a web application that allows users to search for NixOS options and view their documentation. The API serves as the backend for this application, providing the necessary data and functionality to support the search feature.

## Hosted Public Instance

<del>A publicly accessible instance of the Nix Options Search API is hosted at [search-nixos-api.hadi.diy](https://search-nixos-api.hadi.diy). Users can utilize this endpoint to programmatically search for options without the need to set up their own instance. (Now down)</del>

## Data Source

The API leverages the [`nix-json`](https://github.com/anotherhadi/nix-json) project to retrieve and process JSON files containing option definitions. This integration ensures that the API delivers up-to-date and comprehensive information about available options across the supported Nix projects.

## Features

- **Comprehensive Search**: Query options from Nixpkgs, NixOS, Home Manager, nix-darwin, and NUR through a single interface.
- **Public Accessibility**: Utilize the hosted instance for immediate access without any setup.
- **Up-to-date Data**: Integration with `nix-json` ensures that the API reflects the latest options available in the Nix ecosystem.

## Usage

To search for options, send a GET request to the API's endpoint with your query parameters. For example:

```
GET https://search-nixos-api.hadi.diy/search?q=your_option_name
```

The API will respond with a JSON object containing matching options and their details.

## Contributing

Contributions to the Nix Options Search API are welcome. Please refer to the project's repository for guidelines on how to contribute.

## Funding

Maintaining and hosting the Nix Options Search API incurs ongoing server costs. If you find this project useful and would like to support its development and upkeep, please consider making a donation. Contributions can be made through Ko-fi at [ko-fi.com/anotherhadi](https://ko-fi.com/anotherhadi). Your support is greatly appreciated and helps ensure the continued availability and improvement of the API.

## License

This project is licensed under the MIT License.
