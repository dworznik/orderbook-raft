#!/usr/bin/env bash

DIR=$(dirname "$0")
THREADS=2
CONNS=10
DURATION=60

echo "================================"
echo "No latency"
echo "================================"

sh ./scripts/latency.sh 0
wrk -s $DIR/../_benchmark/limit.lua -t$THREADS -c$CONNS -d${DURATION}s http://127.0.0.1:8081/limit

echo

echo "================================"
echo "Latency 5ms"
echo "================================"


sh ./scripts/latency.sh 5
wrk -s $DIR/../_benchmark/limit.lua -t$THREADS -c$CONNS -d${DURATION}s http://127.0.0.1:8081/limit

echo
echo "================================"
echo "Latency 100ms"
echo "================================"


sh ./scripts/latency.sh 100
wrk -s $DIR/../_benchmark/limit.lua -t$THREADS -c$CONNS -d${DURATION}s http://127.0.0.1:8081/limit
