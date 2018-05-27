#!/usr/bin/env bash

wrk -t2 -c256  -d20s -T5 --script=./benchmark/wrk.lua  --latency http://127.0.0.1:8087/invoke
wrk -t2 -c64  -d60s -T5 --script=./benchmark/wrk.lua  --latency http://127.0.0.1:8087/invoke
wrk -t2 -c128 -d60s -T5 --script=./benchmark/wrk.lua  --latency http://127.0.0.1:8087/invoke
wrk -t2 -c256 -d60s -T5 --script=./benchmark/wrk.lua  --latency http://127.0.0.1:8087/invoke

