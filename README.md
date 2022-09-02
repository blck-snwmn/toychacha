# go-chacha

Toy implementation of chacha20 poly1305 written in Go.

See: https://datatracker.ietf.org/doc/html/rfc8439

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
