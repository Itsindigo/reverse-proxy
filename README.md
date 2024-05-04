# Reverse Proxy

todo

1. load route limitations from config, apply to routes as middleware
2. middleware should take in the route path, ip address, and request method. needs to increment a route level counter, and a global counter. if either conunter exceeds the limit, the middleware should return a 429 status code
   3. redis route key should be in the format of `route:<route>:<ip>:<method>`