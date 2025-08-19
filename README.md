# Netkit

Go (generic) graph algorithms and extensible libraries focused on clarity and performance.

- Go 1.21+
- Module: `github.com/elecbug/netkit`

## Install

```powershell
go get github.com/elecbug/netkit@latest
```

## Packages

### Graph algorithm

- [`network-graph`](./network-graph/): Unweighted network analysis library.
  - [`graph`](./network-graph/graph/): Library for creating and building graphs.
  - [`algorithm`](./network-graph/algorithm/): Library containing various graph algorithms.

# Extensible

- [`bimap`](./bimap/): Bidirectional map with O(1) lookups key->value and value->key.
- [`slice`](./slice/): Generic helpers: binary search, stable merge sort, parallel sort, and `IsSorted`.

## Development

- Run tests

```powershell
go test ./...
```

## License

MIT © 2025 elecbug. See [`LICENSE`](./LICENSE).

## Credits

This project reimplements common network algorithms in Go with results validated against NetworkX.
NetworkX is © the NetworkX Developers and distributed under the BSD 3-Clause License.