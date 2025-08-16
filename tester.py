#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import json
from pathlib import Path

import networkx as nx


def _try_parse_3lines(content: str):
    lines = [ln.strip() for ln in content.splitlines() if ln.strip()]
    if len(lines) < 3:
        raise ValueError("3-line format not detected")
    nodes_map = json.loads(lines[0])
    adj_map = json.loads(lines[1])
    is_bidirectional = json.loads(lines[2])
    return nodes_map, adj_map, bool(is_bidirectional)

def _try_parse_flex(content: str):
    """
    - [nodes_map, adj_map, true]
    - {"nodes": {...}, "adj": {...}, "is_bidirectional": true}
    - {"V": {...}, "E": {...}, "undirected": true}
    """
    obj = json.loads(content)
    if isinstance(obj, list) and len(obj) >= 3:
        nodes_map, adj_map, is_bidirectional = obj[0], obj[1], obj[2]
        return nodes_map, adj_map, bool(is_bidirectional)

    if isinstance(obj, dict):
        nkey = next((k for k in ["nodes", "V", "node_map"] if k in obj), None)
        akey = next((k for k in ["adj", "E", "adjacency", "edge_map"] if k in obj), None)
        bkey = next((k for k in ["is_bidirectional", "bidirectional", "undirected"] if k in obj), None)
        if nkey and akey and bkey is not None:
            return obj[nkey], obj[akey], bool(obj[bkey])

    raise ValueError("Flexible JSON format not detected")

def load_graph_file(path: str):
    content = Path(path).read_text(encoding="utf-8")
    try:
        return _try_parse_flex(content)
    except Exception:
        pass
    try:
        return _try_parse_3lines(content)
    except Exception:
        pass
    raise ValueError(
        "Unsupported input format. Supply either 3-line JSON (nodes, adj, is_bidirectional) "
        "or a single JSON array/object with those fields."
    )

def build_nx_graph(nodes_map, adj_map, is_bidirectional: bool):
    def to_id(x):
        try:
            return int(x)
        except Exception:
            return str(x)

    G = nx.Graph() if is_bidirectional else nx.DiGraph()

    nodes = [to_id(k) for k, v in nodes_map.items() if v]
    G.add_nodes_from(nodes)

    for sk, nbrs in adj_map.items():
        u = to_id(sk)
        if u not in G:
            continue
        for tk, flag in nbrs.items():
            if not flag:
                continue
            v = to_id(tk)
            if v == u or v not in G:
                continue
            G.add_edge(u, v)
    return G

def compute_metrics(G: nx.Graph, is_bidirectional: bool):
    n = G.number_of_nodes()

    metrics = {}

    all_shortest_path_length = dict(nx.all_pairs_shortest_path_length(G))
    metrics["average_shortest_path_length"] = all_shortest_path_length
    metrics["clustering_coefficient"] = nx.clustering(G)
    metrics["betweenness_centrality"] = nx.betweenness_centrality(G, normalized=True, endpoints=False)
    
    return metrics

def to_jsonable(d):
    if isinstance(d, dict):
        return {str(k): to_jsonable(v) for k, v in d.items()}
    if isinstance(d, (list, tuple)):
        return [to_jsonable(x) for x in d]
    return d

def main():
    ap = argparse.ArgumentParser(description="Read graph JSON and compute NetworkX metrics.")
    ap.add_argument("--input", "-i", required=True, help="Input graph file path")
    ap.add_argument("--output", "-o", default="", help="Output JSON path (optional)")
    args = ap.parse_args()

    nodes_map, adj_map, is_bidirectional = load_graph_file(args.input)
    G = build_nx_graph(nodes_map, adj_map, is_bidirectional)

    metrics = compute_metrics(G, is_bidirectional)
    out = {
        "is_bidirectional": bool(is_bidirectional),
        "n_nodes": G.number_of_nodes(),
        "n_edges": G.number_of_edges(),
        "metrics": to_jsonable(metrics),
    }

    text = json.dumps(out, ensure_ascii=False, indent=2)
    if args.output:
        Path(args.output).write_text(text, encoding="utf-8")
        print(f"Saved metrics to: {args.output}")
    else:
        print(text)

if __name__ == "__main__":
    main()
