
# orderbook-raft

Proof-of-concept order book server with replicated state using https://github.com/hashicorp/raft

## Running the server locally

### Single instance

```
mkdir node_1_data
SERVER_PORT=8081 RAFT_NODE_ID=node1 RAFT_ADDR=127.0.0.1 RAFT_PORT=1111 RAFT_VOL_DIR=node_1_data go run ./main.go
```

## Docker

`docker-compose.yml` contains a configuration of 10 nodes.

```
docker-compose up
```

Set up the first node as the leader and the remaining 9 nodes as followers by running

```
sh ./_scripts/join.sh
```

## REST API

`http://127.0.0.1:8081` is the base URL of the leader node


### Limit sell order

```
curl -X "POST" "http://127.0.0.1:8081/limit" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "side": "sell",
  "price": 100,
  "quantity": 10
}'
```

### Limit buy order

```
curl -X "POST" "http://127.0.0.1:8081/limit" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "side": "buy",
  "price": 100,
  "quantity": 10
}'
```

### Cancel order

```
curl -X "POST" "http://127.0.0.1:8081/cancel" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "orderId": "d3ec6544-e6ea-4849-9f5b-2ec4458c7c8e"
}'
```

### Market info

```
curl "http://127.0.0.1:8081/depth"
```

## Benchmark

Run
```
sh ./_scripts/benchmark.sh
```


```
================================
No latency
================================
Running 1m test @ http://127.0.0.1:8081/limit
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    15.73ms   11.14ms 241.27ms   95.12%
    Req/Sec   336.90     65.14   480.00     76.88%
  40230 requests in 1.00m, 14.63MB read
Requests/sec:    669.62
Transfer/sec:    249.38KB

================================
Latency 5ms
================================
Running 1m test @ http://127.0.0.1:8081/limit
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    25.85ms   11.57ms 201.51ms   97.58%
    Req/Sec   201.75     26.03   252.00     80.76%
  23989 requests in 1.00m, 8.74MB read
Requests/sec:    399.77
Transfer/sec:    149.09KB

================================
Latency 100ms
================================
Running 1m test @ http://127.0.0.1:8081/limit
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   379.43ms   48.08ms 507.35ms   58.85%
    Req/Sec    14.12      9.58    40.00     67.45%
  1577 requests in 1.00m, 589.64KB read
Requests/sec:     26.25
Transfer/sec:      9.81KB
```
