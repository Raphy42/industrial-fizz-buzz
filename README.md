# FizzBuzz SaaS
_Production grade fizzbuzz_

# Usage
The docker compose stack runs both the api, a small client written in `typescript` using `Deno` which emulates some load, and a godoc container rendering the package documentation.
The api container should be around 16Mo while the client 
```bash
# running the stack
mv .env.example .env
docker compose up -d
```
# Links (assuming default .env config and a running docker-compose stack)
- [openapi.json](http://localhost:8080/openapi.json)
- [fizz buzz top request](http://localhost:8080/api/v1/metrics/request/fizzbuzz)
- [all top requests](http://localhost:8080/api/v1/metrics/request)
- [godoc](http://localhost:6060/pkg/github.com/Raphy42/industrial-fizz-buzz/)

## Implementation
### Third parties
- Mux/Routing/HttpBody: [echo](https://echo.labstack.com/)
- Logging: [zap](https://github.com/uber-go/zap)
- Testing: [testify](https://github.com/stretchr/testify)
- Config: [envconfig](https://github.com/kelseyhightower/envconfig)
- go 1.20+ (generics)
- openapi3.json automatic generation at runtime
### `core` package
Contains various conveniences and helpers.  
It can be refactored into its own package/applicative-framework if the `core/config` package becomes more generalised, and not `fizzbuzz` specific.
- errors: convenience error wrapper
- config: runtime config from environment + sane defaults
- http: generic handler + server abstraction
- logger: logging layer initialisation + helpers
- semconv: formatting keys and naming things
- generics: slice/maps generic utilities
### `cmd` package
Package for everything binary related: servers, migration runners, scripts...
- `server`: the actual http server

## Environment
```bash
# see core/config/env.go for the complete list
# and default values
FIZZBUZZ_MODE=prod|dev
FIZZBUZZ_PORT=8080
# supersedes FIZZBUZZ_PORT
FIZZBUZZ_ADDR=:8080
# defaults to true
FIZZBUZZ_CORS_ENABLED=true|false
```