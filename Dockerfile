FROM scratch
ARG evm_binary
COPY "$evm_binary" /evm
ENTRYPOINT ["/evm"]
