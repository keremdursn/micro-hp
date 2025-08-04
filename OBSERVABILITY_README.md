# 🔍 Observability Stack

Bu mikroservis projesi comprehensive bir observability stack ile donatılmıştır.

## 📊 Stack Bileşenleri

### 1. **Structured Logging (Zap)**
- JSON formatında structured logging
- Correlation ID ile request tracking
- Service, operation ve context bilgileri

### 2. **Metrics (Prometheus + Grafana)**
- HTTP request metrics
- Database operation metrics  
- Custom business metrics
- Real-time dashboards

### 3. **Distributed Tracing (Simple Zipkin)**
- Lightweight HTTP-based tracing
- Service-to-service call tracking
- Database operation tracing
- Request flow visualization

### 4. **Log Aggregation (Loki + Promtail)**
- Centralized log collection
- Log correlation and search
- Integration with Grafana

## 🚀 Başlatma

```bash
# Tüm stack'i başlat
docker compose up --build -d

# Go modules güncelle
cd hospital-shared && go mod tidy
cd ../auth-service && go mod tidy
cd ../hospital-service && go mod tidy
cd ../personnel-service && go mod tidy
```

## 🔗 URL'ler

| Service | URL | Açıklama |
|---------|-----|----------|
| **Grafana** | http://localhost:3000 | Dashboards (admin/admin) |
| **Prometheus** | http://localhost:9090 | Metrics collection |
| **Zipkin** | http://localhost:9411 | Distributed tracing |
| **Loki** | http://localhost:3100 | Log aggregation |

## 📈 Monitoring Özellikleri

### Request Tracking
```bash
# Her request otomatik correlation ID alır
curl -H "X-Correlation-ID: my-trace-123" http://localhost:8081/api/auth/register
```

### Log Format
```json
{
  "timestamp": "2025-08-03T17:00:00Z",
  "level": "info",
  "service": "auth-service",
  "correlation_id": "my-trace-123",
  "method": "POST",
  "path": "/api/auth/register",
  "status_code": 201,
  "duration": "250ms",
  "type": "http_request"
}
```

### Database Tracing
```json
{
  "timestamp": "2025-08-03T17:00:00Z",
  "level": "info",
  "service": "auth-service",
  "correlation_id": "my-trace-123",
  "operation": "INSERT",
  "table": "authorities",
  "duration": "15ms",
  "type": "database"
}
```

### Service Call Logging
```json
{
  "timestamp": "2025-08-03T17:00:00Z",
  "level": "info", 
  "service": "auth-service",
  "correlation_id": "my-trace-123",
  "target_service": "hospital-service",
  "endpoint": "/api/hospital",
  "status_code": 201,
  "duration": "120ms",
  "type": "service_call"
}
```

## 🎯 Kullanım Örnekleri

### 1. Request Tracing
Grafana'da correlation_id ile tüm servisler arasında request'i takip edebilirsiniz.

### 2. Error Investigation
Loki'de error level logları filtreleyerek sorunları analiz edebilirsiniz.

### 3. Performance Analysis
Zipkin'de slow request'leri identify edebilir, bottleneck'leri bulabilirsiniz.

### 4. Business Metrics
Prometheus'da custom metric'ler ile business KPI'ları track edebilirsiniz.

## 🛠️ Development

### Log Seviyesi Değiştirme
```bash
# Debug mode için
export LOG_LEVEL=debug
```

### Custom Metrics Ekleme
```go
import "hospital-shared/logging"

// Info log
logging.GlobalLogger.LogInfo(ctx, "User action", 
    zap.String("action", "profile_update"),
    zap.Uint("user_id", userID))

// Error log  
logging.GlobalLogger.LogError(ctx, err, "Database error",
    zap.String("operation", "user_create"))
```

### Tracing Spans
```go
import "hospital-shared/tracing"

span, ctx := tracing.StartServiceSpan(ctx, "user-service", "create-user")
defer span.Finish()

// Error durumunda
tracing.FinishSpanWithError(span, err)
```

Bu observability stack ile mikroservislerinizin sağlığını, performansını ve davranışlarını real-time olarak monitor edebilirsiniz! 🚀