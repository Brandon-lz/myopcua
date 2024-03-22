#!/usr/bin/env bash

cd `dirname $0`

# start opua server
pids=`ps -aux | grep opcua_simulation_server.py | grep -v grep | awk '{print $2}'`
if [ -n "$pids" ]; then
    echo "try to kill opcua_simulation_server.py: $pids"
    kill -s 9 $pids
fi

# pip install opcua==0.98.3 -i https://mirror.sjtu.edu.cn/pypi/web/simple

python opcua_simulation_server.py