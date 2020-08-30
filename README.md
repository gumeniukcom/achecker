achekcer
=======

Simple domain checker with kafka and pg

# Specification

## JSONRPC 2.0

Use jsonrpc2.0 for rpc 
Full specification: https://www.jsonrpc.org/specification

### example message for check domain

```json
{
	"jsonrpc": "2.0",
	"method": "check_domain",
	"params": {
		"domain": "gumeniuk.com"
	},
	"id": "fa70dd9a-418a-4d38-86a0-a35768303bba"
}
```

### example message for save check domain
```json
{
	"jsonrpc": "2.0",
	"method": "save_check_domain",
	"params": {
		"domain": "http://gumeniuk.com",
		"status_code": 301,
		"error": ""
	},
	"id": "077bc49e-4685-4107-84f7-dde38e619fe5"
}
```

## easyjson

easyjson for fast marshall/unmarshall json
See https://github.com/mailru/easyjson

## go-critic

go-critic for check code
See https://github.com/go-critic/go-critic

# Before start

## Kafka

1. create kafka instance
2. create topic for checks task (eg "checks")
3. create topic for result task (eg "result")
4. remember / write config

## Postgres

1. create instance of pg
2. create user, db
3. Write config/ rember
4. With config, run  ``` go run cmd/add_tables.go```

## Config 

Write config : `config.toml`

eg
```toml
env = "local" # dev | prod
appName = "achecker"
version = "v1"

[logger]
    level = 0 # debug - 0, info - 1, warn - 2, error - 3, fatal - 4, panic - 5

[kafka]
    debug = true
    brokers = ["gumeniuk.com:443"]
    ssl = "true"
    capath = "./kafkaca/ca.pem"
    certpath = "./kafkaca/service.cert"
    keypath = "./kafkaca/service.key"
    group = "achecker"
    initial_offset = "oldest"
    version = "2.6.0"

[resultkafka]
    debug = true
    brokers = ["gumeniuk.com:443"]
    ssl = "true"
    capath = "./kafkaca/ca.pem"
    certpath = "./kafkaca/service.cert"
    keypath = "./kafkaca/service.key"
    group = "achecker"
    initial_offset = "oldest"
    version = "2.6.0"

[postgresql]
connect_string = "postgres://gumeniukcom:supersecurepass@gumeniuk.com:443/defaultdb?sslmode=require"
```

# Build && Run

## sh
```shell script
make
```

```shell script
./achecker 
```

## docker

1. put `ca`, `cert` and `keys` to `./kafkaca`
2. write correct `config.toml` (like eg upper)
3. build ```docker build -t achecker .```
4. run it ```docker run achecker```


# TODO

1. Check result should store check time (i forgot about it)
2. More interfaces
3. Add tests. most project without tests
4. Add migrations
5. Add table `domain`, and change `checks` table to reference to `domain`
6. Better design for tables