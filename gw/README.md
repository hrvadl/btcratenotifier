# Gateway microservice (gw)

This service is the main entrypoint (and only public avaiable) service to the application. Its role is map HTTP -> GRPC requests with some extra logic.

## Available tasks

You can see all available tasks running following command in the root of the repo:

```sh
task
```

You should get a following output:

```sh
task: [default] task --list-all
task: Available tasks for this project:
* default:                Show available tasks
* generate:               Generate (used for mock generation)
* install:                Install all tools
* lint:                   Run golangci-lint
* run:                    Populate env from .env file and run service
* run-with-env:           run service
* test:                   Run tests
* gen:swagger:            Generate swagger docs
* install:godotenv:       Install go dot env lib
* install:gofumpt:        Install gofumpt
* install:lint:           Install golangci-lint
* install:mock:           Install mockgen
* test:cover:             Run tests & show coverage
* test:race:              Run tests with a race flag

```

## How to run?

If you want to run it as a standalone service you need:

1. Populate env vars needed for it in root `.env` file (../.env)
2. Run `task run` from `./gw` dir or `task gw:run` from the root of the repo

## Documentation

You should be able to hit `<BASE_URL>/docs/index.html` and observe swagger API docs.

<img width="1693" alt="image" src="https://github.com/hrvadl/converter/assets/93580374/411b8b23-fc7c-4ab4-8da3-3b1c246196ac">

## Folder structure

1. `pkg` contains possibly reusable package, not binded to this project. Currently it contains only logger utils
2. `internal`contains packages binded to this project.
   2.1. `cfg` contains config which is read from environment vars.
   2.2. `app` is an abstraction with all services initialization
   2.3. `transport` contains all transport layer logic: grpc clients and http handlers.
3. `docs` container swagger-generated API documentation.
4. `cmd` contains entrypoints to the program.
