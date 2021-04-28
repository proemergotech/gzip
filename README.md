# GZIP package for the dliver system

GZIP compression packages for the different frameworks we use in dliver.

## Installation

1. Add dependency to go mod
2. Run go build/run/tidy

```bash
go get -u gitlab.com/proemergotech/gzip v1.0.0
```

## Usage

### Gentleman

```go
    import (
        "gitlab.com/proemergotech/gzip/gentlemangzip"	
    )

    centrifugeClient := client.NewCentrifuge(
        gentleman.New().BaseURL(fmt.Sprintf("%v://%v:%v", cfg.CentrifugeScheme, cfg.CentrifugeHost, cfg.CentrifugePort)).
            Use(gentlemangzip.Response(log.GlobalLogger())).
            Use(gentlemantrace.Middleware(opentracing.GlobalTracer(), log.GlobalLogger())).
            Use(gentlemanlog.Middleware(log.GlobalLogger(), true, true)).
            Use(client.RestErrorMiddleware("centrifuge")).
            Use(
                gentlemanretry.Middleware(
                    gentlemanretry.BackoffTimeout(20*time.Second),
                    gentlemanretry.Logger(log.GlobalLogger()),
                    gentlemanretry.RequestTimeout(5*time.Second),
                ),
            ).
            Use(gentlemangzip.Request(log.GlobalLogger())),
)
```

## Documentation

Private repos don't show up on godoc.org so you have to run it locally.

```
godoc -http=":6060"
```

Then open http://localhost:6060/pkg/gitlab.com/proemergotech/gzip/

## Development

- install go
- check out project to: $GOPATH/src/gitlab.com/proemergotech/gzip
