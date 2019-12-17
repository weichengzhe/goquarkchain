#!/bin/bash

configPath=$1
slaveInfo=`grep -Po 'ID[" :]+\K[^"]+' $configPath | grep S`

# start slave
for value in $slaveInfo
do
	 cmd="./cluster --cluster_config "${configPath}" --service "${value}">> "${value}".log 2>&1 &"
	 echo $cmd
	 eval $cmd
done

sleep 2s

# start master
cmd="./cluster --cluster_config "${configPath}" "

if [  -n "$2" ] ;then
    echo "dsdadsadasd"
    cmd="./cluster --cluster_config "${configPath}" "$2"  "
fi

echo $cmd
eval $cmd