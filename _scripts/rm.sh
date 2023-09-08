for n in {1..10}; do
docker rm orderbook-raft-node-$(printf "%02d" $n)-1
done
