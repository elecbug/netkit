# Netkit

Go (generic) graph algorithms and extensible libraries focused on clarity and performance.

- Go 1.21+
- Module: `github.com/elecbug/netkit`

## Install

```powershell
go get github.com/elecbug/netkit@latest
```

## Packages

### Graph

- [`graph`](./network-graph/graph/): Library for creating and building unweighted graphs.
  - [`standard_graph`](./network-graph/graph/standard_graph/): Library for generating standard graphs like Erdos-Reyni graph.
  - [`algorithm`](./network-graph/algorithm/): Library containing various graph algorithms.

### P2P

- [`p2p`](./p2p/): Library that integrates with graph libraries to form networks and enable p2p broadcast experiments.

### Extensible

- [`bimap`](./bimap/): Bidirectional map with O(1) lookups key->value and value->key.

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