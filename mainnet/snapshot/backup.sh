#!/bin/bash
# script to sync data into s3:

set -ex

echo "[default]" > ~/.aws/credentials
echo "aws_access_key_id="AKIAWRPAX743SNQH6WNV>>~/.aws/credentials
echo "aws_secret_access_key="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>>~/.aws/credentials

echo "[default]"> ~/.aws/config
echo "output = json">>~/.aws/config
echo "region = us-east-1">>~/.aws/config

BACKUP_DIR=/home/ubuntu/backup
rm -rf $BACKUP_DIR
mkdir $BACKUP_DIR

step=60
touch /tmp/status.txt
for (( i = 0; i < 1000000; i=(i+1) )); do
	git pull
	go build
	DATA_DIR=$GOPATH/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/
	DATE=`date +%Y-%m-%d.%H:%M:%S`
	OUTPUT_FILE=$BACKUP_DIR/$DATE.tar.gz
	LATEST_FILE=$BACKUP_DIR/LATEST
	# 3 day's backup
	RETENTION=3
	cd $DATA_DIR
	cd ..
	./stop.sh

	# checkdb
	echo "begin check_db.sh"
	./check_db.sh
	echo "end check_db.sh"
	echo $DATE >> /tmp/status.txt



	echo "begin -syc"
	echo $DATE

	mkdir -p $BACKUP_DIR
	cd $DATA_DIR
	tar cvfz $OUTPUT_FILE ./mainnet/
	cd ..
	./run.sh



	cd $BACKUP_DIR
	sz=$(ls | wc -l)
	# includes `LATEST`
	retention=$((RETENTION + 1))
	if [ "$sz" -gt "$retention" ]; then
		ls -t | tail -$((sz - retention)) | xargs -I {} rm {}
	fi
	cd -

	echo $DATE > $LATEST_FILE

	echo "begin sync-----63"
	#aws s3 sync $BACKUP_DIR s3://qkcmainnet-go/data --acl public-read
	echo "end sync-----65"
	cd $DATA_DIR
	cd ..
	#./backupcn.sh

	sleep $step
done
