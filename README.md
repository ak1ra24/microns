# MicroNS

[![Build Status](https://travis-ci.com/ak1ra24/microns.svg?branch=master)](https://travis-ci.com/ak1ra24/microns) [![Go Report Card](https://goreportcard.com/badge/github.com/ak1ra24/microns)](https://goreportcard.com/report/github.com/ak1ra24/microns) [![codecov](https://codecov.io/gh/ak1ra24/microns/branch/master/graph/badge.svg)](https://codecov.io/gh/ak1ra24/microns) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Docker を用いたルーティングシュミレーション

## Prerequisites

### For Example

Ubuntu18.04

- Docker Install

```
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl software-properties-common -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable test edge" -y
sudo apt-get update
sudo apt-get install docker-ce -y
```

- Go Install

```
sudo add-apt-repository ppa:longsleep/golang-backports -y
sudo apt update
sudo apt install golang-go -y
export GO111MODULE=on
```

- Graphviz & ascii graph for microns image

```
sudo add-apt-repository universe -y
sudo apt update
sudo apt install graphviz cpanminus -y
sudo cpanm Graph::Easy
```

## Usage

1. セットアップ

```
go build
./microns

■■      ■■   ■                      ■■     ■■      
 ■■     ■■                           ■■    ■       
 ■■    ■ ■                           ■ ■   ■       
 ■ ■     ■  ■■   ■■■■  ■■ ■  ■■■■    ■ ■   ■   ■■■■
 ■ ■  ■  ■   ■  ■■  ■   ■   ■■  ■■   ■  ■  ■   ■  ■
 ■ ■  ■  ■   ■  ■       ■   ■    ■   ■   ■ ■    ■  
 ■  ■    ■   ■  ■       ■   ■    ■   ■   ■■■      ■
 ■  ■■   ■   ■  ■■      ■   ■■  ■■   ■    ■■   ■  ■
■■   ■  ■■  ■■   ■■■■  ■■    ■■■■   ■■     ■   ■■■■
```

2. 使い方

```
./microns help
microns

Usage:
  microns [flags]
  microns [command]

Available Commands:
  convert     convert from tinet config file to microns config file
  create      create docker container and ns topology
  delete      delete docker container and ns topology
  help        Help about any command
  image       create network topology image file
  init        A brief description of your command
  recreate    reconfigure router
  status      status docker container and ns topology
  test        Execute test from config

Flags:
  -a, --api              use Docker api
  -c, --config string    config file name
  -h, --help             help for microns
  -m, --microns string   microns config file (default "microns.yaml")
  -o, --output string    topology image (filename only) (default "topo")
  -s, --shell            use shell

Use "microns [command] --help" for more information about a command.
```

- Shell を用いた場合 [**recommend**]

```
sudo ./microns create -s -c examples/basic_ebgp/config.yaml | sudo sh
```

- Docker API を用いた場合

```
sudo ./microns -a -c examples/basic_ebgp/config.yaml
```

## config ファイルのテンプレートを作成

```
sudo ./microns init -c test.yaml

nodes:
  - name:
    image:
    interfaces:
        - inf:
          type:
          ipv4:
          ipv6:
          peernode:
          peerinf:
    volumes:
        - hostvolume:
          containervolume:
    sysctls:
        - sysctl:
node_config:
  - name:
    cmds:
        - cmd:
test:
  - cmds:
        - cmd:

```

## config ファイルを書き換えた場合

一度 delete してから create してください

```
sudo ./microns recreate -s -c examples/basic_ebgp/config.yaml | sudo sh
```

## status を一気に見たい場合

```
sudo ./microns status -c examples/basic_ebgp/config.yaml

----------------------------------------------
                   STATUS
----------------------------------------------
{"name":"R0","status":{"ns":"Found","container":"running"}}
{"name":"R1","status":{"ns":"Found","container":"running"}}
{"name":"R2","status":{"ns":"Found","container":"running"}}
{"name":"R3","status":{"ns":"Found","container":"running"}}
{"name":"C0","status":{"ns":"Found","container":"running"}}
{"name":"C1","status":{"ns":"Found","container":"running"}}
```

## トポロジー図を画像で保存

`-o` で画像のファイル名および dot ファイル名で出力

For Example `-o ebgp => ebgp.dot, ebgp.png`

```
./microns image -c examples/basic_ebgp/config.yaml -o ebgp
```

## トポロジー図を ascii graph で表示

```
./microns image -c examples/basic_ebgp/config.yaml -o ebgp
cat ebgp.dot | graph-easy --from=dot --as_ascii
```

## tinet config -> microns config への変換

```
./microns convert -c tinetcfg.yaml -o micronscfgfile
```
