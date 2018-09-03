FROM glider as builder

RUN mkdir -p "$GOPATH/src/github.com/andrecronje/evm" && \
    git clone 'https://github.com/andrecronje/evm' "$GOPATH/src/github.com/andrecronje/evm" && \
    cd "$GOPATH/src/github.com/andrecronje/evm" && \
    glide install && \
    cd "$GOPATH/src/github.com/andrecronje/evm/cmd/evm" && \
    go build -ldflags "-linkmode external -extldflags -static" -a main.go && \
    mv "$GOPATH/src/github.com/andrecronje/evm/cmd/evm/main" "$GOPATH/bin/evm"

FROM scratch as evm

EXPOSE 1338
EXPOSE 1339
EXPOSE 8080

# cp -r /etc/ssl/certs certs, then add to your `docker build`: `--build-arg ca_certificates=certs`
ARG ca_certificates
ADD "$ca_certificates" /etc/ssl/certs/
COPY --from=0 /go/bin/evm /evm

# `mkdir -p` cheat
COPY 'Dockerfile' "/evm_data_dir/.gitignore"

ENTRYPOINT ["/evm", "run", "--datadir", "/evm_data_dir", "--db", "/evm_data_dir/chaindata", "--pwd", "/evm_data_dir/pwd.txt"]
