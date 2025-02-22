# RoundBalancer - HTTP Round Robin Load Balancer

## Overview

RoundBalancer is an HTTP-based **Round Robin API Load Balancer** that acts as a load balancer for routing HTTP POST requests to multiple instances of an Application API on the backend. It ensures that incoming requests are distributed evenly among available instances using a round-robin algorithm.

## Components

- **Echo Server:** A simple application that accepts POST requests with JSON payload and responds with the same data.
- **LoadBalancer Server:** 
  * Reads backend server configurations from `config/values/proxies_config.json`. 
  * Routes HTTP requests using a round-robin load balancing algorithm. 
  * Periodically checks backend health at configurable intervals.

## Installation & Setup

1. Clone the Repository

```
 git clone https://github.com/onkarbanerjee/roundbalancer.git
 cd roundbalancer
```

2. You will need golang installed, latest is recommended. Once installed, download all dependencies as shown below:

```
go mod tidy
go mod vendor
```

3. You should be set to run and test it now. Mocks are already committed to help make it easier to run.

## How to test it

1. Run multiple instances of echo servers as below:
```
go run cmd/echo/main.go -id=$id -port=$port
```
 Example -
```
go run cmd/echo/main.go -id=1 -port=8080
go run cmd/echo/main.go -id=2 -port=8081
go run cmd/echo/main.go -id=3 -port=8082
```
To know more, refer help - ` go run cmd/echo/main.go --help`.

Any requests to these should be served by them directly now, same should be seen in their logs.

2. Configure the addresses of the running echo servers in `config/values/proxies_config.json`
 
 Example-
```
{
  "backends": [
    {
      "id": "backend-1",
      "port": 8080
    },
    {
      "id": "backend-2",
      "port": 8081
    },
    {
      "id": "backend-3",
      "port": 8082
    }
  ]
}
```

3. Start the loadbalancer server like below.

```
go run cmd/loadbalancer/main.go -port=$port -timeout=$timeout
```
  Example -
```
go run cmd/loadbalancer/main.go -port=9090 -timeout=10
```
To know more, refer help - ` go run cmd/loadbalancer/main.go --help`

4. Send an HTTP POST request to the loadbalancer.
```
curl --location 'http://localhost:9090/echo' \
--header 'Content-Type: application/json' \
--data '{
  "game": "Mobile Legends",
  "gamerID": "GYUTDTE",
  "points": 20
}'
```
This should now be served by the loadbalancer and you should get back 200OK response like below.
```
{"game":"Mobile Legends","gamerID":"GYUTDTE","points":20}
```

5. Monitor the logs of your running echo server instances to check which one is serving the requests. You should see logs like below:
```
{"level":"info","ts":1740243764.542871,"caller":"echo/server.go:53","msg":"Request completed","server_id":"2","method":"POST","url":"/echo"}

```
**Note the `server_id` in the logs** 

6. Try sending multiple requests, you should see them being served by different echo server in round robin manner.
7. You can shutdown any echo server simply by pressing `ctrl+c` . Subsequent requests should skip the shutdown server and go to the echo server.
8. You can start any of the servers again and you should start seeing them included in the rotation again.


## Handling Failures

- If an instance becomes unresponsive, it is temporarily removed from the rotation. 
- When it becomes available again, it automatically rejoins.
- If an instance slows down, a timeout mechanism prevents blocking subsequent requests.

## Assumptions and Future Enhancements

- If the request to backend times out, we simply respond with 502 Bad Gateway. There are no fine tuned configurations for this apart from a basic timeout cmd line argument.
- No retries have been implemented for this project and has been kept out of scope.
- As the name of the project suggests, it only contains a round robin implementation of load balancer. However implementing any other strategy should be possible by just implementing the loadbalancer interface.
- While this specifically serve HTTP requests as of now, the round robin implementation itself could be made more generic to work with any kind of traffic like DNS requests. However that has been kept out of scope for now.
- HTTP responses are just sent directly back to the client. In order to verify which server they being served by, we will need to check the logs for now.
- Liveness check interval is 2 secs by default, but would be trivial to make it configurable.



