#!/bin/sh
trap 'kill $(jobs -p)' TERM

CONFIG=$GODIR/chord-dfs/config/config.dfs.json ./dfs &
CONFIG=$GODIR/chord-dfs/config/config.fd.json ./member &
if [ $ALL = 1 ]
then
    sleep 2
    CONFIG=$GODIR/chord-dfs/config/config.raft.introducer.json ./raft &
fi

wait
