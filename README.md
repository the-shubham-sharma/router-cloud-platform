# Router Cloud Platform

> A production-grade, cloud-native backend platform for remotely managing, monitoring, and processing telemetry from thousands of simulated network devices — built with Go.

Built to demonstrate scalable backend engineering, distributed systems thinking, and real-world infrastructure patterns relevant to ISP and networking automation companies.

**Backend:** https://github.com/the-shubham-sharma/router-cloud-platform  
**Frontend:** https://github.com/the-shubham-sharma/rcp-frontend

---

## Screenshots

### Dashboard Overview
![Dashboard](screenshots/dashboard.png)

### Device Table
![Devices](screenshots/devices.png)

### API Endpoints
![API](screenshots/api.png)

### Architecture
![Architecture](screenshots/architecture.png)

### Grafana Metrics
![Grafana](screenshots/grafana.png)

### RabbitMQ Management
![RabbitMQ](screenshots/rabbitmq.png)

---

## Architecture

```
Device Simulator (Goroutines)
         │
         ▼
  Go + Gin REST API
  ├── JWT Authentication
  ├── Rate Limiting (10 req/s)
  ├── CORS Middleware
  └── Prometheus Metrics
         │
         ▼
   RabbitMQ Queue
   (heartbeat_queue)
         │
         ▼
  Consumer Workers (retry x3)
         │
    ┌────┴────┐
    ▼         ▼
PostgreSQL   Redis
(persistent) (2min TTL cache)
         │
         ▼
  WebSocket Hub
  (live broadcasts)
         │
         ▼
  Grafana + Prometheus
  (observability)
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
| Frontend | React + TypeScript + Vite + Tailwind |
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
- Retry logic (3 attempts with exponential backoff)
- Prometheus metrics (`/metrics` endpoint)
- Grafana dashboard
- Offline device alert detector (checks every 60s)
- Role-based access control (admin / user)
- Admin APIs (all devices, all users, promote user)
- React + TypeScript frontend dashboard

---

## API Reference

### Authentication
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/auth/register` | Register user + JWT | No |
| POST | `/auth/login` | Login + JWT | No |
| GET | `/auth/profile` | Get profile | Yes |

### Devices
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/devices` | Register device | Yes |
| GET | `/devices` | List devices (Redis cached) | Yes |
| GET | `/devices/:id` | Get device | Yes |
| PUT | `/devices/:id` | Update device | Yes |
| DELETE | `/devices/:id` | Delete device | Yes |
| POST | `/devices/:id/heartbeat` | Send heartbeat | Yes |
| GET | `/devices/:id/heartbeat` | Latest heartbeat (cached) | Yes |

### Dashboard
| Method | Endpoint | Description | Auth |
|---|---|---|---|
| GET | `/dashboard/summary` | Aggregated stats | Yes |

### Admin (Role: admin only)
| Method | Endpoint | Description |
|---|---|---|
| GET | `/admin/devices` | All devices across users |
| GET | `/admin/users` | All users |
| PUT | `/admin/users/:id/promote` | Promote user to admin |

### System
| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/ws?token=JWT` | WebSocket live events |

---

## Project Structure

```
router-cloud-platform/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, routes, graceful shutdown
├── internal/
│   ├── alert/
│   │   └── detector.go          # Offline device detection (60s ticker)
│   ├── cache/
│   │   └── redis.go             # Redis connection
│   ├── config/
│   │   └── config.go            # Environment config loader
│   ├── database/
│   │   └── database.go          # PostgreSQL + GORM auto-migrate
│   ├── handlers/
│   │   ├── admin.go             # Admin-only handlers
│   │   ├── auth.go              # Auth handlers
│   │   ├── dashboard.go         # Dashboard summary
│   │   ├── device.go            # Device CRUD
│   │   └── heartbeat.go         # Heartbeat ingestion
│   ├── metrics/
│   │   └── prometheus.go        # Custom Prometheus counters/gauges
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   ├── prometheus.go        # HTTP metrics middleware
│   │   ├── ratelimiter.go       # IP-based rate limiting
│   │   └── rbac.go              # Role-based access control
│   ├── models/
│   │   ├── device.go            # Device model
│   │   ├── heartbeat.go         # Heartbeat model
│   │   ├── metric.go            # Metric model
│   │   └── user.go              # User model with roles
│   ├── queue/
│   │   ├── consumer.go          # RabbitMQ consumer + retry logic
│   │   └── rabbitmq.go          # RabbitMQ connection + publish
│   ├── utils/
│   │   ├── jwt.go               # JWT generate + validate
│   │   └── response.go          # Standard API response wrapper
│   ├── websocket/
│   │   ├── handler.go           # WebSocket upgrade + auth
│   │   └── hub.go               # Hub pattern + broadcast
│   └── worker/
│       └── heartbeat_worker.go  # Goroutine worker pool
├── docker/
│   └── prometheus.yml           # Prometheus scrape config
├── docker-compose.yml           # All 5 infrastructure services
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

### 3. Start all infrastructure
```bash
docker compose up -d
```

This starts 5 containers:
| Container | Service | Port |
|---|---|---|
| rcp_postgres | PostgreSQL | 5432 |
| rcp_redis | Redis | 6379 |
| rcp_rabbitmq | RabbitMQ | 5672, 15672 |
| rcp_prometheus | Prometheus | 9090 |
| rcp_grafana | Grafana | 3001 |

### 4. Run the server
```bash
go run cmd/server/main.go
```

Expected output:
```
Config loaded successfully
Database connected successfully
Database migrated successfully
Redis connected successfully
RabbitMQ connected successfully
RabbitMQ consumer started
Heartbeat worker pool started with 5 workers
Offline device detector started
Server starting on port 8080
```

### 5. Run the frontend
```bash
git clone https://github.com/the-shubham-sharma/rcp-frontend.git
cd rcp-frontend
npm install
npm run dev
```

Open `http://localhost:5173`

---

## Monitoring

### Prometheus
Open `http://localhost:9090`

Custom metrics:
| Metric | Type | Description |
|---|---|---|
| `rcp_heartbeats_total` | Counter | Total heartbeats received |
| `rcp_active_devices` | Gauge | Currently online devices |
| `rcp_http_requests_total` | Counter | Requests by method/path/status |
| `rcp_http_request_duration_seconds` | Histogram | Request latency |
| `rcp_rabbitmq_messages_total` | Counter | Messages published to queue |
| `rcp_worker_jobs_total` | Counter | Jobs processed by worker pool |

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

- **Go** — goroutines, channels, worker pools, interfaces, error handling
- **Backend Engineering** — REST APIs, WebSockets, middleware chain, clean architecture
- **Distributed Systems** — message queues, async processing, caching strategies, event-driven design
- **DevOps** — Docker Compose, multi-container orchestration, health checks
- **Observability** — Prometheus metrics, Grafana dashboards, custom instrumentation
- **Security** — JWT auth, RBAC, rate limiting, bcrypt password hashing
- **Database** — PostgreSQL, GORM, auto-migration, foreign keys, UUID primary keys
- **Reliability** — retry logic, graceful shutdown, offline detection, fallback strategies

---

## Roadmap

### V3 (Planned)
- Kubernetes deployment with HPA
- API Gateway (Kong/Traefik)
- Distributed tracing (Jaeger/OpenTelemetry)
- AI anomaly detection on telemetry data
- Auto-healing simulation
- Firmware update simulation
- CI/CD pipeline (GitHub Actions)
- Multi-region support

---

## Author

**Shubham Sharma**  
GitHub: [@the-shubham-sharma](https://github.com/the-shubham-sharma)
