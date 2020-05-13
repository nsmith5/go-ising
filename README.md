# go-ising

[![Build Status](https://cloud.drone.io/api/badges/nsmith5/go-ising/status.svg)](https://cloud.drone.io/nsmith5/go-ising)

[Ising model](https://en.wikipedia.org/wiki/Ising_model) simulation written in
Go. Streams state as mjpeg movie. To build simply run:

```shell
$ git clone https://github.com/nfsmith5/go-ising
$ pushd go-ising
$ go get
$ go build
```

To run, execute the binary

```shell
$ ./go-ising
2020/05/12 23:04:17 Server binding to :8080...

```

If you point your browser to `localhost:8080` you'll see a movie streaming with
the realtime status of the simulation. Beta (the inverse temperature scaled by
Boltzman's contant) can be controlled with a simple POST require.

```shell
$ # Modify beta to 1
$ curl -X POST -d '1' localhost:8080
```
