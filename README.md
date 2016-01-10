# Redis Go
Redis Go Client &amp; Cluster

## Play with the litte gedget

    $ cd cli/
    $ go install
    $ cd $GOPATH/bin
    $ ./cli 
    127.0.0.1:6379>ping
    +PONG
    127.0.0.1:6379>set foo buzz
    +OK
    127.0.0.1:6379>get foo
    $4
    buzz
    127.0.0.1:6379>lpush bar 0 1 2
    :3
    127.0.0.1:6379>lrange bar 0 -1
    *3
    $1
    2
    $1
    1
    $1
    0
    127.0.0.1:6379>save
    +OK
    127.0.0.1:6379>shutdown
    127.0.0.1:6379>quit