# !/bin/bash
set -e

go test -v ./v2/graph/analyzer/

CURRENT_DIR=$(pwd)
SOURCE_DIR=$(dirname "$0")

cd "$SOURCE_DIR"

if [ ! -d "venv" ]; then
    python3 -m venv venv
fi

source venv/bin/activate
pip3 install networkx numpy matplotlib scipy

python3 compare_metrics.py -i ../graph/analyzer/graph.log -c ../graph/analyzer/metrics.log -r ./report.log -o ./out.log

deactivate

cd "$CURRENT_DIR"