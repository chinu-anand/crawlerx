# CrawlerX

<div align="center">
  <h3>A high-performance web crawler with real-time monitoring</h3>
</div>

<div align="center">
  
[![Go Version](https://img.shields.io/github/go-mod/go-version/chinu-anand/crawlerx)](https://github.com/chinu-anand/crawlerx)
  
</div>

## 🚀 Overview

CrawlerX is a robust, scalable web crawling system built in Go. It allows you to extract information from websites efficiently with real-time progress monitoring through WebSockets. The system uses a distributed architecture with Redis for job queuing and PostgreSQL for data persistence.

```
🔍 Crawl any website
📊 Monitor progress in real-time
🔄 Automatic retry mechanism
📱 RESTful API for easy integration
🐳 Docker ready for quick deployment
```

## 🏗️ Architecture

CrawlerX follows a clean, modular architecture that separates concerns and enables scalability:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│  API Layer  │────▶│  Job Queue  │────▶│   Workers   │
│             │     │   (Redis)   │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
       │                                       │
       │                                       │
       ▼                                       ▼
┌─────────────┐                         ┌─────────────┐
│             │                         │             │
│  WebSocket  │◀────────────────────────│   Crawler   │
│    Hub      │                         │             │
└─────────────┘                         └─────────────┘
       │                                       │
       │                                       │
       ▼                                       ▼
┌─────────────┐                         ┌─────────────┐
│             │                         │             │
│   Clients   │                         │  Database   │
│             │                         │(PostgreSQL) │
└─────────────┘                         └─────────────┘
```

## ✨ Features

- **High-performance crawling**: Efficiently extracts data from web pages
- **Real-time updates**: WebSocket integration for live progress monitoring
- **Retry mechanism**: Automatic retry for failed requests with configurable delay
- **Distributed architecture**: Redis-based job queue for horizontal scaling
- **Persistent storage**: PostgreSQL database for storing crawl results
- **RESTful API**: Simple API for job submission and monitoring
- **Docker support**: Easy deployment with Docker and docker-compose
- **Export functionality**: Export crawl results in various formats

## 🛠️ Tech Stack

- **Backend**: Go (Golang)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Queue**: Redis
- **WebSockets**: Gorilla WebSocket
- **HTML Parsing**: GoQuery
- **Configuration**: Viper
- **Containerization**: Docker

## 📋 API Endpoints

| Method | Endpoint          | Description                                |
|--------|-------------------|--------------------------------------------|
| GET    | /api/v1/health    | Health check endpoint                      |
| GET    | /api/v1/ws        | WebSocket connection for real-time updates |
| POST   | /api/v1/crawl     | Submit a new crawl job                     |
| GET    | /api/v1/status/:id| Get status of a specific job               |
| GET    | /api/v1/jobs      | List all jobs                              |
| GET    | /api/v1/jobs/:id  | Get detailed information about a job       |
| GET    | /api/v1/jobs/export/:id | Export job results                   |
| GET    | /api/v1/jobs/summary    | Get summary of all jobs              |

## 🚀 Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker and docker-compose (for containerized deployment)
- PostgreSQL (if running locally)
- Redis (if running locally)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/chinu-anand/crawlerx.git
   cd crawlerx
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (create a `.env` file):
   ```
   PORT=8080
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/crawlerx?sslmode=disable
   REDIS_URL=redis://localhost:6379/0
   ```

4. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

### Docker Deployment

1. Build and start the containers:
   ```bash
   docker-compose up -d
   ```

2. The application will be available at `http://localhost:8080`

## 📝 Usage Examples

### Submit a Crawl Job

```bash
curl -X POST http://localhost:8080/api/v1/crawl \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

Response:
```json
{
  "message": "Crawl request received",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Get Job Status

```bash
curl http://localhost:8080/api/v1/status/550e8400-e29b-41d4-a716-446655440000
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://example.com",
  "status": "done",
  "title": "Example Domain",
  "description": "Example website for demonstration purposes",
  "links": ["https://www.iana.org/domains/example"]
}
```

## 🔄 WebSocket Integration

Connect to the WebSocket endpoint to receive real-time updates:

```javascript
const socket = new WebSocket('ws://localhost:8080/api/v1/ws');

socket.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Job update:', data);
};
```

## 🧩 Project Structure

```
crawlerx/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── configs/
│   └── config.go           # Configuration management
├── internal/
│   ├── api/                # API handlers and routes
│   ├── crawler/            # Web crawling implementation
│   ├── models/             # Database models
│   ├── queue/              # Job queue implementation
│   ├── worker/             # Worker processing logic
│   └── ws/                 # WebSocket implementation
├── pkg/                    # Shared packages
├── Dockerfile              # Docker configuration
├── docker-compose.yaml     # Docker Compose configuration
├── go.mod                  # Go module definition
└── README.md               # Project documentation
```

## 📊 Performance

CrawlerX is designed for high performance and scalability:

- **Concurrent processing**: Multiple workers can process jobs simultaneously
- **Efficient HTML parsing**: Using GoQuery for optimized DOM traversal
- **Minimal memory footprint**: Careful resource management
- **Horizontal scaling**: Add more worker instances as needed