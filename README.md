# aes

## Introduction

AES-256 implementation in Go using the `crypto/aes` package from the standard
library. The logging is handled through `go.uber.org/zap` and the command line
interface was built with `github.com/spf13/cobra`.

## Installation

### Requirements

The script is idempotent and will either use (if existing in the directory) or
generate a `key` and `nonce` files containing the bytes used in the algorithm.

### Build from source

```sh
git clone https://github.com/piotrostr/aed && cd aed
go build .
./aes
```

## Usage

```sh
$ ./aes

AES is a symmetric encryption algorithm that can be used to encrypt and decrypt data.
It is used to encrypt data at rest at Google.

Usage:
  aes [flags]
  aes [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  decrypt     Decrypts a file using AES-256-GCM
  encrypt     Encrypts a file using AES-256-GCM
  help        Help about any command
  stdin       Encrypt input from stdin

Flags:
  -h, --help   help for aes

Use "aes [command] --help" for more information about a command.
```

The script generates `.enc` files after encryption and given an `.enc` file it
writes the decoded file without the suffix.

## Disclaimer

This is not production software and it is not affiliated with Google. I wrote
it for educational purposes, not as part of my employment. I have not tested
the security and protocol compliance (yet) so please be careful when dealing
with sensitive data!

## TODOs

Will add a config later so that the package can be used as a command line tool,
for now the `key` and `nonce` files are generated in the working directory of
the package. They will be specifiable in the future, also planning to allow stdin.

## License

MIT
