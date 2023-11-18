#!/usr/bin/env bash
MAX="$@"
for ((i=1 ; i<=$MAX ; i++)); do
    ./bin/darwin/worker --queue test1 --wfid "schedule"$i --exit true
done
