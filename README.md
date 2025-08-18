# GO - Data Structure PacKaGe

Go (generic) data structures and graph algorithms with a focus on clarity and performance.

- Go 1.21+
- Module: `github.com/elecbug/go-dspkg`

## Install

```powershell
go get github.com/elecbug/go-dspkg@latest
```

## Packages

- `bimap`: Bidirectional map with O(1) lookups key->value and value->key.
- `slice`: Generic helpers: binary search, stable merge sort, parallel sort, and `IsSorted`.
- `network-graph`: Unweighted network analysis library.

## Development

- Run tests

```powershell
go test ./...
```

## License

MIT © 2025 elecbug. See `LICENSE`.

## Credits

This project reimplements common network algorithms in Go with results validated against NetworkX.
NetworkX is © the NetworkX Developers and distributed under the BSD 3-Clause License.