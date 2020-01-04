# `faas-akash provider`

<img alt="" src="https://camo.githubusercontent.com/cf01eefb5b6905f3774376d6d1ed55b8f052d211/68747470733a2f2f626c6f672e616c6578656c6c69732e696f2f636f6e74656e742f696d616765732f323031372f30382f666161735f736964652e706e67" width="600px">

[![version](https://img.shields.io/github/tag/vitwit/faas-akash.svg)](https://github.com/vitwit/faas-akash/releases/latest)
[![Build Status](https://circleci.com/gh/vitwit/faas-akash.svg?style=svg)](https://circleci.com/gh/vitwit/faas-akash)
[![codecov](https://codecov.io/gh/vitwit/faas-akash/branch/master/graph/badge.svg)](https://codecov.io/gh/vitwit/faas-akash)
[![Go Report Card](https://goreportcard.com/badge/github.com/openfaas/faas)](https://goreportcard.com/report/github.com/openfaas/faas)[![GoDoc](https://godoc.org/github.com/vitwit/go-faas?status.svg)](https://godoc.org/github.com/vitwit/go-faas) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


## Contents

* [Requirements](#requirements)
* [Installation](#installation)
* [Getting Started](#getting-started)

Requirements
------------

- [Go](https://golang.org/doc/install) 1.13+ (to build the provider plugin)
- [faas-cli](https://github.com/openfaas/faas-cli/releases) 
- [akash-cli](https://github.com/ovrclk/akash/releases) (make sure to put akash-cli in your $PATH or /usr/local/bin)

Installation
------------

* Clone this repository and cd into the directory
* Run `make build`, it will generate a file named `faas-akash`  
* run `./faas-akash` and enjoy OpenFaas :handshaking: Akash Network 

Getting Started
---------------
A Makefile has been provided to make it easier to work with faas-akash provider.

* Clone this repo and cd into it
* run `make build`, this will download all the dependencies and build the project into a binary named `faas-akash`.
* an example `config.yaml.example` file has been provided, rename this file to config.yaml and set the correct values for 
Port number, Read and Write timeouts
* run `./faas-akash` and faas-akash should be up and running
* open another terminal and export the faas-akash url as
`export OPENFAAS_URL=127.0.0.1:<port-from-config.yaml>`
* now faas-cli should be able to interact with faas-akash provider


Provider Configuration
------------

* by default, faas-akash provider assumes a `config.yaml` to be present in $HOME/.faas-akash/
* few fields can be set using this config file or environmental variables

| Property                      | Description                                                                                                           | Required    |
| ----------------------------- | ------------------------------------------------- | ---------- |
| `port`                        | OpenFaas Gateway Port                             | `No`       |
| `readTimeout`                 | Request body ReadTimeout in seconds               | `No`       |
| `writeTimeout`                | Request Body WriteTimeout in seconds              | `No`       |

### Example

```yaml
port: 8090
readTimeout: 180
writeTimeout: 180
```

Contribution Guidelines
-----------------------
 
To contribute, please read the contribution guidelines: [Contributing to faas-akash Provider](.github/CONTRIBUTING.md)

[install-go]: https://golang.org/doc/install#install
