# Router Cloud Platform

A production-grade, cloud-native backend platform for remotely managing, monitoring, and processing telemetry from thousands of simulated network devices — built with Go.

> Built to demonstrate scalable backend engineering, distributed systems thinking, and real-world infrastructure patterns relevant to ISP and networking automation companies.

---

## Architecture

```
Device Simulator (Goroutines)
         │
         ▼
  Go + Gin REST API
  ├── JWT Authentication
  ├── Rate Limiting
  ├── CORS Middleware
  └── Prometheus Metrics
         │
         ▼
   RabbitMQ Queue
   (heartbeat_queue)
         │
         ▼
  Consumer Workers
         │
    ┌────┴────┐
    ▼         ▼
PostgreSQL   Redis
(persistent) (cache)
         │
         ▼
  WebSocket Hub
  (live device updates)
         │
         ▼
  Grafana Dashboard
  (Prometheus metrics)
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.21 |
| Framework | Gin |
| Database | PostgreSQL 15 |
| Cache | Redis 7 |
| Message Queue | RabbitMQ 3 |
| Metrics | Prometheus + Grafana |
| Auth | JWT (golang-jwt/jwt) |
| ORM | GORM |
| Real-time | WebSocket (gorilla/websocket) |
| Containerization | Docker Compose |

---

## Features

### V1 — Core Backend
- JWT Authentication (register, login, profile)
- Device Management CRUD APIs
- Heartbeat ingestion with Redis caching
- WebSocket live device status updates
- Dashboard summary API
- PostgreSQL with auto-migration
- Docker Compose for all infrastructure

### V2 — Production Patterns
- CORS middleware
- Graceful shutdown (5s drain)
- Rate limiting (10 req/s per IP, burst 20)
- Worker pool (5 goroutines, queue of 100)
- RabbitMQ message queue for heartbeat processing
- Retry logic (3 attempts with backoff)
- Prometheus metrics (`/metrics` endpoint)
- Grafana dashboard
- Offline device alert detector (checks every 60s)
- Role-based access control (admin / user)
- Admin APIs (all devices, all users, promote user)

---

## API Reference

### Authentication
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/auth/register` | Register user | No |
| POST | `/auth/login` | Login + JWT | No |
| GET | `/auth/profile` | Get profile | Yes |

### Devices
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/devices` | Register device | Yes |
| GET | `/devices` | List devices (cached) | Yes |
| GET | `/devices/:id` | Get device | Yes |
| PUT | `/devices/:id` | Update device | Yes |
| DELETE | `/devices/:id` | Delete device | Yes |
| POST | `/devices/:id/heartbeat` | Send heartbeat | Yes |
| GET | `/devices/:id/heartbeat` | Latest heartbeat | Yes |

### Dashboard
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| GET | `/dashboard/summary` | Aggregated stats | Yes |

### Admin
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| GET | `/admin/devices` | All devices | Admin |
| GET | `/admin/users` | All users | Admin |
| PUT | `/admin/users/:id/promote` | Promote to admin | Admin |

### System
| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/ws` | WebSocket (token param) |

---

## Project Structure

```
router-cloud-platform/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── alert/
│   │   └── detector.go          # Offline device detection
│   ├── cache/
│   │   └── redis.go             # Redis connection
│   ├── config/
│   │   └── config.go            # Environment config
│   ├── database/
│   │   └── database.go          # PostgreSQL + auto-migrate
│   ├── handlers/
│   │   ├── auth.go              # Auth handlers
│   │   ├── admin.go             # Admin handlers
│   │   ├── dashboard.go         # Dashboard handlers
│   │   ├── device.go            # Device handlers
│   │   └── heartbeat.go         # Heartbeat handlers
│   ├── metrics/
│   │   └── prometheus.go        # Custom Prometheus metrics
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   ├── prometheus.go        # HTTP metrics middleware
│   │   ├── ratelimiter.go       # Rate limiting
│   │   └── rbac.go              # Role-based access control
│   ├── models/
│   │   ├── device.go            # Device model
│   │   ├── heartbeat.go         # Heartbeat model
│   │   ├── metric.go            # Metric model
│   │   └── user.go              # User model with roles
│   ├── queue/
│   │   ├── consumer.go          # RabbitMQ consumer
│   │   └── rabbitmq.go          # RabbitMQ connection + publish
│   ├── utils/
│   │   ├── jwt.go               # JWT generate + validate
│   │   └── response.go          # Standard API response
│   ├── websocket/
│   │   ├── handler.go           # WebSocket handler
│   │   └── hub.go               # WebSocket hub + broadcast
│   └── worker/
│       └── heartbeat_worker.go  # Worker pool
├── docker/
│   └── prometheus.yml           # Prometheus scrape config
├── docker-compose.yml           # All infrastructure
├── .env.example                 # Environment variables template
└── README.md
```

---

## Getting Started

### Prerequisites
- Go 1.21+
- Docker Desktop

### 1. Clone the repo
```bash
git clone https://github.com/the-shubham-sharma/router-cloud-platform.git
cd router-cloud-platform
```

### 2. Set up environment
```bash
cp .env.example .env
```

### 3. Start infrastructure
```bash
docker compose up -d
```

This starts:
- PostgreSQL on `:5432`
- Redis on `:6379`
- RabbitMQ on `:5672` (management UI on `:15672`)
- Prometheus on `:9090`
- Grafana on `:3001`

### 4. Run the server
```bash
go run cmd/server/main.go
```

### 5. Test the API
```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Shubham","email":"shubham@test.com","password":"123456"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"shubham@test.com","password":"123456"}'
```

---

## Monitoring

### Prometheus
Open `http://localhost:9090`

Custom metrics exposed:
- `rcp_heartbeats_total` — total heartbeats received
- `rcp_active_devices` — currently online devices
- `rcp_http_requests_total` — requests by method/path/status
- `rcp_http_request_duration_seconds` — request latency histogram
- `rcp_rabbitmq_messages_total` — messages published to queue
- `rcp_worker_jobs_total` — jobs processed by worker pool

### Grafana
Open `http://localhost:3001` (admin/admin)

Add Prometheus data source: `http://rcp_prometheus:9090`

### RabbitMQ Management
Open `http://localhost:15672` (rcpuser/rcppassword)

---

## Environment Variables

```env
SERVER_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=rcpuser
DB_PASSWORD=rcppassword
DB_NAME=rcpdb

REDIS_HOST=localhost
REDIS_PORT=6379

RABBITMQ_USER=rcpuser
RABBITMQ_PASSWORD=rcppassword
RABBITMQ_URL=amqp://rcpuser:rcppassword@localhost:5672/

JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=24
```

---

## Skills Demonstrated

- **Go** — goroutines, channels, worker pools, interfaces
- **Backend Engineering** — REST APIs, WebSockets, middleware chain
- **Distributed Systems** — message queues, async processing, caching strategies
- **DevOps** — Docker Compose, multi-container orchestration
- **Observability** — Prometheus metrics, Grafana dashboards
- **Security** — JWT auth, RBAC, rate limiting
- **Database** — PostgreSQL, GORM, auto-migration, foreign keys
- **Reliability** — retry logic, graceful shutdown, offline detection

---

## Roadmap

### V3 (Planned)
- Kubernetes deployment
- API Gateway
- Distributed tracing (Jaeger)
- AI anomaly detection on telemetry
- Auto-healing simulation
- Firmware update simulation
- CI/CD pipeline
- Multi-region support
