## Conibase Tickers feed

### Use

Fast run
```
$ docker-compose -f docker/docker-compose.yaml  up
```
or

```
$ MYSQL_URL="root:example@/test" go run cmd/feeder.go -symbol ETH-BTC  -symbol BTC-USD
```

### How it work
- open connection to mysql
- create table
- open one websocket connection
- create gorutines every ticker-symbol
- write ticker into mysql table