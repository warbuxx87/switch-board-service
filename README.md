# Switchboard Service

> Go service that controls which exchanges, trade pairs and service keys can access a distributed market-data network. Echo · GORM/MySQL · JWT · Swagger.

<p>
<img src="https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
<img src="https://img.shields.io/badge/MySQL-4479A1?style=flat-square&logo=mysql&logoColor=white" alt="MySQL" />
<img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker" />
<img src="https://img.shields.io/badge/Swagger-85EA2D?style=flat-square&logo=swagger&logoColor=black" alt="Swagger" />
<img src="https://img.shields.io/badge/Sentry-362D59?style=flat-square&logo=sentry&logoColor=white" alt="Sentry" />
</p>

The authorization hub of the platform. It keeps the whitelist of exchanges and trade pairs and issues the service keys that the data relays, discovery service and price-index service use to authenticate. It also includes a small admin UI (login and dashboard templates) for operators.

## Where it fits

Part of Alkemi's distributed market-data platform: a set of Go microservices that collect completed trades from exchange containers, relay them to oracles and a VWAP price index, and settle results on-chain.

```mermaid
flowchart LR
	EX[Exchange containers] -->|completed trades| DR[Data Relay Service]
	SB[Switchboard Service] -. whitelist + service keys .-> DR
	SB -. service keys .-> PI
	DS[Discovery Service] -. relay registry .-> OR
	DR -->|REST / WebSocket| PI[Price Index Service<br/>VWAP]
	DR -->|REST / WebSocket| OR[Oracles / Book-keeper nodes]
	OR -->|settlement votes| SC[(Settlement contract)]
	ADM[Switchboard Dashboard] --> SB
```

## How it works

- `cmd/switch-board-service` is the entry point. It runs an Echo HTTP server with session middleware.
- `internal/app/switch-board-service/models` defines the `User`, `Exchange`, `TradePair` and `ServiceKey` models, persisted to MySQL through GORM with auto-migration.
- Endpoints use JWT and API-key auth, and Swagger docs are generated into `docs/swagger`.
- Logging goes through logrus, with error reporting to Sentry.

## Stack

Go, MySQL, Docker, Swagger, Sentry

## Running locally

```bash
docker-compose up -d database          # MySQL 5.7
cp configs/config.json.example configs/config.json
go run ./cmd/switch-board-service
```

## Credits

Built at Alkemi by Raynaldo Rivera and Lekan Adigun (2019), as part of the platform I founded and led.

---
<sub>Part of <a href="https://github.com/warbuxx87">Ryan Breen's</a> work · <a href="https://www.linkedin.com/in/breen-ryan/">LinkedIn</a></sub>
