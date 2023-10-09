# adctf_scoreboard_exporter

Prometheus exporter for Attack Defense CTF scoreboards.

Supported Attack Defense CTF scoreboard APIs:

* `faustv1` (old Faust CTF scoreboard API)
* `faustv2` (new Faust CTF scoreboard-v2 API)

## Running

**Option 1. Compile from source**

```shell
git clone git@github.com:boxmein/adctf_scoreboard_exporter.git
cd adctf_scoreboard_exporter
go build -o scoreboard_exporter ./cmd/scoreboard_exporter
./scoreboard_exporter --listenAddr :5001 faustv2 --base-url https://2023.faustctf.net
```

**Option 2. Run with Docker**

```shell
git clone git@github.com:boxmein/adctf_scoreboard_exporter.git
cd adctf_scoreboard_exporter
docker build -t scoreboard_exporter .
docker run --rm -it scoreboard_exporter -- --listenAddr :5001 faustv2 --base-url https://2023.faustctf.net
```

**Option 3. Run with Docker Compose**

```shell
git clone git@github.com:boxmein/adctf_scoreboard_exporter.git
cd adctf_scoreboard_exporter
docker compose build
docker compose up
```

## CLI flags

Help can be found on:

```
./scoreboard_exporter --help
./scoreboard_exporter faustv1 --help
./scoreboard_exporter faustv2 --help
```

Example to pull metrics from faustv2 API on 2023.faustctf.net:

```shell
./scoreboard_exporter --listenAddr :5001 faustv2 --base-url https://2023.faustctf.net
```

Example to pull metrics from faustv1 API on 2023.faustctf.net:

```shell
./scoreboard_exporter --listenAddr :5001 faustv1 --base-url https://2023.faustctf.net
```

Customizing API endpoints for faustv1:

```shell
./scoreboard_exporter --listenAddr :5001 faustv1 \
  --scoreboard-url https://2023.faustctf.net/competition/scoreboard.json \
  --status-url https://2023.faustctf.net/competition/status.json
```

Customizing API endpoints for faustv2:

```shell
./scoreboard_exporter --listenAddr :5001 faustv2 \
  --current-url https://2023.faustctf.net/competition/scoreboard-v2/scoreboard_current.json \
  --round-url https://2023.faustctf.net/competition/scoreboard-v2/scoreboard_round_%d.json \
  --teams-url https://2023.faustctf.net/competition/scoreboard-v2/scoreboard_teams.json
```


## Exported metrics

The metrics are also documented in the Prometheus HELP comments.

All metrics are annotated with the `{team, service}` labels to distinguish 
between each team's and service's specific points.

Metric                   | Example  | Meaning
-------------------------|----------|---
scoreboard_tick          | 100      | Current tick
scoreboard_offense       | 199.203  | Offense points
scoreboard_defense       | 402.1    | Defense points
scoreboard_sla           | 2502.22  | SLA points
scoreboard_captures      | 44       | Flags captured
scoreboard_stolen        | 95       | Flags lost / stolen

## Support matrix

Not all APIs support all the metrics.

Metric                   | faustv1  | faustv2
-------------------------|----------|----------
scoreboard_tick          | YES      | YES
scoreboard_offense       | YES      | YES
scoreboard_defense       | YES      | YES
scoreboard_sla           | YES      | YES
scoreboard_captures      | NO       | YES
scoreboard_stolen        | NO       | YES
