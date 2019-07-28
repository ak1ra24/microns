## MicroNS
Dockerを用いたルーティングシュミレーション

## Prerequisites
* Go
    * export GO111MODULE=on

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
  create      create docker container and ns topology
  delete      delete docker container and ns topology
  help        Help about any command
  image       create network topology image file
  status      status docker container and ns topology

Flags:
  -a, --api             use Docker api
  -c, --config string   config file name
  -h, --help            help for microns
  -o, --output string   topology image (filename only) (default "topo")
  -s, --shell           use shell

Use "microns [command] --help" for more information about a command.
```

* Shellを用いた場合 [**recommend**]
```
sudo ./microns create -s -c examples/basic_ebgp/config.yaml | sudo sh
```

* Docker APIを用いた場合
```
sudo ./microns -a -c examples/basic_ebgp/config.yaml
```

## configファイルを書き換えた場合
一度deleteしてからcreateしてください
```
sudo ./microns delete -s -c examples/basic_ebgp/config.yaml | sudo sh
sudo ./microns create -s -c examples/basic_ebgp/config.yaml | sudo sh
```

## statusを一気に見たい場合
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
