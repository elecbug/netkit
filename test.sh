go test -v ./network-graph/

cd script

python3 compare_metrics.py -i ../network-graph/directional.graph.log -c ../network-graph/directional.metrics.log -r ./directional.report.log -o ./directional.out.log
python3 compare_metrics.py -i ../network-graph/bidirectional.graph.log -c ../network-graph/bidirectional.metrics.log -r ./bidirectional.report.log -o ./bidirectional.out.log

cd ..