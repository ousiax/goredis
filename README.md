# Redis Go
Redis Go Client &amp; Cluster

## Play with the litte gedget

    $ cd cli/
    $ go install
    $ cd $GOPATH/bin
    $ ./cli 
    127.0.0.1:6379>ping
    PONG
    127.0.0.1:6379>set foo buzz
    OK
    127.0.0.1:6379>get foo
    buzz
    127.0.0.1:6379>lpush bar 0 1 2
    3
    127.0.0.1:6379>lrange bar 0 -1
    [2 1 0]
    127.0.0.1:6379>save
    OK
    127.0.0.1:6379>quit

***Note***: you can also type `./cli -p=<port> -h=<hostname>`.