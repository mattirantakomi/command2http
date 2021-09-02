# command2http

start:

    $ COMMAND="ping google.com" command2http
    listen :8080

then:

    $ curl localhost:8080
    PING google.com (142.250.74.78): 56 data bytes
    64 bytes from 142.250.74.78: icmp_seq=0 ttl=59 time=24.822 ms
    64 bytes from 142.250.74.78: icmp_seq=1 ttl=59 time=16.339 ms
    64 bytes from 142.250.74.78: icmp_seq=2 ttl=59 time=19.575 ms
    ...


## install

    go get github.com/matti/command2http
