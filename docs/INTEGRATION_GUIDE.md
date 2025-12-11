# Infrastructure Integration Guide

This guide documents how to use the integrated infrastructure components in this boilerplate. All components are designed to be modular and can be enabled/disable via `config.yaml`.

## 1. Redis

### Configuration (`config.yaml`)
```yaml
redis:
  enabled: true
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```

### Usage (Code)
The `RedisManager` provides a wrapper around `go-redis`.

```go
// Inject RedisManager into your service
type MyService struct {
    redis *infrastructure.RedisManager
}

func (s *MyService) Example() {
    ctx := context.Background()

    // SET
    err := s.redis.Set(ctx, "my-key", "my-value", time.Minute*10)

    // GET
    val, err := s.redis.Get(ctx, "my-key")
    
    // DELETE
    err = s.redis.Delete(ctx, "my-key")
}
```

---

## 2. Postgres

### Configuration (`config.yaml`)
```yaml
postgres:
  enabled: true
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "mydb"
  sslmode: "disable"
  max_open_conns: 10
  max_idle_conns: 5
```

### Usage (Code)
The `PostgresManager` wraps `sqlx.DB`, providing struct mapping and helper methods.

```go
// Inject PostgresManager
type MyService struct {
    db *infrastructure.PostgresManager
}

func (s *MyService) Example() {
    // Access underlying sqlx.DB
    var users []User
    err := s.db.DB.Select(&users, "SELECT * FROM users WHERE active = $1", true)
    
    // Using transaction helper (if implemented in your manager extensions) or usage of standard sqlx patterns
    tx, err := s.db.DB.Beginx()
    // ...
}
```

---

## 3. Kafka

### Configuration (`config.yaml`)
```yaml
kafka:
  enabled: true
  brokers: ["localhost:9092"]
  topic: "my-topic"
  group_id: "my-group"
```

### Usage (Code)
The `KafkaManager` handles producing messages.

```go
// Inject KafkaManager
type MyService struct {
    kafka *infrastructure.KafkaManager
}

func (s *MyService) SendNotification() {
    // Publish a message
    err := s.kafka.Publish("notification-topic", []byte("Hello Kafka"))
    
    // Publish with Key (if supported by your specific implementation extension, default Publish typically sends value)
}
```

---

## 4. MinIO (Object Storage)

### Configuration (`config.yaml`)
```yaml
monitoring:
  minio:
    enabled: true
    endpoint: "localhost:9000"
    access_key: "minioadmin"
    secret_key: "minioadmin"
    use_ssl: false
    bucket: "my-bucket"
    region: "us-east-1"
```

### Usage (Code)
The `MinIOManager` simplifies file uploads and URL retrieval.

```go
// Inject MinIOManager
type MyService struct {
    storage *infrastructure.MinIOManager
}

func (s *MyService) UploadAvatar(fileHeader *multipart.FileHeader) {
    file, _ := fileHeader.Open()
    defer file.Close()

    // Upload
    info, err := s.storage.UploadFile(context.Background(), "avatars/user-1.jpg", file, fileHeader.Size, "image/jpeg")

    // Get Presigned URL (for private buckets) or direct URL
    url := s.storage.GetFileUrl("avatars/user-1.jpg")
}
```

---

## 5. Cron Jobs

### Configuration (`config.yaml`)
Cron jobs can be defined in config for simple logging/testing, or registered in code for logic.

```yaml
cron:
  enabled: true
  jobs:
    "cleanup_logs": "0 0 * * *"   # Run at midnight
    "health_check": "*/5 * * * *" # Run every 5 minutes
```

### Usage (Code)
The `CronManager` allows dynamic job registration.

```go
// Inject CronManager
type MyService struct {
    cron *infrastructure.CronManager
}

func (s *MyService) InitJobs() {
    // Register a new job
    id, err := s.cron.AddJob("database_backup", "0 3 * * *", func() {
        fmt.Println("Performing database backup...")
        // Call service logic here
    })
}
```
