# AnonBoard

AnonBoard is a simple anonymous message board API built with Go and PostgreSQL.

## Features
- Post and retrieve anonymous messages
- REST API with JSON responses
- Uses PostgreSQL for data storage

## Prerequisites
Before running the project, ensure you have the following installed:
- [Go](https://go.dev/dl/) (latest version)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [PostgreSQL](https://www.postgresql.org/download/) (if running without Docker)

## Installation & Setup

### 1. Clone the repository
```sh
 git clone https://github.com/yourusername/go-anonboard.git
 cd go-anonboard
```

### 2. Install dependencies
```sh
go mod tidy
```

### 3. Start PostgreSQL using Docker Compose
```sh
docker compose up -d
```
This will start a PostgreSQL instance with the following credentials:
- **User:** `localuser`
- **Password:** `localpassword`
- **Database:** `localdb`
- **Port:** `5432`

### 4. Configure the application
Create a configuration file `config.yaml` in the config directory:
```yaml
env: "local" # local, prod
postgres_dsn: "postgres://localuser:localpassword@0.0.0.0:5432/localdb"
http_server:
  address: "localhost:8080"
  timeout: 4s
  idle_timeout: 60s
```

### 5. Run the application
```sh
go run main.go
```
The server will start on `http://localhost:8080`.

## API Endpoints

### 1. Create a Post
**Request:**
```sh
curl -X POST http://localhost:8080/api/post \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello, world!"}'
```

**Response:**
```json
{
  "status": "OK",
  "post": {
    "id": 1,
    "nanoid": "snqdclm76f",
    "message": "Hello, world!",
    "created_at": "2025-02-22T13:23:41.13693+05:00"
  }
}
```

### 2. Retrieve All Posts
**Request:**
```sh
curl http://localhost:8080/api/post
```

**Response:**
```json
{
  "status": "OK",
  "post": [
    {
      "id": 1,
      "nanoid": "snqdclm76f",
      "message": "Hello, world!",
      "created_at": "2025-02-22T13:23:41.13693+05:00"
    }
    {
      "id": 2,
      "nanoid": "lsrf0algsz",
      "message": "Hello, world!",
      "created_at": "2025-02-22T13:26:45.642632+05:00"
    }
  ]
}
```

## Stopping the Application
To stop the application, press `CTRL+C` in the terminal running `go run main.go`.

To stop and remove the PostgreSQL container:
```sh
docker compose down
```

## Future Improvements
- Client
- WebSocket support
- Better error handling

