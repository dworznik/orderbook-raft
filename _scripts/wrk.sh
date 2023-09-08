#!/usr/bin/env bash

wrk -s $DIR/../_benchmark/limit.lua -t2 -c10 -d60s http://127.0.0.1:8081/limit
