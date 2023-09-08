#!/usr/bin/env bash

for n in {2..10}; do
id=$(printf "%02d" $n)
curl --location --request POST 'localhost:8081/raft/join' --header 'Content-Type: application/json' --data-raw "{ \"node_id\": \"node-${id}\", \"raft_address\": \"node-${id}:1111\" }"
echo
echo
sleep 1
done
