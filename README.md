# Reverse Proxy

todo

1. load route limitations from config, apply to routes as middleware
2. middleware should take in the route path, ip address, and request method. needs to increment a route level counter, and a global counter. if either conunter exceeds the limit, the middleware should return a 429 status code
   3. redis route key should be in the format of `route:<route>:<ip>:<method>`


- Create refiller
  - infinite loop, every 1 second, refill all keys by 1
- Redis Key Schema
  - On Request - key: `route_bucket:<method>:<route>:<ip>:` -> hash -> decrement if > 0 -> return 429 if 0



REDIS CHEATSHEET

Redis CLI: `docker exec -it proxy_server_redis redis-cli -a boop`

All keys: KEYS *

Del all keys: flushall


## Package Documentation

`task start-servers` will start a `godoc` server on port `6060` for viewing package documentation.

Documentation for the project can be found [here](http://localhost:6060/pkg/github.com/itsindigo/reverse-proxy/?m=all).


## TODO 

- Implement different rate limiting strategies that can be configured via the proxy config file.
- Look at mTLS/Ingress/Firewalls
- Implement an additional rate limiter that tracks a user's requests across all routes.
- Add a boat load of tests including a long running integration test that hits the proxy with a bunch of requests.
- Maybe stick it in a k8s cluster with Services/Ingresses/NetworkPolicies
