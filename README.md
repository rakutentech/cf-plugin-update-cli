# cf-plugin-update-cli [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: /LICENSE

`cf-plugin-update-cli` is a [cloudfoundry/cli](https://github.com/cloudfoundry/cli) plugin. It allows you to update cli to the latest version. 

## Demo

The following domo shows updating from `v6.14.0` to `v6.16.1` (the latest ver at 201603).

![gif](/doc/update-cli.gif)

## Install

To install this plugin, use `go get` (make sure you have already setup golang enviroment like `$GOPATH`),

```bash
$ go get -d github.com/tcnksm/cf-plugin-update-cli
$ cd $GOPATH/src/github.com/tcnksm/cf-plugin-update-cli
$ make install # if you have already installed, then run `make uninstall` before
```

Or you can install it from [my experimental plugin repository](https://t-plugins.au-syd.mybluemix.net/ui/).

```bash
$ cf add-plugin-repo tcnksm https://t-plugins.au-syd.mybluemix.net
$ cf install-plugin -r tcnksm update-cli
```

Since this plugin is still immature, it's not on [Community Plugin Repo](http://plugins.cloudfoundry.org/ui/). 

## Usage

To update, run the following command,

```bash
$ cf update-cli
```

It will be atumatically detect your OS/Arch and installs appropriate binary and replace with old one. If permission denied to write new binary, then use `sudo`. 

## Contribution

1. Fork ([https://github.com/tcnksm/cf-plugin-update-cli/fork](https://github.com/tcnksm/cf-plugin-update-cli/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[Taichi Nakashima](https://github.com/tcnksm)
