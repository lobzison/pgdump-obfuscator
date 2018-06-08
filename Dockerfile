FROM golang:1.10.2-stretch

COPY . /pgdump-obfuscator

RUN cd /pgdump-obfuscator &&\
    go clean &&\
    go build . &&\
    mv /pgdump-obfuscator/pgdump-obfuscator /usr/sbin/pgdump-obfuscator &&\
    rm -rf /pgdump-obfuscator
