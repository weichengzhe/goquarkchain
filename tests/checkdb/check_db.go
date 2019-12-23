package main

import (
	"encoding/json"
	"fmt"
	"github.com/QuarkChain/goquarkchain/tests/loadtest/deployer/deploy"
	"io/ioutil"
	"os"
	"time"
)
var(
	configPath = "./checkDBConfig.json"
)
type LocalConfig struct {
	IP          string `json:"IP"`
	Port        uint64 `json:"Port"`
	User        string `json:"User"`
	Password    string `json:"Password"`
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadConfig(filePth string) *LocalConfig {
	var config LocalConfig
	f, err := os.Open(filePth)
	CheckErr(err)

	buffer, err := ioutil.ReadAll(f)
	CheckErr(err)
	err = json.Unmarshal(buffer, &config)
	CheckErr(err)
	return &config
}

func main()  {
	config := LoadConfig(configPath)
	fmt.Println("config",config)

	session := deploy.NewSSHConnect(config.User, config.Password, config.IP, int(config.Port))

	session.RunCmdIgnoreErr("docker stop $(docker ps -a|grep checkdb |awk '{print $1}')")
	session.RunCmdIgnoreErr("docker  rm $(docker ps -a|grep checkdb |awk '{print $1}')")
	session.RunCmd("docker run -itd --name checkdb --network=host quarkchaindocker/goquarkchain")

	session.RunCmd("docker exec -it checkdb  /bin/bash -c  'rm finish.txt && rm data.tar.gz  '")
	fmt.Println("50---")
	session.RunCmdIgnoreErr("docker exec -it checkdb  /bin/bash -c  'curl https://qkcmainnet-go.s3.amazonaws.com/data/2019-12-22.21:04:06.tar.gz --output data.tar.gz '")
	fmt.Println("51---")
	session.RunCmd("docker exec -it checkdb  /bin/bash -c  ' tar xvfz data.tar.gz &&  rm -rf /tmp/mainnet && mv mainnet /tmp && echo ok > finish.txt'")
	fmt.Println("52---")
	for true{
		status:=session.RunCmdAndGetOutPut("docker exec -it checkdb  /bin/bash -c  'cat finish.txt '")
		fmt.Println("status",status)
		if status=="ok"{
			fmt.Println("OOOOOOOOOOOOOO")
		}else{
			fmt.Println("RRRRRRRRRRRR")
		}
		time.Sleep(5*time.Second)
	}
	session.RunCmd("docker exec -it checkdb  /bin/bash -c  'mkdir -p /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data   '")
	session.RunCmd("docker exec -it checkdb  /bin/bash -c  'mv /tmp/mainnet /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/   '")

	session.RunCmd("docker exec -it checkdb  /bin/bash -c  'cd /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster && rm -rf *.log'")
	session.RunCmd("docker exec -it checkdb  /bin/bash -c  ' cd /go/src/github.com/QuarkChain/goquarkchain/cmd/cluster && pwd && go build && chmod +x ./run_cluster.sh && ./run_cluster.sh  ../../mainnet/singularity/cluster_config_template.json --check_db   '")
}
