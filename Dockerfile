FROM scratch

ARG evm_binary
ARG evm_dir
ARG evm_con_data_dir='/evm_data_dir'
ARG evm_con_bin='/evm'

COPY "$evm_binary" "$evm_con_bin"
COPY "$evm_dir" "$evm_con_data_dir"

EXPOSE 1338
EXPOSE 1339
EXPOSE 8080

ENTRYPOINT ["$evm_con_bin", "run", "--datadir", "$evm_con_data_dir"]
