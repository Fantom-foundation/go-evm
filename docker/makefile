BABBLE_VERSION=0.2.1
DUMMY_VERSION=0.2.1
WATCHER_VERSION=0.0.2
GLIDER_VERSION=0.0.2

all: evm-image dummy-image watcher-image glider-image

demo: evm-image dummy-image watcher-image

evm-image:
	CGO_ENABLED=0 go build \
		--ldflags "-X github.com/andrecronje/evm/version.GitCommit=`git rev-parse HEAD`" \
		-o evm/evm ../cmd/evm/
	docker build --no-cache=true -t andrecronje/evm:$(BABBLE_VERSION) evm/

dummy-image:
	CGO_ENABLED=0 go build \
		-o dummy/dummy ../cmd/dummy_client/
	docker build --no-cache=true -t andrecronje/dummy:$(DUMMY_VERSION) dummy/

watcher-image:
	docker build --no-cache=true -t andrecronje/watcher:$(WATCHER_VERSION) -t andrecronje/watcher:latest watcher/

glider-image:
	docker build --no-cache=true -t andrecronje/glider:$(GLIDER_VERSION) glider/

.PHONY: all demo evm-image dummy-image watcher-image glider-image
