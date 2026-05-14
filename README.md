# toychacha
[![Go test&lint](https://github.com/blck-snwmn/toychacha/actions/workflows/test.yaml/badge.svg)](https://github.com/blck-snwmn/toychacha/actions/workflows/test.yaml)
[![CodeQL](https://github.com/blck-snwmn/toychacha/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/blck-snwmn/toychacha/actions/workflows/github-code-scanning/codeql)

Toy implementation of chacha20 poly1305 written in Go.

See: https://datatracker.ietf.org/doc/html/rfc8439

## Development

CLI tools (`golangci-lint`, `lefthook`) are managed by [aqua](https://aquaproj.github.io/) with versions pinned in [aqua.yaml](aqua.yaml).

### Install tools

Install aqua itself first (see the [aqua installation guide](https://aquaproj.github.io/docs/install)), then install the pinned tools:

```
aqua install
```

### Set up git hooks

[lefthook](lefthook.yml) runs `golangci-lint` on staged `*.go` files before each commit. Register the hooks once after cloning:

```
lefthook install
```

### Lint

```
golangci-lint run --enable=gosec
```

## Test

```
go test
```

### Benchmark

```
go test -bench . -benchmem
```

### Coverage

```
go test -v -coverpkg=. ./...
```

## WASI

### Build

```
tinygo build -o wasm.wasm -target wasi --no-debug ./cmd/main.go
```

### Run

```bash
wasmtime wasm.wasm "Ladies and Gentlemen of the class of '99: If I could offer you only one tip for the future, sunscreen would be it."
```
