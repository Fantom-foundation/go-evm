FROM scratch

ARG evm_binary
ARG evm_dir
ARG evm_con_data_dir='/evm_data_dir'
ARG evm_con_bin='/evm'

COPY "$evm_binary" "$evm_con_bin"
COPY "$evm_dir" "$evm_con_data_dir"
COPY /usr_bins/mkdir /mkdir
RUN mkdir /bin
COPY /usr_bins/mv /bin/mv
COPY /usr_bins/rm /bin/rm
RUN /bin/mv /mkdir /bin/mkdir
COPY sh /bin/sh
COPY ls /bin/ls

EXPOSE 1338
EXPOSE 1339
EXPOSE 8080

RUN [/bin/ls]
ENTRYPOINT [$evm_con_bin, "run", "--datadir", $evm_con_data_dir]
