# Traffic Balancer

<p align="center">
  <img src="https://github.com/AAVision/traffic-balancer/blob/e634803d0c108daa1cc6892837abf4d0d21bcd4a/logo.gif"  style="border-radius:9px;"/>
</p>

**Load Balancing Algorithm**
- least-time
- weighted-round-robin
- connection-per-second
- round-robin

## Configuration 

Updating the `config.yaml` file:
```bash

# least-time
# weighted-round-robin
# connection-per-second
# round-robin
algorithm: "weighted-round-robin" # Algorithm used
port: 3030 # Port that the reverse proxy will run on
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


## LICENSE :balance_scale:

This project is licensed under the MIT License. See the [LICENSE](https://github.com/AAVision/traffic-balancer/blob/main/LICENSE) file for details.