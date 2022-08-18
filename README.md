# aes

## Disclaimer

Currently only supports payload of multiples of 16 bytes, not functional
more-so for the learning/experimental purposes.

## Usage

```sh
go run . \
  --payload examplepayloadddd
```

## TODOs

- Use the `cipher.NewGCM(block)` instead using the `cipher.Block` directly.
  Thanks to that the padding and authentication are handled by the package.
- Enable support for files and directories as well as binary data.
