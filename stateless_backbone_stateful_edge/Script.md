# Script

This script generates the operator, account and user JWTs for the scenario. The `nats` CLI tool is required, along with the `nats-server` 2.9.16 binary, which is the latest release.

```bash
nsc add operator --generate-signing-key --sys --name statlessbackbone

nsc edit operator --require-signing-keys --account-jwt-server-url nats://localhost:4222

nsc add account acc-us-east
nsc edit account acc-us-east --sk generate

nsc add account acc-global
nsc edit account acc-global --sk generate

nsc add account acc-eu-lon
nsc edit account acc-eu-lon --sk generate

nsc add user user-eu-lon -a acc-eu-lon 
nsc add user user-us-east -a acc-us-east
nsc add user user-global -a acc-global
nsc add user user-global-eu-lon -a acc-global
nsc add user user-global-us-east -a acc-global

nats context save fc-us-east --nsc nsc://statlessbackbone/acc-us-east/user-us-east
nats context save fc-eu-lon --nsc nsc://statlessbackbone/acc-eu-lon/user-eu-lon
nats context save fc-global --nsc nsc://statlessbackbone/acc-global/user-global
nats context save fc-global-eu-lon --nsc nsc://statlessbackbone/acc-global/user-global-eu-lon
nats context save fc-global-us-east --nsc nsc://statlessbackbone/acc-global/user-global-us-east
nats context save fc-sys --nsc nsc://statlessbackbone/SYS/SYS

# Geneerate the resolver.conf for the server config
nsc generate config --nats-resolver --sys-account SYS > resolver.conf

# Export credentials for the users
nsc generate creds -a acc-global -n user-global-eu-lon -o user-global-eu-lon.creds
nsc generate creds -a acc-global -n user-global-us-east -o user-global-us-east.creds
nsc generate creds -a SYS -n SYS -o SYS.creds 

nsc edit account acc-global --js-mem-storage 1G --js-disk-storage 100G --js-streams 10 --js-consumer 100

# Ok, we're done here. Start the nodes and push.
nsc push --all && nsc pull --all
```
