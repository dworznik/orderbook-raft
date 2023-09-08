#!/usr/bin/env bash

for n in {1..10}; do
docker exec orderbook-raft-node-$(printf "%02d" $n)-1 tc qdisc del dev eth0 root netem 2> /dev/null
docker exec orderbook-raft-node-$(printf "%02d" $n)-1 tc qdisc add dev eth0 root netem delay ${1}ms
done
