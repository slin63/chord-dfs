#!/bin/sh
trap 'kill $(jobs -p)' TERM

./dfs &
./member &

wait
