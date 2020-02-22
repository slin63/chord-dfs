#!/bin/sh
trap 'kill $(jobs -p)' TERM

# start service in background here
go run src/main.go &
go run chord-failure-detector/src/main.go &

wait
