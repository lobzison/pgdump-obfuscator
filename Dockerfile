FROM golang:1.9-stretch

COPY . /pgdump-obfuscator

RUN cd /pgdump-obfuscator &&\
    go build . &&\
    mv /pgdump-obfuscator/pgdump-obfuscator /usr/sbin/pgdump-obfuscator &&\
    rm -rf /pgdump-obfuscator
