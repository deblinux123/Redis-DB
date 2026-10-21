# 14-Day Redis with Golang Roadmap

A hands-on roadmap for learning Redis with Go, starting with the fundamentals and progressing toward production-ready backend development.

The goal is to understand Redis, integrate it with Go applications, and build practical systems using caching, data structures, messaging, and concurrency.

## 🎯 Goals

By the end of these 14 days, you will be able to:

- Understand Redis architecture and core data structures.
- Connect Redis to Go applications.
- Perform CRUD operations using Redis.
- Implement caching and expiration strategies.
- Use Redis for sessions, rate limiting, and distributed locking.
- Work with Pub/Sub and Streams.
- Handle errors, timeouts, and connection management.
- Build a practical Redis-powered Go backend.

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **Database:** Redis
- **Redis Client:** `github.com/redis/go-redis/v9`
- **Containerization:** Docker
- **API Framework:** Fiber (optional)
- **Testing:** Go testing package
- **Configuration:** Environment variables and `.env`

## 📅 14-Day Learning Plan

### Day 01 — Redis Fundamentals

**Topics**
- What is Redis?
- In-memory data storage
- Redis vs. traditional databases
- Redis use cases
- Redis architecture and persistence overview

**Practical Tasks**
- Run Redis with Docker.
- Connect using `redis-cli`.
- Execute `PING`, `SET`, `GET`, and `DEL`.
- Inspect basic server information.

**Deliverable:** A running Redis instance and a collection of basic Redis commands.

### Day 02 — Redis Setup with Golang

**Topics**
- Redis Go client
- Creating a Redis client
- Connecting to Redis
- Context and timeouts
- Environment-based configuration

**Practical Tasks**
- Initialize a Go module.
- Install `go-redis/v9`.
- Connect Go to Redis.
- Implement a connection health check.
- Handle connection errors correctly.

**Deliverable:** A reusable Redis client package.

### Day 03 — Strings and Basic CRUD

**Topics**
- Redis strings
- `SET`, `GET`, `DEL`, and `EXISTS`
- `INCR` and `DECR`
- Key expiration with `EXPIRE` and `TTL`
- Atomic operations

**Practical Tasks**
- Store and retrieve user data.
- Implement counters.
- Set expiration times.
- Handle missing keys.

**Deliverable:** A Go application demonstrating Redis CRUD operations.

### Day 04 — Hashes and Structured Data

**Topics**
- Redis hashes
- `HSET`, `HGET`, `HGETALL`, and `HDEL`
- Representing objects in Redis
- Redis data modeling

**Practical Tasks**
- Store user profiles using hashes.
- Update individual fields.
- Retrieve complete profiles.
- Compare hashes with JSON stored in strings.

**Deliverable:** A simple user repository backed by Redis.

### Day 05 — Lists, Sets, and Sorted Sets

**Topics**
- Lists: `LPUSH`, `RPUSH`, `LRANGE`, and `LPOP`
- Sets: `SADD`, `SMEMBERS`, and `SISMEMBER`
- Sorted sets: `ZADD`, `ZRANGE`, and ranking
- Choosing the right data structure

**Practical Tasks**
- Build a simple task queue using lists.
- Track unique users with sets.
- Create a leaderboard using sorted sets.

**Deliverable:** A small application demonstrating three Redis data structures.

### Day 06 — Expiration and Cache Design

**Topics**
- TTL and key expiration
- Cache-aside pattern
- Cache hits and misses
- Cache invalidation
- Cache stampede basics

**Practical Tasks**
- Build a Go service that caches data.
- Read from Redis before querying the underlying data source.
- Set TTL values.
- Invalidate cached data after updates.

**Deliverable:** A reusable caching layer for a Go application.

### Day 07 — Redis and REST APIs with Fiber

**Topics**
- Redis integration with Fiber
- Request validation
- Repository and service layers
- HTTP error handling
- JSON responses

**Practical Tasks**
- Build a REST API for managing users or products.
- Implement create, read, update, and delete endpoints.
- Store records in Redis.
- Test endpoints using `curl`.

**Deliverable:** A Redis-backed Go REST API.

### Day 08 — Transactions and Pipelining

**Topics**
- Redis transactions with `MULTI` and `EXEC`
- `TxPipelined` and pipelining in Go
- Reducing network round trips
- Atomicity and transaction limitations
- When to use Lua scripts

**Practical Tasks**
- Execute multiple Redis commands efficiently.
- Compare individual commands with pipelining.
- Implement a multi-command operation.
- Explore error handling during transactions.

**Deliverable:** A Go program demonstrating pipelining and transactional operations.

### Day 09 — Pub/Sub Messaging

**Topics**
- Publishers and subscribers
- Redis Pub/Sub
- Channels and message delivery
- Reconnection and subscription lifecycle
- Pub/Sub limitations

**Practical Tasks**
- Create a publisher in Go.
- Create a subscriber in a separate process.
- Send real-time notifications.
- Handle cancellation and shutdown with `context.Context`.

**Deliverable:** A simple real-time messaging system.

### Day 10 — Redis Streams and Consumer Groups

**Topics**
- Redis Streams
- `XADD` and `XREAD`
- Consumer groups
- Message acknowledgment
- Pending messages and recovery

**Practical Tasks**
- Build a background job processing demo.
- Publish jobs to a Redis Stream.
- Process jobs using a consumer group.
- Acknowledge successful messages.

**Deliverable:** A basic event-processing worker.

### Day 11 — Rate Limiting

**Topics**
- Rate-limiting algorithms
- Fixed-window rate limiting
- Sliding-window concepts
- Atomic counters
- Redis Lua scripts

**Practical Tasks**
- Implement a per-user or per-IP rate limiter.
- Add rate limiting middleware to Fiber.
- Return HTTP `429 Too Many Requests` when limits are exceeded.
- Test concurrent requests.

**Deliverable:** Reusable Redis-backed rate-limiting middleware.

### Day 12 — Distributed Locks and Concurrency

**Topics**
- Distributed coordination
- `SET NX` with expiration
- Lock ownership and safe release
- Race conditions
- Lock expiration and failure scenarios

**Practical Tasks**
- Implement a basic distributed lock.
- Use a unique token for each lock owner.
- Release locks safely, preferably with an atomic Lua script.
- Test concurrent workers.

**Deliverable:** A distributed-lock demonstration with documented limitations.

> A basic Redis lock is a learning exercise. Production-grade distributed coordination requires careful handling of expiration, failures, and ownership.

### Day 13 — Persistence, Reliability, and Performance

**Topics**
- RDB snapshots and AOF
- Redis memory management
- Eviction policies
- Connection pooling
- Timeouts and retries
- Monitoring and slow commands

**Practical Tasks**
- Configure Redis persistence.
- Inspect memory usage and server statistics.
- Test behavior after a Redis restart.
- Add timeouts and error handling to the Go client.
- Benchmark common operations.

**Deliverable:** A Redis reliability and performance report.

### Day 14 — Final Project

**Project:** Redis-Powered Go Backend

Build a small backend that combines the concepts from the previous days.

**Required Features**
- Go REST API using Fiber.
- Redis client and configuration package.
- CRUD operations.
- Cache-aside with TTL.
- Rate limiting.
- Background processing with Redis Streams.
- Structured logging and error handling.
- Docker Compose setup.
- Unit tests and integration tests.
- README documentation.

**Suggested Architecture**

```text
Client
  |
  v
Go + Fiber API
  |
  +------> Redis Cache
  |
  +------> Redis Data Structures
  |
  +------> Redis Streams
                 |
                 v
             Go Worker
```

**Deliverable:** A documented, runnable GitHub repository.

## 📂 Suggested Project Structure

```text
redis-go-roadmap/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── config/
│   ├── redisclient/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── cache/
│   └── worker/
├── tests/
├── docker-compose.yml
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

Keep the structure simple during the first few days. Introduce additional packages as the project grows.

## 📦 Installation

Install the Redis Go client:

```bash
go get github.com/redis/go-redis/v9
```

Example `docker-compose.yml` for local development:

```yaml
services:
  redis:
    image: redis:7-alpine
    container_name: redis-roadmap
    ports:
      - "127.0.0.1:6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    restart: unless-stopped

volumes:
  redis_data:
```

Start Redis:

```bash
docker compose up -d
```

Check its status:

```bash
docker compose ps
docker exec -it redis-roadmap redis-cli PING
```

Expected response:

```text
PONG
```

This Compose configuration is intended for local development, not an internet-facing production deployment.

## 📚 Learning Resources

- [Redis Documentation](https://redis.io/docs/latest/)
- [Redis University](https://university.redis.io/)
- [go-redis GitHub Repository](https://github.com/redis/go-redis)
- [Go Documentation](https://go.dev/doc/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [Docker Documentation](https://docs.docker.com/)

## ✅ Progress Tracker

- ✅ Day 01 — Redis Fundamentals
- ✅ Day 02 — Redis Setup with Golang
- ✅ Day 03 — Strings and Basic CRUD
- ✅ Day 04 — Hashes and Structured Data
- [ ] Day 05 — Lists, Sets, and Sorted Sets
- [ ] Day 06 — Expiration and Cache Design
- [ ] Day 07 — Redis and REST APIs with Fiber
- [ ] Day 08 — Transactions and Pipelining
- [ ] Day 09 — Pub/Sub Messaging
- [ ] Day 10 — Redis Streams and Consumer Groups
- [ ] Day 11 — Rate Limiting
- [ ] Day 12 — Distributed Locks and Concurrency
- [ ] Day 13 — Persistence, Reliability, and Performance
- [ ] Day 14 — Final Project

## 🏁 Final Goal

The objective is not just to memorize Redis commands. It is to learn how to use Redis effectively in real Go backend systems, understand its trade-offs, and build software that is maintainable, observable, and reliable.

**Learn by building. Test every feature. Document every day.**
