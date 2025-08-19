#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import json
from pathlib import Path
from typing import Any, Dict, Tuple

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
    Accept flexible formats:
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


def compute_metrics(G: nx.Graph, is_bidirectional: bool) -> Dict[str, Any]:
    metrics: Dict[str, Any] = {}

    # NOTE: keep key names aligned with your Go outputs for easy comparison
    metrics["betweenness_centrality"] = nx.betweenness_centrality(G)
    metrics["closeness_centrality"] = nx.closeness_centrality(G)  # wf_improved=True semantics in recent NX
    metrics["clustering_coefficient"] = nx.clustering(G)          # for DiGraph, NX uses underlying undirected
    metrics["degree_centrality"] = nx.degree_centrality(G)
    metrics["edge_betweenness_centrality"] = nx.edge_betweenness_centrality(G)
    metrics["eigenvector_centrality"] = nx.eigenvector_centrality(G)
    # print(metrics["edge_betweenness_centrality"])  # debug output
    metrics["page_rank"] = nx.pagerank(G, weight=None)            # unweighted
    metrics["shortest_paths"] = dict(nx.all_pairs_shortest_path_length(G))

    # If you also want degree_centrality comparison, uncomment:

    return metrics


def to_jsonable(d):
    if isinstance(d, dict):
        return {str(k): to_jsonable(v) for k, v in d.items()}
    if isinstance(d, (list, tuple)):
        return [to_jsonable(x) for x in d]
    return d


# -------- comparison helpers --------

NUMERIC_NODE_METRICS = {
    "betweenness_centrality",
    "closeness_centrality",
    "clustering_coefficient",
    "degree_centrality",
    "edge_betweenness_centrality",
    "eigenvector_centrality",
    "page_rank",
    "shortest_paths"
    # "degree_centrality",
}

def _load_metrics_obj(path: str) -> Dict[str, Any]:
    """
    Accepts either:
      - {"metrics": {...}}  (preferred)
      - {"betweenness_centrality": {...}, "closeness_centrality": {...}, ...}
    Returns the metrics dict.
    """
    obj = json.loads(Path(path).read_text(encoding="utf-8"))
    if isinstance(obj, dict) and "metrics" in obj and isinstance(obj["metrics"], dict):
        return obj["metrics"]
    if isinstance(obj, dict):
        # assume direct metrics at root
        return obj
    raise ValueError(f"Unsupported metrics JSON structure in: {path}")


def _safe_rel_err(diff: float, ref: float, eps: float = 1e-15) -> float:
    denom = abs(ref)
    if denom < eps:
        denom = eps
    return abs(diff) / denom


def compare_metric_maps(name: str, ref: Dict[str, float], cmp_: Dict[str, float], include_per_node: bool) -> Dict[str, Any]:
    if name == "edge_betweenness_centrality":
        ref_s = {str(k): float(v) for k, v in ref.items()}
        cmp_s = {f"({str(k1)}, {str(k2)})": float(v) for k1, v1 in cmp_.items() for k2, v in v1.items()}
    elif name == "shortest_paths":
        ref_s = {f"({str(k1)}, {str(k2)})": float(v) for k1, v1 in ref.items() for k2, v in v1.items()}
        cmp_s = {f"({str(k1)}, {str(k2)})": float(v) for k1, v1 in cmp_.items() for k2, v in v1.items()}
    else:
        # align keys as strings
        ref_s = {str(k): float(v) for k, v in ref.items()}
        cmp_s = {str(k): float(v) for k, v in cmp_.items()}


    common = sorted(set(ref_s.keys()) & set(cmp_s.keys()), key=lambda x: (len(x), x))
    miss_in_cmp = sorted(set(ref_s.keys()) - set(cmp_s.keys()))
    miss_in_ref = sorted(set(cmp_s.keys()) - set(ref_s.keys()))

    n = len(common)
    mae = rmse = mape = mse = l1 = l2 = 0.0
    max_abs_err = -1.0
    max_abs_err_node = None
    mean_signed = 0.0

    per_node = {}

    for k in common:
        r = ref_s[k]
        c = cmp_s[k]
        diff = c - r
        ad = abs(diff)
        l1 += ad
        mse += diff * diff
        mean_signed += diff
        if ad > max_abs_err:
            max_abs_err = ad
            max_abs_err_node = k
        mae += ad
        mape += _safe_rel_err(diff, r)

        if include_per_node:
            per_node[k] = {"ref": r, "cmp": c, "abs_error": ad, "signed_error": diff}

    if n > 0:
        mae /= n
        mape /= n
        mse /= n
        rmse = mse ** 0.5
        l2 = (sum((cmp_s[k] - ref_s[k]) ** 2 for k in common)) ** 0.5
        mean_signed /= n

    return {
        "n_ref": len(ref_s),
        "n_cmp": len(cmp_s),
        "n_common": n,
        "missing_in_compare": miss_in_cmp,
        "missing_in_reference": miss_in_ref,
        "mae": mae,
        "rmse": rmse,
        "max_abs_error": max_abs_err if max_abs_err_node is not None else 0.0,
        "max_abs_error_node": max_abs_err_node,
        "mape": mape,  # mean absolute percentage error (safe denom)
        "mean_signed_error": mean_signed,
        "l1_error": l1,
        "l2_error": l2,
        "per_node": per_node if include_per_node else None,
    }


def compare_shortest_path_length(ref: Dict[str, Dict[str, int]], cmp_: Dict[str, Dict[str, int]]) -> Dict[str, Any]:
    result: Dict[str, Any] = {}

    for s, v in ref.items():
        for e, _ in v.items():
            if s not in cmp_ or e not in cmp_[s]:
                result[s] = {"ref": ref[s][e], "cmp": None}
            elif ref[s][e] != cmp_[s][e]:
                result[s] = {"ref": ref[s][e], "cmp": cmp_[s][e]}
    return result


def compare_metrics(ref_metrics: Dict[str, Any], cmp_metrics: Dict[str, Any], include_per_node: bool) -> Dict[str, Any]:
    report: Dict[str, Any] = {"metrics_compared": []}

    for name in sorted(NUMERIC_NODE_METRICS):
        if name in ref_metrics and name in cmp_metrics:
            # Only compare node->float maps
            if isinstance(ref_metrics[name], dict) and isinstance(cmp_metrics[name], dict):
                report[name] = compare_metric_maps(name, ref_metrics[name], cmp_metrics[name], include_per_node)
                report["metrics_compared"].append(name)
        # else: silently skip if either missing

    return report


def main():
    ap = argparse.ArgumentParser(description="Read graph JSON, compute NetworkX metrics, and optionally compare to another JSON.")
    ap.add_argument("--input", "-i", required=True, help="Input graph file path")
    ap.add_argument("--output", "-o", default="", help="Output JSON path for computed metrics (optional)")
    ap.add_argument("--compare", "-c", default="", help="Path to comparison metrics JSON (optional)")
    ap.add_argument("--report", "-r", default="", help="Where to save the comparison report JSON (optional)")
    ap.add_argument("--per-node", action="store_true", help="Include per-node errors in the comparison report")
    args = ap.parse_args()

    nodes_map, adj_map, is_bidirectional = load_graph_file(args.input)
    G = build_nx_graph(nodes_map, adj_map, is_bidirectional)

    computed = compute_metrics(G, is_bidirectional)
    out = {
        "is_bidirectional": bool(is_bidirectional),
        "n_nodes": G.number_of_nodes(),
        "n_edges": G.number_of_edges(),
        "metrics": to_jsonable(computed),
    }

    # Save or print computed metrics
    text = json.dumps(out, ensure_ascii=False, indent=2)
    if args.output:
        Path(args.output).write_text(text, encoding="utf-8")
        print(f"Saved metrics to: {args.output}")
    else:
        print(text)

    # Optional: comparison
    if args.compare:
        try:
            cmp_metrics_raw = _load_metrics_obj(args.compare)
        except Exception as e:
            raise SystemExit(f"[compare] Failed to load comparison metrics: {e}")

        report = compare_metrics(out["metrics"], cmp_metrics_raw, args.per_node)
        report_text = json.dumps(report, ensure_ascii=False, indent=2)

        if args.report:
            Path(args.report).write_text(report_text, encoding="utf-8")
            print(f"Saved comparison report to: {args.report}")
        else:
            print(report_text)


if __name__ == "__main__":
    main()
