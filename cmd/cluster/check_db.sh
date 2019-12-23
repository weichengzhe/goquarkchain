#!/bin/bash

step=60 #间隔的秒数，不能大于60
touch /tmp/status.txt
for (( i = 0; i < 60; i=(i+1) )); do
 rm data.tar.gz
 curl https://qkcmainnet-go.s3.amazonaws.com/data/2019-12-22.21:04:06.tar.gz --output data.tar.gz
 rm -rf /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data
 mkdir -p /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data
 tar xvfz data.tar.gz  && mv mainnet /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data
 rm -rf *.log
 git pull
 go build
 chmod +x ./run_cluster.sh && ./run_cluster.sh  ../../mainnet/singularity/cluster_config_template.json --check_db
 DATE=`date +%Y-%m-%d.%H:%M:%S`
 echo $DATE >> /tmp/status.txt
 sleep $step
done
