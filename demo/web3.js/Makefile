.PHONY: web3_image start stop exec


web3_image:
	test -z $$(docker image ls -q "evm_web3") \
	&& docker build -t "evm_web3" -f ./Dockerfile.web3 ./ \
	|| echo Docker image "evm_web3" already exists

start:
	go run ../../../evm/cmd/evm --log=debug solo & echo $$! > evm.pid

stop:
	pkill -P $$(cat evm.pid)
	rm -f evm.pid

exec: web3_image
	docker run --rm -ti -v ${PWD}:/tmp --network="host" "evm_web3" node /tmp/example.js
