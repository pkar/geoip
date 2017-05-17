MaxMind GeoLite API
===================

Usage
=====

This assumes Go is installed and the env GOPATH set. The package should be placed into:

```
$GOPATH/src/github.com/pkar/geoip
```

```
$ go run cmd/geoip/main.go -listen=localhost:8899 -location=./GeoLiteCity-Location.csv -blocks=./GeoLiteCity-Blocks.csv
... some numbers for 20-30s, hopefully you don't run out of memory (uses about 850mb)

2017/05/11 20:59:29 listening on localhost:8899
```

Lookup time is fairly quick locally, about 3ms but haven't had a chance to load test.

```
$ ron hs http://localhost:8899/location?ip=73.240.122.151

Connected to 127.0.0.1:8899

HTTP/1.1 200 OK
Content-Length: 156
Content-Type: application/json
Date: Fri, 12 May 2017 06:47:21 GMT

Body discarded

   DNS Lookup   TCP Connection   Server Processing   Content Transfer
[       1ms  |           0ms  |              0ms  |             0ms  ]
             |                |                   |                  |
    namelookup:1        ms      |                   |                  |
                        connect:2        ms         |                  |
                                      starttransfer:2        ms        |
                                                                 total:3        ms
```

```
$ hey http://localhost:8899/location?ip=73.240.122.151
All requests done.

Summary:
  Total:	0.0187 secs
  Slowest:	0.0131 secs
  Fastest:	0.0001 secs
  Average:	0.0033 secs
  Requests/sec:	10723.7462
  Total data:	31200 bytes
  Size/request:	156 bytes

Status code distribution:
  [200]	200 responses

Response time histogram:
  0.000 [1]	|
  0.001 [86]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.003 [45]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.004 [18]	|∎∎∎∎∎∎∎∎
  0.005 [0]	|
  0.007 [0]	|
  0.008 [10]	|∎∎∎∎∎
  0.009 [27]	|∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.010 [7]	|∎∎∎
  0.012 [2]	|∎
  0.013 [4]	|∎∎

Latency distribution:
  10% in 0.0006 secs
  25% in 0.0009 secs
  50% in 0.0016 secs
  75% in 0.0072 secs
  90% in 0.0088 secs
  95% in 0.0094 secs
  99% in 0.0130 secs
```

Block Lookup
============

This library assumes the GeoLiteCity-Blocks.csv are ip range blocks in sorted order and contiguous in
order to do faster (log n) look ups by having an array of range boundaries (16777216 - 3758096383) and associated ids.

With this assumption you can do a binary search by using the endIpNum in ranges.

Another assumption is startIpNum begins with 16777216. If the csv changes it will lead to incorrect
values in the initial range.

```
startIpNum,endIpNum,locId
"16777216","16777471","609013"
"16777472","16778239","104084"
```

```
ranges := []int{-1, 16777215, 16777471, 16778239, ...}
ids    := []int{nil, nil, blockn, blocka, ...}
```

API
===

The API provides a single endpoint for looking up a location by either location id or ip address IPv4.

```
$ curl localhost:8899/location?locId=24107
{"locId":24107,"country":"US","region":"NY","city":"Fonda","postalCode":"12068","latitude":42.9508,"longitude":-74.3937,"metroCode":532,"areaCode":"518"}
```

```
$ curl localhost:8899/location?ip=73.240.122.151
{"locId":4004,"country":"US","region":"OR","city":"Portland","postalCode":"97202","latitude":45.4763,"longitude":-122.6408,"metroCode":820,"areaCode":"503"}
```


DB
==

The input assumes unzipped csv files from this GeoLiteCity link
http://geolite.maxmind.com/download/geoip/database/GeoLiteCity_CSV/GeoLiteCity-latest.zip

GeoLiteCity-Location.csv

```
Copyright (c) 2012 MaxMind LLC.  All Rights Reserved.
locId,country,region,city,postalCode,latitude,longitude,metroCode,areaCode
1,"O1","","","",0.0000,0.0000,,
```

GeoLiteCity-Blocks.csv

```
Copyright (c) 2017 MaxMind Inc.  All Rights Reserved.
startIpNum,endIpNum,locId
"16777216","16777471","609013"
```

Caveats
=======

- Startup load times are slowish, it loads both csv's into memory as well as its serialized locations. But memory is cheap these days.
