#!/bin/bash

cpu=$1

process="qch-run-$cpu.sh"

#ps ax | grep "$process"

#ps ax | grep "$process" > /dev/null
#if [ $? -eq 0 ]; then
#    echo "Process $process is running."
#else
#    echo "Run $process."
#    bash $process
#fi

echo $process
pgrep -f "$process"
if [[ -n `pgrep -f "$process"` ]]; 
then
    echo "Process $process is running"
else 
    echo "Run $process"
    bash $process
fi
