# 💰 Expenses Tracker API

> A personal-finance REST API written in Go — wallets, deposits, withdrawals, peer-to-peer transfers and a fully auditable transaction ledger.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-1.12-00B5A0?style=flat-square)](https://gin-gonic.com)
[![GORM](https://img.shields.io/badge/GORM-1.31-CC2927?style=flat-square)](https://gorm.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![Lint](https://img.shields.io/badge/lint-golangci--lint-7B68EE?style=flat-square)](https://golangci-lint.run)
[![Status](https://img.shields.io/badge/status-in%20development-F5A623?style=flat-square)](#-roadmap)

---

## 📌 Table of Contents

- [Overview](#-overview)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Configuration](#-configuration)
- [Make Targets](#-make-targets)
- [API Reference](#-api-reference)
- [Data Model](#-data-model)
- [Money Handling](#-money-handling)
- [Development Workflow](#-development-workflow)
- [Roadmap](#-roadmap)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 Overview

**Expenses Tracker API** gives every user a wallet and an immutable record of everything
that moves through it. Instead of storing a balance that gets overwritten, the service
records **transactions** — the balance is the story those transactions tell.

| Domain | What it handles |
| :--- | :--- |
| 🔐 **Auth** | Registration, login, stateless JWT access + refresh tokens |
| 👛 **Wallet** | One wallet per user, current balance, deposits and withdrawals |
| 🧾 **Transactions** | Append-only ledger — filterable by type, paginated, never mutated |
| 🔁 **Transfers** | Atomic peer-to-peer transfers: debit and credit succeed or fail together |

**Engineering highlights**

- 🧱 **Layered architecture** — handlers, services and repositories are separated, so the
  business rules can be unit-tested without an HTTP server or a database.
- ⏱️ **Hardened HTTP server** — explicit read, write and idle timeouts guard against slow
  clients and connection exhaustion.
- 📋 **Structured logging** — `zerolog`, pretty console output in dev, JSON in `release`.
- ✅ **Validation at the edge** — `binding` tags reject malformed payloads before any
  business logic runs.
- 🔒 **Transactional integrity** — money movements run inside a database transaction, so a
  transfer can never leave one side credited and the other untouched.
- ⚡ **Hot reload** with `air`, **strict linting** with `golangci-lint`.

---

## 🧰 Tech Stack

| Layer | Technology | Why |
| :--- | :--- | :--- |
| Language | **Go 1.26** | Static typing, fast builds, single-binary deployment |
| HTTP | **Gin 1.12** | High-performance router with a mature middleware ecosystem |
| ORM | **GORM 1.31** + `driver/postgres` | Relations, hooks and transactions with little ceremony |
| Database | **PostgreSQL 12** | ACID guarantees — non-negotiable for financial data |
| Auth | **JWT** | Stateless access + refresh token pair |
| Logging | **rs/zerolog** | Zero-allocation structured logging |
| Config | **joho/godotenv** | `.env` loading with sane in-code defaults |
| Migrations | **golang-migrate** | Versioned, reversible SQL migrations |
| Live reload | **air** | Sub-second feedback loop in development |
| Quality | **golangci-lint**, `gofmt`, `goimports` | Consistent, vetted codebase |

---

## 🏗 Architecture

```mermaid
flowchart TD
    Client([Client])
    Client -->|HTTP/JSON| Router

    subgraph API["cmd/api — HTTP server"]
        Router["Gin Engine"]
        MW["Middleware<br/>Logger · Recovery · CORS · Auth"]
        Router --> MW
    end

    subgraph Internal["internal/"]
        Handlers["handlers<br/>bind · validate · respond"]
        Services["services<br/>wallet · transfer rules"]
        Repos["repositories<br/>GORM queries"]
        Models["models<br/>User · Wallet · Transaction"]
        Utils["utils<br/>jwt · env · response"]
    end

    MW --> Handlers
    Handlers --> Services
    Handlers --> Utils
    Services --> Repos
    Repos --> Models
    Repos -->|"GORM (tx)"| DB[(PostgreSQL)]

    Config["config<br/>.env → typed structs"] -.-> API
    Config -.-> Internal
    Logger["logger<br/>zerolog"] -.-> API
    Logger -.-> Internal
```

**A transfer, step by step**

```mermaid
sequenceDiagram
    participant C as Client
    participant H as TransferHandler
    participant S as TransferService
    participant DB as PostgreSQL

    C->>H: POST /api/v1/transfers
    H->>H: bind + validate payload
    H->>S: Transfer(from, to, amount)
    S->>DB: BEGIN
    S->>DB: SELECT wallet FOR UPDATE (sender)
    S->>S: assert sufficient balance
    S->>DB: debit sender + INSERT transaction
    S->>DB: credit receiver + INSERT transaction
    S->>DB: COMMIT
    S-->>H: transfer record
    H-->>C: 201 Created
```

---

## 📁 Project Structure

```
expense-tracker/
├── cmd/
│   └── api/
│       └── main.go            # Entrypoint: config → logger → db → server
├── internal/
│   ├── config/
│   │   └── config.go          # Env → typed config (Server, Database, JWT)
│   ├── database/
│   │   └── database.go        # GORM + PostgreSQL connection factory
│   ├── logger/
│   │   └── logger.go          # zerolog setup (console in dev, JSON in release)
│   ├── server/
│   │   └── server.go          # Engine, middleware, routes, HTTP timeouts
│   └── utils/
│       └── env.go             # Environment helpers with defaults
├── db/
│   └── migrations/            # golang-migrate SQL files (NNNN_name.up/down.sql)
├── docker/
│   └── docker-compose.yml     # PostgreSQL + LocalStack (S3, SQS)
├── .air.toml                  # Hot-reload configuration
├── .golangci.yml              # Linter configuration
├── Makefile                   # Developer commands
├── go.mod / go.sum
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

| Tool | Version | Install |
| :--- | :--- | :--- |
| Go | 1.26+ | https://go.dev/dl |
| Docker + Compose | latest | https://docs.docker.com/get-docker |
| `golang-migrate` | latest | `brew install golang-migrate` |
| `air` *(optional)* | latest | `go install github.com/air-verse/air@latest` |
| `golangci-lint` *(optional)* | latest | `brew install golangci-lint` |

### 1 — Clone

```bash
git clone git@github.com:boriskamtou96/Expenses-Tracker-API.git
cd Expenses-Tracker-API
```

### 2 — Install dependencies

```bash
go mod download
```

### 3 — Create your environment file

```bash
cp .env.example .env   # or create .env from the table below
```

### 4 — Start PostgreSQL

```bash
make up          # PostgreSQL on :5432, LocalStack on :4566
```

Create the database:

```bash
docker exec -it postgres createdb --username=postgres --owner=postgres expense_tracker
```

### 5 — Run the migrations

```bash
set -a && source .env && set +a   # export DB_* for the Makefile
make migrate-up
```

### 6 — Run the API

```bash
make run         # go run ./cmd/api
# or
make air         # hot reload
```

### 7 — Verify

```bash
curl -s http://localhost:8080/health
```

```json
{ "app": "expense tracker", "status": "Ok" }
```

🎉 The API is live on **http://localhost:8080**.

---

## ⚙️ Configuration

All configuration is read from the environment; `.env` is loaded automatically at startup
and every key falls back to a safe default.

### Server

| Variable | Default | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | Port the HTTP server binds to |
| `GIN_MODE` | `debug` | `debug` or `release` — also switches the logger to JSON |

### Database

| Variable | Default | Description |
| :--- | :--- | :--- |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `expense_tracker` | Database name |
| `DB_SSLMODE` | `disable` | `disable` locally, `require` in production |

### Authentication

| Variable | Default | Description |
| :--- | :--- | :--- |
| `JWT_SECRET` | `your_jwt_secret` | HMAC signing secret — **must** be overridden in production |
| `JWT_EXPIRES_IN` | `15m` | Access-token lifetime (Go duration: `15m`, `1h`, `24h`) |
| `REFRESH_TOKEN_EXPIRES_IN` | `168h` | Refresh-token lifetime (Go durations have no `d` unit — use `168h` for 7 days) |

<details>
<summary>📄 <b>Sample <code>.env</code></b></summary>

```dotenv
# Server
PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=expense_tracker
DB_SSLMODE=disable
DB_ADDR=postgres://postgres:postgres@localhost:5432/expense_tracker?sslmode=disable

# Auth
JWT_SECRET=change-me-to-a-long-random-string
JWT_EXPIRES_IN=15m
REFRESH_TOKEN_EXPIRES_IN=168h
```

> ⚠️ `.env` is git-ignored. Never commit real secrets.

</details>

---

## 🛠 Make Targets

```bash
make help                    # List every available target
```

| Target | Description |
| :--- | :--- |
| `make up` | Start PostgreSQL + LocalStack via Docker Compose |
| `make down` | Stop and remove the containers |
| `make migration <name>` | Scaffold a new sequential migration pair |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down <n>` | Roll back the last `n` migrations |
| `make build` | Compile the binary to `bin/app` |
| `make run` | Run the API (`go run ./cmd/api`) |
| `make air` | Run with hot reload |
| `make lint` | Run `golangci-lint` across the module |
| `make fix-lint` | Run the linter with `--fix` |
| `make format` | `gofmt -s -w .` + `goimports -w .` |

> ℹ️ Migration targets build the DSN from `DB_*`, so export your `.env` first:
> `set -a && source .env && set +a`

---

## 📡 API Reference

Base URL: `http://localhost:8080`

### System

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | — | Liveness probe |

### Authentication

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | — | Create an account and its wallet |
| `POST` | `/api/v1/auth/login` | — | Obtain an access + refresh token pair |

### Wallet

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/wallet` | 🔒 | Current user's wallet and balance |
| `POST` | `/api/v1/wallet/deposit` | 🔒 | Credit the wallet |
| `POST` | `/api/v1/wallet/withdraw` | 🔒 | Debit the wallet (balance must cover it) |

### Transactions

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/transactions` | 🔒 | List the ledger |
| `GET` | `/api/v1/transactions/:id` | 🔒 | Get one transaction |

**Query parameters**

| Parameter | Example | Description |
| :--- | :--- | :--- |
| `type` | `?type=withdrawal` | Filter by `deposit`, `withdrawal`, `transfer_in`, `transfer_out` |
| `page` | `?page=1` | Page number (default `1`) |
| `limit` | `?limit=20` | Items per page (default `20`) |
| `from` / `to` | `?from=2026-01-01&to=2026-01-31` | Date range |

### Transfers

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/transfers` | 🔒 | Move funds to another user's wallet |

**Legend** — `—` public · 🔒 authenticated

> **Status.** `/health` is live today; the routes above are registered in
> `internal/server/server.go` as the contract this API is being built against, and their
> handlers are landing next. See the [roadmap](#-roadmap).

### Authenticating a request

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"supersecret"}'
```

```bash
curl http://localhost:8080/api/v1/wallet \
  -H "Authorization: Bearer <accessToken>"
```

<details>
<summary>📥 <b>Example payloads</b></summary>

**Register**

```json
{
  "email": "jane@example.com",
  "password": "supersecret",
  "first_name": "Jane",
  "last_name": "Doe"
}
```

**Deposit**

```json
{ "amount": "150.00", "description": "January salary" }
```

**Withdraw**

```json
{ "amount": "42.50", "description": "Groceries" }
```

**Transfer**

```json
{ "recipient_email": "john@example.com", "amount": "25.00", "description": "Dinner split" }
```

**Wallet response**

```json
{
  "success": true,
  "message": "Wallet retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "balance": "1082.50",
    "currency": "XAF"
  }
}
```

**Paginated transactions**

```json
{
  "success": true,
  "message": "Transactions retrieved successfully",
  "data": [
    {
      "id": 42,
      "type": "withdrawal",
      "amount": "42.50",
      "balance_after": "1082.50",
      "description": "Groceries",
      "created_at": "2026-08-21T10:14:00Z"
    }
  ],
  "meta": { "page": 1, "limit": 20, "total": 137, "totalPages": 7 }
}
```

</details>

---

## 🗄 Data Model

```mermaid
erDiagram
    USER ||--|| WALLET : owns
    WALLET ||--o{ TRANSACTION : records
    USER ||--o{ TRANSFER : sends
    USER ||--o{ TRANSFER : receives
    TRANSFER ||--o{ TRANSACTION : produces

    USER {
        uint     id PK
        string   email UK
        string   password_hash
        string   first_name
        string   last_name
        bool     is_active
        datetime created_at
    }
    WALLET {
        uint     id PK
        uint     user_id FK
        decimal  balance
        string   currency
        datetime updated_at
    }
    TRANSACTION {
        uint     id PK
        uint     wallet_id FK
        string   type
        decimal  amount
        decimal  balance_after
        string   description
        string   reference UK
        datetime created_at
    }
    TRANSFER {
        uint     id PK
        uint     sender_id FK
        uint     recipient_id FK
        decimal  amount
        string   status
        datetime created_at
    }
```

**Transaction types**

| Type | Effect on balance | Created by |
| :--- | :--- | :--- |
| `deposit` | ➕ credit | `POST /wallet/deposit` |
| `withdrawal` | ➖ debit | `POST /wallet/withdraw` |
| `transfer_out` | ➖ debit | `POST /transfers` (sender side) |
| `transfer_in` | ➕ credit | `POST /transfers` (recipient side) |

> 💡 Transactions are **append-only**. Nothing is ever updated or deleted — a correction is
> a new, opposite entry. `balance_after` snapshots the wallet balance at that instant, which
> makes the ledger auditable and reconcilable without replaying history.

---

## 💵 Money Handling

Financial data punishes shortcuts, so this project follows three rules:

1. **Never `float64` for money.** Amounts are stored as PostgreSQL `NUMERIC(19,4)` and
   handled in Go as `shopspring/decimal` — binary floats cannot represent `0.10` exactly,
   and rounding drift in a ledger is a bug you find months later.
2. **Every mutation runs in a transaction.** A transfer debits and credits inside a single
   `BEGIN … COMMIT`; a failure at any step rolls the whole thing back.
3. **Row locks over read-modify-write.** Wallet rows are read with `SELECT … FOR UPDATE`, so
   two concurrent withdrawals cannot both read the same balance and overdraw the account.

---

## 👨‍💻 Development Workflow

### Hot reload

```bash
make air
```

`air` watches every `.go` file, rebuilds into `./bin/main` and restarts the process.
Build errors are written to `build-errors.log`.

### Creating a migration

```bash
make migration create_wallets_table
# → db/migrations/000002_create_wallets_table.up.sql
# → db/migrations/000002_create_wallets_table.down.sql
```

Write the SQL, then:

```bash
make migrate-up          # apply
make migrate-down 1      # roll back the last one
```

### Before every commit

```bash
make format
make lint
go test ./...
```

### Commit convention

```
feat(wallet): add withdrawal endpoint
fix(transfer): roll back on insufficient balance
chore(deps): bump gorm to 1.31.2
```

---

## 🗺 Roadmap

**Foundation** — done

- [x] Project layout (`cmd/` + `internal/`)
- [x] Typed configuration loaded from the environment
- [x] Structured logging with `zerolog`
- [x] PostgreSQL connection via GORM
- [x] Gin engine with Logger, Recovery and CORS middleware
- [x] HTTP server with read / write / idle timeouts
- [x] Health-check endpoint
- [x] Route contract for auth, wallet, transactions and transfers
- [x] Docker Compose stack (PostgreSQL + LocalStack)
- [x] `Makefile`, `air` hot reload, `golangci-lint`

**In progress**

- [ ] GORM models (`User`, `Wallet`, `Transaction`, `Transfer`) and the first migration set
- [ ] Repository layer
- [ ] Auth service — bcrypt hashing, register, login, refresh
- [ ] JWT authentication middleware
- [ ] Wallet service — deposit, withdraw, balance
- [ ] Transaction listing with filters and pagination
- [ ] Transfer service with row locking and rollback

**Next**

- [ ] Graceful shutdown on `SIGINT` / `SIGTERM`
- [ ] Idempotency keys on money-moving endpoints
- [ ] Spending categories, budgets and monthly reports
- [ ] Unit tests on services, integration tests on handlers (testcontainers)
- [ ] OpenAPI / Swagger documentation
- [ ] Rate limiting and request-ID middleware
- [ ] Multi-stage `Dockerfile` and GitHub Actions CI
- [ ] Receipt attachments stored on S3
- [ ] CSV / PDF statement export

---

## 🤝 Contributing

1. Fork the repository
2. Create a branch — `git checkout -b feat/my-feature`
3. Run `make format && make lint && go test ./...`
4. Commit — `git commit -m "feat: add my feature"`
5. Push and open a Pull Request

---

## 📄 License

Released under the **MIT License**.

---

## 👤 Author

**Boris Kamtou** — [@boriskamtou96](https://github.com/boriskamtou96)

<p align="center">
  <sub>Built with Go 💙 — if this project helped you, consider leaving a ⭐</sub>
</p>
