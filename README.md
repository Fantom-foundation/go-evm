## Build Details

[![version](https://img.shields.io/github/tag/Fantom-foundation/go-evm.svg?style=flat-square&logo=github
)](https://github.com/Fantom-foundation/go-evm/releases/latest)  
[![appveyor](https://img.shields.io/appveyor/ci/andrecronje/go-evm.svg?style=flat-square&logo=appveyor)](https://ci.appveyor.com/project/andrecronje/go-evm)  
[![license](https://img.shields.io/github/license/Fantom-foundation/go-evm.svg?style=flat-square&logo=github)](LICENSE.md)  
[![libraries.io dependencies](https://img.shields.io/librariesio/github/Fantom-foundation/go-evm.svg?style=flat-square&logo=librariesio)](https://libraries.io/github/Fantom-foundation/go-evm)  

## Code Quality

[![Go Report Card](https://goreportcard.com/badge/github.com/Fantom-foundation/go-evm?style=flat-square&logo=goreportcard)](https://goreportcard.com/report/github.com/Fantom-foundation/go-evm)  
[![GolangCI](https://golangci.com/badges/github.com/Fantom-foundation/go-evm.svg?style=flat-square&logo=golangci)](https://golangci.com/r/github.com/Fantom-foundation/go-evm)   
[![Code Climate Maintainability Grade](https://img.shields.io/codeclimate/maintainability/Fantom-foundation/go-evm.svg?style=flat-square&logo=codeclimate)](https://codeclimate.com/github/Fantom-foundation/go-evm)  
[![Code Climate Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/Fantom-foundation/go-evm.svg?style=flat-square&logo=codeclimate)](https://codeclimate.com/github/Fantom-foundation/go-evm)  
[![Code Climate Technical Dept](https://img.shields.io/codeclimate/tech-debt/Fantom-foundation/go-evm.svg?style=flat-square&logo=codeclimate)](https://codeclimate.com/github/Fantom-foundation/go-evm)  
[![Codacy code quality](https://img.shields.io/codacy/grade/c8c27910210f4b23bcbbe8c60338b1d5.svg?style=flat-square&logo=codacy)](https://app.codacy.com/project/andrecronje/go-evm/dashboard)  
[![cii best practices](https://img.shields.io/cii/level/2409.svg?style=flat-square&logo=cci)](https://bestpractices.coreinfrastructure.org/en/projects/2409)  
[![cii percentage](https://img.shields.io/cii/percentage/2409.svg?style=flat-square&logo=cci)](https://bestpractices.coreinfrastructure.org/en/projects/2409)  
  
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square&logo=godoc)](https://godoc.org/github.com/Fantom-foundation/go-evm)   

[Documentation](https://github.com/Fantom-foundation/fantom-documentation/wiki).  

[![Sonarcloud](https://sonarcloud.io/api/project_badges/quality_gate?project=Fantom-foundation_go-evm)](https://sonarcloud.io/dashboard?id=Fantom-foundation_go-evm)  
  
## GitHub


[![Commit Activity](https://img.shields.io/github/commit-activity/w/Fantom-foundation/go-evm.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/go-evm/commits/master)  
[![Last Commit](https://img.shields.io/github/last-commit/Fantom-foundation/go-evm.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/go-evm/commits/master)  
[![Contributors](https://img.shields.io/github/contributors/Fantom-foundation/go-evm.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/go-evm/graphs/contributors)  
[![Issues][github-issues-image]][github-issues-url]  
[![LoC](https://tokei.rs/b1/github/Fantom-foundation/go-evm?category=lines)](https://github.com/Fantom-foundation/go-evm)  

[![Throughput Graph](https://graphs.waffle.io/Fantom-foundation/go-evm/throughput.svg)](https://waffle.io/Fantom-foundation/go-evm/metrics/throughput)  

## Social

[![](https://img.shields.io/gitter/room/nwjs/nw.js.svg?style=flat-square)](https://gitter.im/fantom-foundation)    
[![twitter][twitter-image]][twitter-url]  


[codecov-image]: https://codecov.io/gh/fantom-foundation/go-evm/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/fantom-foundation/go-evm
[twitter-image]: https://img.shields.io/twitter/follow/FantomFDN.svg?style=social
[twitter-url]: https://twitter.com/intent/follow?screen_name=FantomFDN
[github-issues-image]: https://img.shields.io/github/issues/Fantom-foundation/go-evm.svg?style=flat-square&logo=github
[github-issues-url]: https://github.com/Fantom-foundation/go-evm/issues


## Usage

The **lachesis_addr** option specifies the endpoint where the consensus node is listening  
to the VM.

The **proxy_addr** option specifies the endpoint where the VM is listening for consensus events.  

```
NAME:
   evm run -

USAGE:
   evm run [command options] [arguments...]

OPTIONS:
   --datadir value        Directory for the databases and keystore (default: "$HOME/.evm")
   --lachesis_addr value  IP:Port of consensus node (default: "127.0.0.1:1338")
   --proxy_addr value     IP:Port to bind Proxy server (default: "127.0.0.1:1339")
   --api_addr value       IP:Port to bind API server (default: ":8080")
   --log_level value      Debug, info, warn, error, fatal, panic (default: "debug")
   --pwd value            Password file to unlock accounts (default: "$HOME/.evm/pwd.txt")
   --db value             Database file (default: "$HOME/.evm/chaindata")
   --cache value          Megabytes of memory allocated to internal caching (min 16MB / database forced) (default: 128)
```

## Configuration

The application writes data and reads configuration from the directory specified  
by the --datadir flag. The directory structure **MUST** be as follows:
```
host:~/.evm$ tree
eth
├── genesis.json
└── keystore
    ├── [Ethereum Key File]
    ├── ...
    ├── ...
    ├── [Ethereum Key File]


```
The Ethereum genesis file defines Ethereum accounts . This file is useful to
predefine a set of accounts that own all the initial Ether at the inception
of the network.  

Example Ethereum genesis.json defining two account:
```json
{
   "alloc": {
        "6cC5F688a315f3dC28A7781717a9A798a59fDA7b": {
            "balance": "1000000000000000000"
        },
        "408d0D182a0397b334a4465Fbe37f3888eE579A7  ": {
            "balance": "1000000000000000000"
        }
   }
}
```

### Get controlled accounts

example:
```bash
host:~$ curl http://[api_addr]/accounts -s | json_pp
{
   "accounts" : [
      {
         "address" : "0x6cC5F688a315f3dC28A7781717a9A798a59fDA7b",
         "balance" : 1000000000000000000,
         "nonce": 0
      }
   ]
}
```
### Get any account

```bash
host:~$ curl http://[api_addr]/account/0x629007eb99ff5c3539ada8a5800847eacfc25727 -s | json_pp
{
    "address":"0x629007eb99ff5c3539ada8a5800847eacfc25727",
    "balance":1000000000000000000,
    "nonce":0
}
```

### Send transactions from controlled accounts

example: Send Ether between accounts  
```bash
host:~$ curl -X POST http://[api_addr]/tx -d '{"from":"0x629007eb99ff5c3539ada8a5800847eacfc25727","to":"0xe32e14de8b81d8d3aedacb1868619c74a68feab0","value":6666}' -s | json_pp
{
   "txHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf"
}
```

### Get Transaction receipt
example:
```bash
host:~$ curl http://[api_addr]/tx/0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf -s | json_pp
{
   "to" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
   "root" : "0xc8f90911c9280651a0cd84116826d31773e902e48cb9a15b7bb1e7a6abc850c5",
   "gasUsed" : "0x5208",
   "from" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
   "transactionHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf",
   "logs" : [],
   "cumulativeGasUsed" : "0x5208",
   "contractAddress" : null,
   "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
}

```
### Send raw signed transactions

example:
```bash
host:~$ curl -X POST http://[api_addr]/rawtx -d '0xf8628080830f424094564686380e267d1572ee409368e1d42081562a8e8201f48026a022b4f68bfbd4f4c309524ebdbf4bac858e0ad65fd06108c934b45a6da88b92f7a046433c388997fd7b02eb7128f4d2401ef2d10d574c42edf15875a43ee51a1993' -s | json_pp
{
    "txHash":"0x5496489c606d74ad7435568393fa2c4619e64497267f80864109277631aa849d"
}
```  
