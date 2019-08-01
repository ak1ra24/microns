## MicroNS
Dockerを用いたルーティングシュミレーション

## Prerequisites
* Go
    * export GO111MODULE=on

Ubuntu18.04
* Go Install
```
sudo add-apt-repository ppa:longsleep/golang-backports -y
sudo apt update
sudo apt install golang-go
```

* Graphviz & ascii graph for microns image
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
  create      create docker container and ns topology
  delete      delete docker container and ns topology
  help        Help about any command
  image       create network topology image file
  init        A brief description of your command
  recreate    reconfigure router
  status      status docker container and ns topology
  test        Execute test from config

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

## configファイルのテンプレートを作成
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

## configファイルを書き換えた場合
一度deleteしてからcreateしてください
```
sudo ./microns recreate -s -c examples/basic_ebgp/config.yaml | sudo sh
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

## トポロジー図を画像で保存
`-o` で画像のファイル名およびdotファイル名で出力

For Example `-o ebgp => ebgp.dot, ebgp.png`

```
./microns image -c examples/basic_ebgp/config.yaml -o ebgp
```

## トポロジー図をascii graphで表示
```
./microns image -c examples/basic_ebgp/config.yaml -o ebgp
cat ebgp.dot | graph-easy --from=dot --as_ascii
```
