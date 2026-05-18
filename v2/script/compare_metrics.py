#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import json
import math
from pathlib import Path
from typing import Any, Dict, Iterable, List, Optional, Tuple

import networkx as nx

MAP_TYPE = "map"
EDGE_MAP_TYPE = "edge_map"
SINGLE_VALUE_TYPE = "single"
CLUSTERING_TYPE = "clustering"
SHORTEST_PATHS_TYPE = "shortest_paths"

# Keep these names aligned with the Go Analyzer JSON output.
METRIC_TYPES = {
    "betweenness_centrality": MAP_TYPE,
    "closeness_centrality": MAP_TYPE,
    "clustering_coefficient": CLUSTERING_TYPE,
    "degree_assortativity_coefficient": SINGLE_VALUE_TYPE,
    "degree_centrality": MAP_TYPE,
    "diameter": SINGLE_VALUE_TYPE,
    "diameter_weight": SINGLE_VALUE_TYPE,
    "edge_betweenness_centrality": EDGE_MAP_TYPE,
    "eigenvector_centrality": MAP_TYPE,
    "modularity": SINGLE_VALUE_TYPE,
    "page_rank": MAP_TYPE,
    "shortest_paths": SHORTEST_PATHS_TYPE,
}


def node_sort_key(x: Any) -> Tuple[int, Any]:
    s = str(x)
    try:
        return (0, int(s))
    except Exception:
        return (1, s)


def to_node_id(x: Any) -> str:
    return str(x)


def parse_bool(v: Any, default: bool = False) -> bool:
    if v is None:
        return default
    if isinstance(v, bool):
        return v
    if isinstance(v, str):
        return v.strip().lower() in {"1", "t", "true", "yes", "y"}
    return bool(v)


def try_parse_latest(content: str) -> Tuple[Dict[str, Dict[str, Any]], bool, bool]:
    """
    Latest graph format:

      {
        "nodes": {
          "0": {"1": 1, "2": 1},
          "1": {"0": 1}
        },
        "directed": false,
        "weighted": true
      }

    In this format, "nodes" is actually an adjacency map.
    """
    obj = json.loads(content)
    if not isinstance(obj, dict):
        raise ValueError("latest format is not a JSON object")

    if "nodes" not in obj or not isinstance(obj["nodes"], dict):
        raise ValueError("latest format requires object field 'nodes' as adjacency map")

    if "directed" not in obj and "weighted" not in obj:
        raise ValueError("latest format requires at least 'directed' or 'weighted' marker")

    adj = obj["nodes"]
    directed = parse_bool(obj.get("directed", False))
    weighted = parse_bool(obj.get("weighted", False))
    return adj, directed, weighted


def try_parse_legacy_3lines(content: str) -> Tuple[Dict[str, Any], Dict[str, Dict[str, Any]], bool, bool]:
    """
    Legacy 3-line format:
      line 1: node existence map
      line 2: adjacency map
      line 3: undirected / bidirectional bool
    """
    lines = [ln.strip() for ln in content.splitlines() if ln.strip()]
    if len(lines) < 3:
        raise ValueError("legacy 3-line format not detected")

    nodes_map = json.loads(lines[0])
    adj_map = json.loads(lines[1])
    undirected = parse_bool(json.loads(lines[2]), False)
    return nodes_map, adj_map, not undirected, False


def try_parse_legacy_flex(content: str) -> Tuple[Dict[str, Any], Dict[str, Dict[str, Any]], bool, bool]:
    """
    Older flexible formats:
      [nodes_map, adj_map, true]
      {"nodes": {...}, "adj": {...}, "is_bidirectional": true}
      {"V": {...}, "E": {...}, "undirected": true}
    """
    obj = json.loads(content)

    if isinstance(obj, list) and len(obj) >= 3:
        nodes_map, adj_map, undirected = obj[0], obj[1], obj[2]
        return nodes_map, adj_map, not parse_bool(undirected), False

    if isinstance(obj, dict):
        nkey = next((k for k in ["V", "node_map"] if k in obj), None)
        akey = next((k for k in ["adj", "E", "adjacency", "edge_map"] if k in obj), None)
        bkey = next((k for k in ["is_bidirectional", "bidirectional", "undirected"] if k in obj), None)

        # Avoid treating latest {"nodes": adjacency, "directed": ...} as this format.
        if nkey and akey and bkey is not None:
            return obj[nkey], obj[akey], not parse_bool(obj[bkey]), False

        if "nodes" in obj and akey and bkey is not None:
            return obj["nodes"], obj[akey], not parse_bool(obj[bkey]), False

    raise ValueError("legacy flexible format not detected")


def load_graph_file(path: str) -> Tuple[nx.Graph, bool, bool]:
    content = Path(path).read_text(encoding="utf-8")

    # 1) Latest format first.
    try:
        adj, directed, weighted = try_parse_latest(content)
        G = build_graph_from_latest_adj(adj, directed=directed, weighted=weighted)
        return G, directed, weighted
    except Exception:
        pass

    # 2) Legacy flexible object / array.
    try:
        nodes_map, adj_map, directed, weighted = try_parse_legacy_flex(content)
        G = build_graph_from_legacy(nodes_map, adj_map, directed=directed)
        return G, directed, weighted
    except Exception:
        pass

    # 3) Legacy 3-line.
    try:
        nodes_map, adj_map, directed, weighted = try_parse_legacy_3lines(content)
        G = build_graph_from_legacy(nodes_map, adj_map, directed=directed)
        return G, directed, weighted
    except Exception as e:
        raise ValueError(f"unsupported graph format: {e}")


def build_graph_from_latest_adj(adj: Dict[str, Dict[str, Any]], directed: bool, weighted: bool) -> nx.Graph:
    G = nx.DiGraph() if directed else nx.Graph()

    # Add both explicit source nodes and neighbor-only nodes.
    for u in adj.keys():
        G.add_node(to_node_id(u))
    for nbrs in adj.values():
        if not isinstance(nbrs, dict):
            continue
        for v in nbrs.keys():
            G.add_node(to_node_id(v))

    for u_raw, nbrs in adj.items():
        if not isinstance(nbrs, dict):
            continue
        u = to_node_id(u_raw)
        for v_raw, val in nbrs.items():
            if not val:
                continue
            v = to_node_id(v_raw)
            if u == v:
                continue

            if weighted:
                try:
                    w = float(val)
                except Exception:
                    w = 1.0
                G.add_edge(u, v, weight=w)
            else:
                G.add_edge(u, v)

    return G


def build_graph_from_legacy(nodes_map: Dict[str, Any], adj_map: Dict[str, Dict[str, Any]], directed: bool) -> nx.Graph:
    G = nx.DiGraph() if directed else nx.Graph()

    for u, exists in nodes_map.items():
        if exists:
            G.add_node(to_node_id(u))

    for u_raw, nbrs in adj_map.items():
        u = to_node_id(u_raw)
        if u not in G or not isinstance(nbrs, dict):
            continue

        for v_raw, flag in nbrs.items():
            if not flag:
                continue
            v = to_node_id(v_raw)
            if u == v or v not in G:
                continue
            G.add_edge(u, v)

    return G


def largest_component_subgraph(G: nx.Graph) -> nx.Graph:
    if G.number_of_nodes() == 0:
        return G.copy()

    if G.is_directed():
        comps = nx.weakly_connected_components(G)
    else:
        comps = nx.connected_components(G)

    comp = max(comps, key=len)
    return G.subgraph(comp).copy()


def reachable_diameter(G: nx.Graph, weight: Optional[str] = None) -> float:
    """
    Diameter over reachable node pairs only.

    This matches the practical behavior often used for disconnected experiment
    graphs: unreachable pairs are ignored rather than treated as infinity.
    """
    max_dist = 0.0

    if weight is None:
        all_lengths = nx.all_pairs_shortest_path_length(G)
    else:
        all_lengths = nx.all_pairs_dijkstra_path_length(G, weight=weight)

    for _, dist_map in all_lengths:
        for _, d in dist_map.items():
            if d > max_dist:
                max_dist = float(d)

    return max_dist


def compute_shortest_path_counts_for_keys(G: nx.Graph, keys: Iterable[str], weight: Optional[str]) -> Dict[str, int]:
    """
    Go's latest shortest_paths output is often:
      "s->t": [{}, {}, ...]
    where the array length is the number of shortest paths.
    This helper computes only the requested pair keys to avoid exploding output size.
    """
    out: Dict[str, int] = {}

    for key in keys:
        if "->" not in key:
            continue
        s, t = key.split("->", 1)
        if s not in G or t not in G:
            continue

        try:
            if weight is None:
                paths = nx.all_shortest_paths(G, s, t)
            else:
                paths = nx.all_shortest_paths(G, s, t, weight=weight)
            out[key] = sum(1 for _ in paths)
        except (nx.NetworkXNoPath, nx.NodeNotFound):
            continue

    return out


def compute_metrics(G: nx.Graph, directed: bool, weighted: bool, shortest_path_keys: Optional[Iterable[str]] = None) -> Dict[str, Any]:
    metrics: Dict[str, Any] = {}
    weight = "weight" if weighted else None

    # Go BetweennessCentrality uses cached shortest paths; use weights when the graph is weighted.
    metrics["betweenness_centrality"] = nx.betweenness_centrality(G, weight=weight, normalized=True)

    # NetworkX directed closeness is inward by default, matching the current Go Reverse=false semantics.
    metrics["closeness_centrality"] = nx.closeness_centrality(G, distance=weight)

    local_cc = nx.clustering(G, weight=None)
    avg_cc = sum(local_cc.values()) / len(local_cc) if local_cc else 0.0
    metrics["clustering_coefficient"] = {
        "average": avg_cc,
        "global": avg_cc,
        "local": local_cc,
    }

    try:
        dac = nx.degree_assortativity_coefficient(G)
        metrics["degree_assortativity_coefficient"] = 0.0 if math.isnan(dac) else dac
    except Exception:
        metrics["degree_assortativity_coefficient"] = 0.0

    metrics["degree_centrality"] = nx.degree_centrality(G)

    metrics["diameter"] = int(reachable_diameter(G, weight=None))
    metrics["diameter_weight"] = reachable_diameter(G, weight=weight) if weighted else metrics["diameter"]

    # Go EdgeBetweennessCentrality implementation is unweighted Brandes.
    metrics["edge_betweenness_centrality"] = edge_betweenness_nested(nx.edge_betweenness_centrality(G, weight=None, normalized=True), directed)

    try:
        metrics["eigenvector_centrality"] = nx.eigenvector_centrality(G, weight=None, max_iter=1000, tol=1e-06)
    except nx.PowerIterationFailedConvergence:
        metrics["eigenvector_centrality"] = nx.eigenvector_centrality_numpy(G, weight=None)

    UG = G if not G.is_directed() else nx.Graph(G)
    try:
        communities = nx.algorithms.community.greedy_modularity_communities(UG)
        metrics["modularity"] = nx.algorithms.community.modularity(UG, communities)
    except Exception:
        metrics["modularity"] = 0.0

    metrics["page_rank"] = nx.pagerank(G, alpha=0.85, weight=None)

    if shortest_path_keys:
        metrics["shortest_paths"] = compute_shortest_path_counts_for_keys(G, shortest_path_keys, weight=weight)
    else:
        metrics["shortest_paths"] = {}

    return metrics


def edge_betweenness_nested(edge_map: Dict[Any, float], directed: bool) -> Dict[str, Dict[str, float]]:
    out: Dict[str, Dict[str, float]] = {}

    for edge, val in edge_map.items():
        u, v = edge
        su, sv = str(u), str(v)

        if not directed and node_sort_key(sv) < node_sort_key(su):
            su, sv = sv, su

        out.setdefault(su, {})[sv] = float(val)

    return out


def to_jsonable(obj: Any) -> Any:
    if isinstance(obj, dict):
        return {str(k): to_jsonable(v) for k, v in obj.items()}
    if isinstance(obj, (list, tuple, set)):
        return [to_jsonable(x) for x in obj]
    if isinstance(obj, float):
        if math.isnan(obj) or math.isinf(obj):
            return 0.0
    return obj


def load_metrics_obj(path: str) -> Dict[str, Any]:
    obj = json.loads(Path(path).read_text(encoding="utf-8"))
    if isinstance(obj, dict) and "metrics" in obj and isinstance(obj["metrics"], dict):
        return obj["metrics"]
    if isinstance(obj, dict):
        return obj
    raise ValueError(f"unsupported metrics JSON structure in: {path}")


def safe_float(x: Any, default: float = 0.0) -> float:
    try:
        return float(x)
    except Exception:
        return default


def safe_rel_err(diff: float, ref: float, eps: float = 1e-15) -> float:
    return abs(diff) / max(abs(ref), eps)


def canonical_edge_key(u: Any, v: Any, directed: bool) -> str:
    """
    Returns a comparable edge key.

    For undirected graphs, both u->v and v->u are normalized to the same
    min(u,v)->max(u,v) key using node_sort_key().
    For directed graphs, the original order is preserved.
    """
    su, sv = str(u), str(v)

    if not directed and node_sort_key(sv) < node_sort_key(su):
        su, sv = sv, su

    return f"{su}->{sv}"


def parse_edge_key(key: Any) -> Optional[Tuple[str, str]]:
    """
    Parses common edge-key string formats.

    Supported examples:
      - "u->v"
      - "u,v"
      - "(u, v)"
      - "[u, v]"

    Returns None if the key does not look like an edge key.
    """
    s = str(key).strip()

    if "->" in s:
        u, v = s.split("->", 1)
        return u.strip(), v.strip()

    if s.startswith("(") and s.endswith(")"):
        s = s[1:-1].strip()
    elif s.startswith("[") and s.endswith("]"):
        s = s[1:-1].strip()

    if "," in s:
        u, v = s.split(",", 1)
        return u.strip().strip('"').strip("'"), v.strip().strip('"').strip("'")

    return None


def flatten_numeric_map(obj: Any, prefix: str = "") -> Dict[str, float]:
    """
    Converts nested numeric maps to flat "a->b->c" keys.

    This is used for node maps and generic nested maps. Edge maps should use
    flatten_edge_map() instead, because undirected edges need canonical keys.
    """
    out: Dict[str, float] = {}

    if isinstance(obj, dict):
        for k, v in obj.items():
            key = str(k) if not prefix else f"{prefix}->{k}"
            if isinstance(v, dict):
                out.update(flatten_numeric_map(v, key))
            elif isinstance(v, (int, float)):
                out[key] = float(v)
    elif isinstance(obj, (int, float)):
        out[prefix or "value"] = float(obj)

    return out


def flatten_edge_map(obj: Any, directed: bool) -> Dict[str, float]:
    """
    Flattens edge betweenness maps to comparable edge keys.

    Supports:
      - latest / Go nested form: {"u": {"v": value}}
      - flat form: {"u->v": value}
      - tuple-like string form: {"(u, v)": value}

    For undirected graphs, both directions are canonicalized so that
    "4->24" and "24->4" compare as the same edge.
    """
    out: Dict[str, float] = {}

    if not isinstance(obj, dict):
        return out

    for k, v in obj.items():
        # Nested form: {"u": {"v": score}}
        if isinstance(v, dict):
            u = str(k)
            for k2, v2 in v.items():
                if isinstance(v2, (int, float)):
                    edge_key = canonical_edge_key(u, k2, directed)
                    out[edge_key] = float(v2)
            continue

        # Flat form: {"u->v": score} or {"(u, v)": score}
        if isinstance(v, (int, float)):
            parsed = parse_edge_key(k)
            if parsed is None:
                # Keep unknown key as-is rather than dropping data.
                out[str(k)] = float(v)
                continue

            u, vv = parsed
            edge_key = canonical_edge_key(u, vv, directed)
            out[edge_key] = float(v)

    return out


def normalize_clustering(obj: Any) -> Dict[str, float]:
    """
    Latest Go format:
      {
        "average": <float>,
        "global": <float>,
        "local": {"0": <float>, ...}
      }

    Legacy format may be just {"0": <float>, ...}.
    """
    if not isinstance(obj, dict):
        return {}

    out: Dict[str, float] = {}

    if "average" in obj and isinstance(obj["average"], (int, float)):
        out["average"] = float(obj["average"])
    if "global" in obj and isinstance(obj["global"], (int, float)):
        out["global"] = float(obj["global"])

    local = obj.get("local")
    if isinstance(local, dict):
        for k, v in local.items():
            if isinstance(v, (int, float)):
                out[f"local->{k}"] = float(v)
    else:
        for k, v in obj.items():
            if isinstance(v, (int, float)):
                out[f"local->{k}"] = float(v)

    return out


def normalize_shortest_paths(obj: Any) -> Dict[str, float]:
    """
    Supports:
      - latest Go: {"s->t": [{}, {}, ...]} where len(list) is path count
      - pair length: {"s->t": 3}
      - old nested length map: {"s": {"t": 3}}
    """
    out: Dict[str, float] = {}

    if not isinstance(obj, dict):
        return out

    for k, v in obj.items():
        sk = str(k)

        if "->" in sk:
            if isinstance(v, list):
                out[sk] = float(len(v))
            elif isinstance(v, (int, float)):
                out[sk] = float(v)
            elif isinstance(v, dict) and "distance" in v:
                out[sk] = safe_float(v["distance"])
            continue

        if isinstance(v, dict):
            for k2, v2 in v.items():
                key = f"{sk}->{k2}"
                if isinstance(v2, list):
                    out[key] = float(len(v2))
                elif isinstance(v2, (int, float)):
                    out[key] = float(v2)
                elif isinstance(v2, dict) and "distance" in v2:
                    out[key] = safe_float(v2["distance"])

    return out


def normalize_metric(name: str, typ: str, obj: Any, directed: bool) -> Dict[str, float]:
    if typ == SINGLE_VALUE_TYPE:
        return {"value": safe_float(obj)}

    if typ == CLUSTERING_TYPE:
        return normalize_clustering(obj)

    if typ == SHORTEST_PATHS_TYPE:
        return normalize_shortest_paths(obj)

    if typ == EDGE_MAP_TYPE:
        return flatten_edge_map(obj, directed=directed)

    if typ == MAP_TYPE:
        return flatten_numeric_map(obj)

    raise ValueError(f"unknown metric type: {name} / {typ}")


def compare_flat_maps(ref_s: Dict[str, float], cmp_s: Dict[str, float], include_per_node: bool) -> Dict[str, Any]:
    common = sorted(set(ref_s) & set(cmp_s), key=lambda x: (len(x), x))
    miss_in_cmp = sorted(set(ref_s) - set(cmp_s), key=lambda x: (len(x), x))
    miss_in_ref = sorted(set(cmp_s) - set(ref_s), key=lambda x: (len(x), x))

    n = len(common)
    mae = mse = mape = l1 = mean_signed = 0.0
    max_abs_err = -1.0
    max_abs_err_key = None
    per_node: Dict[str, Any] = {}

    for k in common:
        r = ref_s[k]
        c = cmp_s[k]
        diff = c - r
        ad = abs(diff)

        mae += ad
        mse += diff * diff
        mape += safe_rel_err(diff, r)
        l1 += ad
        mean_signed += diff

        if ad > max_abs_err:
            max_abs_err = ad
            max_abs_err_key = k

        if include_per_node:
            per_node[k] = {
                "ref": r,
                "cmp": c,
                "abs_error": ad,
                "signed_error": diff,
                "rel_error": safe_rel_err(diff, r),
            }

    if n > 0:
        mae /= n
        mse /= n
        mape /= n
        mean_signed /= n

    rmse = math.sqrt(mse) if n > 0 else 0.0
    l2 = math.sqrt(sum((cmp_s[k] - ref_s[k]) ** 2 for k in common)) if n > 0 else 0.0

    return {
        "n_ref": len(ref_s),
        "n_cmp": len(cmp_s),
        "n_common": n,
        "missing_in_compare": miss_in_cmp,
        "missing_in_reference": miss_in_ref,
        "mae": mae,
        "rmse": rmse,
        "max_abs_error": max_abs_err if max_abs_err_key is not None else 0.0,
        "max_abs_error_key": max_abs_err_key,
        "mape": mape,
        "mean_signed_error": mean_signed,
        "l1_error": l1,
        "l2_error": l2,
        "per_node": per_node if include_per_node else None,
    }


def shortest_path_keys_from_metrics(metrics: Dict[str, Any]) -> List[str]:
    sp = metrics.get("shortest_paths")
    if not isinstance(sp, dict):
        return []

    keys: List[str] = []
    for k, v in sp.items():
        sk = str(k)
        if "->" in sk:
            keys.append(sk)
        elif isinstance(v, dict):
            for k2 in v.keys():
                keys.append(f"{sk}->{k2}")
    return keys


def compare_metrics(
    ref_metrics: Dict[str, Any],
    cmp_metrics: Dict[str, Any],
    include_per_node: bool,
    directed: bool,
) -> Dict[str, Any]:
    report: Dict[str, Any] = {"metrics_compared": [], "metrics_skipped": []}

    for name, typ in METRIC_TYPES.items():
        if name not in ref_metrics or name not in cmp_metrics:
            report["metrics_skipped"].append(name)
            continue

        ref_s = normalize_metric(name, typ, ref_metrics[name], directed=directed)
        cmp_s = normalize_metric(name, typ, cmp_metrics[name], directed=directed)

        report[name] = compare_flat_maps(ref_s, cmp_s, include_per_node)
        report["metrics_compared"].append(name)

    return report


def main() -> None:
    ap = argparse.ArgumentParser(
        description="Compute NetworkX reference metrics for netkit graph logs and compare them with Go Analyzer metrics."
    )
    ap.add_argument("--input", "-i", required=True, help="Input graph file path")
    ap.add_argument("--output", "-o", default="", help="Output JSON path for computed reference metrics")
    ap.add_argument("--compare", "-c", default="", help="Path to Go metrics JSON for comparison")
    ap.add_argument("--report", "-r", default="", help="Output JSON path for comparison report")
    ap.add_argument("--per-node", action="store_true", help="Include per-node / per-edge errors in report")
    args = ap.parse_args()

    G, directed, weighted = load_graph_file(args.input)

    cmp_metrics: Dict[str, Any] = {}
    sp_keys: List[str] = []
    if args.compare:
        cmp_metrics = load_metrics_obj(args.compare)
        sp_keys = shortest_path_keys_from_metrics(cmp_metrics)

    computed = compute_metrics(G, directed=directed, weighted=weighted, shortest_path_keys=sp_keys)

    out = {
        "directed": directed,
        "weighted": weighted,
        "is_bidirectional": not directed,
        "n_nodes": G.number_of_nodes(),
        "n_edges": G.number_of_edges(),
        "metrics": to_jsonable(computed),
    }

    text = json.dumps(out, ensure_ascii=False, indent=2)

    if args.output:
        Path(args.output).write_text(text, encoding="utf-8")
        print(f"Saved metrics to: {args.output}")
    else:
        print(text)

    if args.compare:
        report = compare_metrics(out["metrics"], cmp_metrics, args.per_node, directed=directed)
        report_text = json.dumps(report, ensure_ascii=False, indent=2)

        if args.report:
            Path(args.report).write_text(report_text, encoding="utf-8")
            print(f"Saved comparison report to: {args.report}")
        else:
            print(report_text)


if __name__ == "__main__":
    main()
