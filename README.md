# Kan's round robin load balancer

This project is a simple Go stateless repo and contains 2 different APIs.
- A simple api that will return whatever JSON string it is given. We will spin up 3 instances of this on bootup in different ports. It's meant for testing the behaviour of our load balancer. If you use a seed that's not 0, we will use a fixed random behaviour to simulate the error and latency of the different intances. We can use these fixed random behaviours in our benchmarks to test out our scoring algorithm and optimise it be a more performant load balancer.
- A load balancer api that will round robin requests to the designated simple API ports

## Design of the round robin connection pool

I started out building a normal round robin api through all the different ports spun up in the simple api. Then I introduced a scoring system to ensure our connection pool favours higher performing connections than lower ones whilst continuing to follow the round robin nature of traffic distribution.

The scoring system has the following characteristics:
- The lower the score, the better the connection health
- Will penalise connections with slower latency than the lastest system average
- Will slightly favour connections with a faster latency than the latest system average
- Will penalise connections responding with errors or a non successful response
- Will slightly favour connections with no recent errors
- Once the score is high enough to reach the penalty threshold, we will penalise the connection by blocking traffic to it for n times through the round robin. The number of times to block will depend how high the connection's score is above the threshold.

## How to run

Please ensure you have Go installed on your computer.

### Running the APIs

```go run cmd/load-balancer-api/main.go```

```go run cmd/simple-api/main.go```

### Testing

```go test ./...```

### Benchmarking

```go run testing/benchmark.go```

## Benchmarking guide

I built a benchmarking script in Go to test a certain number of requests using a defined number of workers. We need to spin up the simple api and the load balancer api with a seeded number then run benchmark file to call against our load balancer api.

We should test with various seeded numbers (different error and latency simulations) and different instance numbers to find the optimised settings for our scoring system.

## Room for improvement had I had more time
- I want to add health check endpoints on the simple API. We can continuously monitor for 200 status and a low latency response. We can use background threads to periodically call health check and use it in our connection's score calculation.
- Improve the simulator so instances can fail or improve over time. Currently the simulator can only simulate the same behaviour of an instance over time. This is not realistic as our connections can become faster or slower and start throwing errors at anytime in real life.
