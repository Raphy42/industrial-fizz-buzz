# FizzBuzz SaaS
_Production grade fizzbuzz_

# Usage
```bash
mv .env.example .env
docker compose up -d
```

## Implementation
### Third parties
- Mux/Routing/HttpBody: [echo](https://echo.labstack.com/)
- Logging: [zap](https://github.com/uber-go/zap)
- Testing: [testify](https://github.com/stretchr/testify)
- Config: [envconfig](https://github.com/kelseyhightower/envconfig)
- go 1.20+ (generics)
### `core` package
Contains various conveniences and helpers.  
It can be refactored into its own package/applicative-framework if the `core/config` package becomes more generalised, and not `fizzbuzz` specific.
- errors: convenience error wrapper
- config: runtime config from environment + sane defaults
- http: handler + server abstraction
- logger: logging layer initialisation + helpers
- semconv: formatting keys and naming things
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