# Service Discovery Demo

> **Disclaimer** This code is a PoC and it's absolutely not production Ready.
> 
> The code is the result of 1 week learning Go in spare time + 1 day of research on Zeroconf
## Overview

In this PoC we implement the following algorithm:

- look for a service instance in the local network
- if the service exists, use it
- if no instances are found, promote itself as service

## Build
The command
```shell
go build .
```
produces the `service_discovery` executable

During development it's possible to use
```
go run .
```
to build and run the executable in a single step

## Usage
Just run with 
```
./service_discovery
```

It's possible to specify the TCP port used by the service with
```
./service_discovery -port <number>
```
the default port is 8080

Use 
```
./service_discover -h
```

to show the command line options

