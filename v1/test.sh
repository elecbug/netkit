go test -v ./graph/

cd script

python3 compare_metrics.py -i ../graph/directional.graph.log -c ../graph/directional.metrics.log -r ./directional.report.log -o ./directional.out.log
python3 compare_metrics.py -i ../graph/bidirectional.graph.log -c ../graph/bidirectional.metrics.log -r ./bidirectional.report.log -o ./bidirectional.out.log

cd ..