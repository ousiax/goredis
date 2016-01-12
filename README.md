# Redis Go
Redis Go Client &amp; Cluster

### Play with the litte gedget

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

### A simple example of the naive `redis.Client`

    package main
    
    import (
        "flag"
        "fmt"
        "github.com/qqbuby/goredis/redis"
        "strconv"
    )
    
    func main() {
        hostnamePtr := flag.String("h", "127.0.0.1", "Server hostname (default: 127.0.0.1).")
        portPtr := flag.Int("p", 6379, "Server port (default: 6379).")
        flag.Parse()
    
        network := "tcp"
        address := *hostnamePtr + ":" + strconv.Itoa(*portPtr)
    
        client, err := redis.NewClient(network, address)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer client.Close()
    
        var key string = "key"
        client.Set(key, "hello world")
        v, _ := client.Get(key)
        fmt.Println(v)
    }
    