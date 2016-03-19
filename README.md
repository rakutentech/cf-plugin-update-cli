# cf-plugin-update-cli [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: /LICENSE

`cf-plugin-update-cli` is a [cloudfoundry/cli](https://github.com/cloudfoundry/cli) plugin. It allows you to update cli to the latest version. 

## Install

To install this plugin, use `go get` (make sure you have already setup golang enviroment like `$GOPATH`),

```bash
$ go get -d github.com/tcnksm/cf-plugin-update-cli
$ cd $GOPATH/src/github.com/tcnksm/cf-plugin-update-cli
$ make install # if you have already installed, then run `make uninstall` before
```

Since this plugin is still immature and PoC, it's not uploaded on [Community Plugin Repo](http://plugins.cloudfoundry.org/ui/). But in future, I'll add this plugin there and make it more easy to install.

## Usage

```bash
$ cf update
```

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
