# Traffic Balancer

<p align="center">
  <img src="https://i.postimg.cc/prDRd08h/logo.gif" style="border-radius:10px;"/>
</p>

**Load Balancing Algorithm**
- least-time
- weighted-round-robin
- connection-per-time
- round-robin

## Description :books:

- least-time: The load balancer will analyse and monitor the lower average latency of the provided servers and redirect the request to that server.
- weighted-round-robin: On booting the load balance, it will read from the configuration file the weight of every server and redirect the request to the server with the highest weight rate.
- connection-per-second: On booting the load balancer, it will check the number of connections on every server, and it will make that server inactive if the number of requests redirected per 1 minute exceeded the provided configuration of that server.
- round-robin: It will redirect the request per index.

**The load balancer checks the server's health status every 1 minute.**

## Configuration :construction:

Updating the `config/config.yaml` file:
```bash

# least-time
# weighted-round-robin
# connection-per-time
# round-robin
algorithm: "weighted-round-robin" # Algorithm to be used
port: 3030 # Port that the reverse proxy will run on
strict: true # Strict mode for black-listing IPs
log: true # Save logs to file in log folder
access-control-allow-origin: "*" # CORS HEADERS
max-body-size: 1024 # Maximum body size!
xss-protection: true # In default true it will add xss protection header to the requested sent from load balancer to the server.
servers: #List of servers.
  - 
    host: "http://localhost:9876"
    weight: 0.1
    connections: 1
  
  - 
    host: "http://localhost"
    weight: 0.9
    connections: 100

```

## Usage :rocket:

```bash
go build
```

```bash
./traffic-balancer 
Load Balancer started at :3030
```

:warning: **don't forget to add your servers in the `config/config.yaml` file**

## TO-DO :man_playing_handball:	
- ~~Adding more `Access-Control-Allow` header.~~
- Add more load balancer algorithms.
- Write unit tests.
- ~~Add XSS Header protection.~~
- ~~Add Body size.~~
- Add MTLS.
- Add Max File Size!
- ~~Strict mode for blocking black listing IPs.~~
- Gzip Compression.

## LICENSE :balance_scale:

This project is licensed under the MIT License. See the [LICENSE](https://github.com/AAVision/traffic-balancer/blob/main/LICENSE) file for details.