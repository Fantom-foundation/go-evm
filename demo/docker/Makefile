.PHONY: all build run evm-xgo

all: build

build: evm
	docker build -t "evm_node" -f ./Dockerfile.node ./
	docker build -t "evm_web3" -f ./Dockerfile.web3 ./

evm:
	go build -a --ldflags '-extldflags "-static"' ../../../evm/cmd/evm


evm-xgo:
	go get github.com/karalabe/xgo
	xgo --targets=linux/amd64 -v /cmd/evm


start:
	docker run --rm -d --name="evm_node_1" -p 8545:8545 -p 8546:8546 "evm_node"

exec:
	docker run --rm -ti --link=evm_node_1:evm -v ${PWD}:/tmp "evm_web3" node /tmp/script.js

stop:
	docker stop "evm_node_1"
    