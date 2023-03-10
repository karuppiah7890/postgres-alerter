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

`postgres-alerter` is just a simple tool and does not run as a service, it just runs once and then exits. to keep it running, run it as a cron job or using `watch` command continuously every few seconds or minutes, whatever interval you wish
