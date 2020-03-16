#!/bin/sh
trap 'kill $(jobs -p)' TERM

CONFIG=$GODIR/chord-dfs/config.json ./dfs &
CONFIG=$GODIR/chord-failure-detector/config.json ./member &
if [ $ALL = 1 ]
then
    CONFIG=$GODIR/raft-consensus/config.json ./raft &
fi

wait
