
# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install

# install compiles and places the binary in GOPATH/bin
install:
	go install \
	 	--ldflags '-extldflags "-static"' \
		./cmd/evm

# build compiles and places the binary in /build
build:
	go build \
		--ldflags '-extldflags "-static"' \
		-o build/evm ./cmd/evm/

# dist builds binaries for all platforms and packages them for distribution
dist:
	@BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/dist.sh'"

test:
	glide novendor | xargs go test

.PHONY: vendor install build test
