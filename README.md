# Traffic Balancer

<p align="center">
  <img src="https://i.postimg.cc/prDRd08h/logo.gif" style="border-radius:9px;"/>
</p>

**Load Balancing Algorithm**
- least-time
- weighted-round-robin
- connection-per-time
- round-robin

## Configuration :construction:

Updating the `config/config.yaml` file:
```bash

# least-time
# weighted-round-robin
# connection-per-time
# round-robin
algorithm: "weighted-round-robin" # Algorithm to be used
port: 3030 # Port that the reverse proxy will run on
strict: true # strict mode for black-listing IPs
log: true # save logs to file in log folder
servers: #list of servers.
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

## LICENSE :balance_scale:

This project is licensed under the MIT License. See the [LICENSE](https://github.com/AAVision/traffic-balancer/blob/main/LICENSE) file for details.