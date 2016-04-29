sufttun
===

extremely fast & secure udp tunnel based on suft

Modified from [kcptun](https://github.com/xtaci/kcptun) and changed kcp to [suft](https://github.com/spance/suft) (Small-scale UDP Fast Transmission).

```
client <--> suft-client  <--> suft-server  <--> server
    tcp <-> udp                        udp <-> tcp
```


usage
---
```
server:
go get github.com/messyidea/sufttun/server

server -l "addr:port" -r "addr:port" -b 10 -key "your key" -tuncrypt true

-l: local addr
-r: remote addr
-b: max bandwidth of sending in mbps
-key: your key
-tuncrypt: encrypt


client:
go get github.com/messyidea/sufttun/client

client -l "addr:port" -r "addr:port" -b 10 -key "your key" -tuncrypt true

-l: local addr
-r: remote addr
-b: max bandwidth of sending in mbps
-key: your key
-tuncrypt: encrypt

```


sample
---
```
server -l ":porta" -r "127.0.0.1:portb" -b 10 -key "yourkey" -tuncrypt true
client -l ":portc" -r "server:portd" -b 10 -key "yourkey" -tuncrypt true
```

license
---
MIT
