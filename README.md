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
