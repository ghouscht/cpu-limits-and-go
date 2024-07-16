# CPU limits & Go
Lightning talk with demo on how CPU limits affect your Go application in containerized environments. The talk was
held at the [Bärner Go Talks 2024 no. 3](https://www.meetup.com/berner-go-meetup/events/301976804)

The slides can be found here: [CPU limits & Go](https://docs.google.com/presentation/d/18qBafo1hClKQYnM1ltsmBt91sBQfi7j-DjKKPdz4vyk)

## Build and run the example
[goreleaser](https://github.com/goreleaser/goreleaser) is used to build a container image out of the example code. After
installing goreleaser, you can build a snapshot release like this:
```shell
goreleaser release --snapshot --clean
```
Now you can run the container e.g. using docker:
```shell
docker run -p '8080:8080' --detach --cpus 1 --name cpu-limits-and-go ghouscht/cpu-limits-and-go
```
Please note that in the example above the container is limited to 1 CPU.

## Experiments
After you have built and started the container, you can run some experiments to see how the CPU limits affect the latency
of the application.

First of all you can check how many CPUs Go detects and uses:
```shell
$ curl http://localhost:8080/maxProcs
Go is using 6 CPUs and there are 6 CPUs available
```
Please note that the number of CPUs might be different on your machine, depending on the number of CPUs available. However,
you can already see that Go detects all available CPUs and seems to miss the CPU limit (In this example the container was
started using `--cpus 1` flag).

We can now use a load testing tool like [hey](https://github.com/rakyll/hey) to send some requests to the application and
see how the latency changes with different settings of `GOMAXPROCS`.

### default GOMAXPROCS
- CPU limit = 1
- GOMAXPROCS = 6

```shell
$ hey http://127.0.0.1:8080/isPrime/666666

Summary:
  Total:        1.4884 secs
  Slowest:      0.8120 secs
  Fastest:      0.0106 secs
  Average:      0.2936 secs
  Requests/sec: 134.3720

  Total data:   5400 bytes
  Size/request: 27 bytes

Response time histogram:
  0.011 [1]     |■
  0.091 [22]    |■■■■■■■■■■■■■■■■■■■■■
  0.171 [30]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.251 [35]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.331 [41]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.411 [29]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.491 [10]    |■■■■■■■■■■
  0.572 [16]    |■■■■■■■■■■■■■■■■
  0.652 [5]     |■■■■■
  0.732 [5]     |■■■■■
  0.812 [6]     |■■■■■■


Latency distribution:
  10% in 0.0800 secs
  25% in 0.1406 secs
  50% in 0.2940 secs
  75% in 0.4003 secs
  90% in 0.5105 secs
  95% in 0.7005 secs
  99% in 0.8093 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0004 secs, 0.0106 secs, 0.8120 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0001 secs
  resp wait:    0.2930 secs, 0.0106 secs, 0.8120 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0005 secs

Status code distribution:
  [200] 200 responses
```

### GOMAXPROCS = 3
- CPU limit = 1
- GOMAXPROCS = 3

The program allows you to configure `GOMAXPROCS` using a REST like endpoint. Let's change the number of CPUs Go uses to 3:
```shell
$  curl -X POST http://localhost:8080/maxProcs/3
Go is now using 3 CPUs, previously it was using 6 CPUs
```

Now we can run the same test as before:
```shell
$ hey http://127.0.0.1:8080/isPrime/666666

Summary:
  Total:        0.7244 secs
  Slowest:      0.5284 secs
  Fastest:      0.0050 secs
  Average:      0.1195 secs
  Requests/sec: 276.1089
  
  Total data:   5400 bytes
  Size/request: 27 bytes

Response time histogram:
  0.005 [1]     |
  0.057 [51]    |■■■■■■■■■■■■■■■■■■■■■■
  0.110 [91]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.162 [18]    |■■■■■■■■
  0.214 [10]    |■■■■
  0.267 [5]     |■■
  0.319 [7]     |■■■
  0.371 [2]     |■
  0.424 [8]     |■■■■
  0.476 [2]     |■
  0.528 [5]     |■■


Latency distribution:
  10% in 0.0144 secs
  25% in 0.0357 secs
  50% in 0.0928 secs
  75% in 0.1135 secs
  90% in 0.3111 secs
  95% in 0.4154 secs
  99% in 0.5134 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0004 secs, 0.0050 secs, 0.5284 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0014 secs
  resp wait:    0.1190 secs, 0.0050 secs, 0.5268 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0002 secs

Status code distribution:
  [200] 200 responses
```

### GOMAXPROCS = CPU limit
- CPU limit = 1
- GOMAXPROCS = 1

```shell
$ curl -X POST http://localhost:8080/maxProcs/1
Go is now using 1 CPUs, previously it was using 3 CPUs
$ hey http://127.0.0.1:8080/isPrime/666666

Summary:
  Total:        0.7630 secs
  Slowest:      0.2098 secs
  Fastest:      0.0105 secs
  Average:      0.1666 secs
  Requests/sec: 262.1121

  Total data:   5400 bytes
  Size/request: 27 bytes

Response time histogram:
  0.011 [1]     |
  0.030 [4]     |■■
  0.050 [4]     |■■
  0.070 [6]     |■■■
  0.090 [6]     |■■■
  0.110 [7]     |■■■
  0.130 [6]     |■■■
  0.150 [5]     |■■
  0.170 [7]     |■■■
  0.190 [92]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.210 [62]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10% in 0.0880 secs
  25% in 0.1763 secs
  50% in 0.1844 secs
  75% in 0.1918 secs
  90% in 0.1959 secs
  95% in 0.1981 secs
  99% in 0.2049 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0005 secs, 0.0105 secs, 0.2098 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0005 secs
  resp wait:    0.1660 secs, 0.0099 secs, 0.2097 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0031 secs

Status code distribution:
  [200] 200 responses
```
