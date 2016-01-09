# Redis Go
Redis Go Client &amp; Cluster

## Play with the litte gedget
    $ cd cli/
    # go run cli.go
    127.0.0.1:6379>set mykey myvalue
    +OK
    127.0.0.1:6379>ping
    +PONG
    127.0.0.1:6379>save
    +PONG
    127.0.0.1:6379>shutdown
    +OK
    127.0.0.1:6379>exit
