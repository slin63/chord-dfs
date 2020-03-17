#!/bin/sh
trap 'kill $(jobs -p)' TERM

CONFIG=$GODIR/chord-dfs/config/config.dfs.json ./dfs &
CONFIG=$GODIR/chord-dfs/config/config.fd.json ./member &
if [ $ALL = 1 ]
then
    CONFIG=$GODIR/chord-dfs/config/config.raft.json ./raft &
fi

wait
