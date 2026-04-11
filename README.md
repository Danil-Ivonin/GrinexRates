# GrinexRates

A gRPC service that fetches USDT/RUB market rates from the Grinex exchange,
computes order book statistics (topN, avgNM), and persists every result to PostgreSQL.

## Configuration

All configuration can be supplied via CLI flags or environment variables and config file.
**CLI flags take precedence** over environment variables when both are set.

| Environment Variable | CLI Flag             | Default      | Required | Description                                      |
|---------------------|----------------------|--------------|----------|--------------------------------------------------|
| `DB_HOST`           | `--db-host`          | `postgres`   | No       | PostgreSQL server hostname or IP address         |
| `DB_PORT`           | `--db-port`          | `5432`       | No       | PostgreSQL server port                           |
| `DB_USER`           | `--db-user`          | `postgres`   | No       | PostgreSQL login user                            |
| `DB_PASSWORD`       | `--db-password`      | _(empty)_    | Yes*     | PostgreSQL login password                        |
| `DB_NAME`           | `--db-name`          | `usdt_rate`  | No       | PostgreSQL database name                         |
| `DB_SSLMODE`        | `--db-sslmode`       | `disable`    | No       | PostgreSQL SSL mode (`disable`, `require`, etc.) |
| `GRPC_PORT`         | `--grpc-port`        | `50051`      | No       | Port the gRPC server listens on                  |
| `PROMETHEUS_PORT`   | `--prometheus-port`  | `9090`       | No       | Metrics server port                              |

| Convig variable               | CLI Flag           | Default                                               | Description                    |
|-------------------------------|--------------------|-------------------------------------------------------|--------------------------------|
| `grinex.url`                  | `--grinex-url`     | `https://grinex.io/api/v1/spot/depth?symbol=usdta7a5` | `Grinex rates endpoint`        |
| `grinex.timeout`              | `--grinex-timeout` | `10s`                                                 | `Grinex HTTP request timeout`  |

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- `grpcurl` for testing (optional): https://github.com/fullstorydev/grpcurl

### 1. Clone and configure

The defaults in `.env` work out of the box with `docker-compose`. You only need to
edit `.env` if you want a custom password or database name.

### 2. Build the Docker image

```bash
make build        # compile the local binary (optional, for local dev)
make docker-build # build the Docker image used by docker-compose
```

### 3. Start the service

```bash
docker-compose up -d
```

Both containers start. The app waits for PostgreSQL to pass its health check before
connecting. Schema migrations run automatically on first start.

### 4. Test the gRPC endpoint

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Call GetRates (N=3, M=5)
grpcurl -plaintext -d '{"n": 3, "m": 5}' localhost:50051 rates.v1.RatesService/GetRates
```

Expected response fields: `ask_price`, `bid_price`, `top_n`, `avg_nm`, `fetched_at`.

### 5. Run unit tests

```bash
make test
```

### 6. Stop the service

```bash
docker-compose down
```

## gRPC API

The service exposes a single RPC:

```
service RatesService {
  rpc GetRates(GetRatesRequest) returns (GetRatesResponse);
}
```

- `GetRatesRequest`: `n` (int32, 1-indexed position for topN), `m` (int32, upper position for avgNM)
- `GetRatesResponse`: `ask_price`, `bid_price`, `top_n`, `avg_nm` (strings), `fetched_at` (int64 Unix seconds)
