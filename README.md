What
====

Streaming obfuscator of sensitive data in PostgreSQL dumps (pg_dump).

    $ git clone https://github.com/yuriidanyliak/pgdump-obfuscator
    $ cd pgdump-obfuscator
    $ go test
    $ go build // if you changed config.go
    $ go install
    $ pg_dump [...] | pgdump-obfuscator [-c params] > dump.sql
    
Example command to export dump with proper format with pg_dump:

    $ pg_dump -h <HOST> -U <USER> -p 5432 -O -x --no-security-labels --no-tablespaces -W <DATABASE> > servicing_dump_production.sql


Configuration
====

```
Example:

```
pgdump-obfuscator -c auth_user:email:email -c auth_user:password:bytes -c address_useraddress:phone_number:digits -c address_useraddress:land_line:digits
```
