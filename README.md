# cf-plugin-update-cli [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: /LICENSE

`cf-plugin-update-cli` is a [cloudfoundry/cli](https://github.com/cloudfoundry/cli) plugin. It allows you to update cli to the latest version. 

## Demo

The following domo shows updating from `v6.14.0` to `v6.16.1` (the latest ver at 201603).

![gif](/doc/update-cli.gif)

## Install

To install this plugin, use `cf` command. It's hosted on [Community Plugin Repo](http://plugins.cloudfoundry.org/ui/). 

```bash
$ cf install-plugin -r CF-Community "update-cli"
```

## Usage

To update, run the following command,

```bash
$ cf update-cli
```

It will be atumatically detect your OS/Arch and installs appropriate binary and replace with old one. If permission denied to write new binary, then use `sudo`. 

## Author

[Taichi Nakashima](https://github.com/tcnksm)
