# postgres-alerter

This is a tool to alert about Postgres being down. Alerts are currently sent to Slack only

## Running

Just build the tool using `make` and then run it

```bash
make
```

```bash
./postgres-alerter
```

`postgres-alerter` runs as a long running service - it keeps running and does check every 1 second
