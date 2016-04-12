# docker-ls

Lists all mapped ports of your running Docker containers. Containers listed as
hosting HTTP services will have clickable links.

## Usage

```sh
go build && sudo ./docker-ls [config.yml] &
```

`sudo` is required because `/var/run/docker.sock` is locked down.

## Config

All configuration is optional. If you don't supply a file, the defaults are used.

```yaml
# run the http server on this port (default: 80)
port: 80

# override the hostname used in the output (default: autodetect)
host: echo.sh

# omit any containers based on any version of these images (default: none)
# note that these are KEYS, not array elements
blacklist:
  docker.movio.co/zookeeper:

# prepend these protocols to these ports of these images (default: none)
protocols:
  grafana/grafana:
    3000: http
  prom/prometheus:
    9090: http
```

## Output (in HTML)

```md
## Containers

- **kafka** (docker.movio.co/kafka:2.0.1)
  - echo.sh:9092
- **grafana** (grafana/grafana)
  - [http://echo.sh:3000](http://echo.sh:3000)
- **prometheus** (prom/prometheus)
  - [http://echo.sh:9090](http://echo.sh:9090)
```
