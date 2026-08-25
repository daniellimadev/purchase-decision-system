# 🛒 Purchase Decision System

Purchase decision system in Golang that consumes SQS events and allows approving, rejecting, or ignoring purchases.

## 🚀 Features

- ✅ **SQS event consumption** (required)
- ✅ Purchase approval
- ✅ Purchase rejection
- ✅ Purchase ignoring
- ✅ SNS event publishing
- ✅ Complete REST API
- ✅ Worker to process SQS queue
- ✅ Business rule validation
- ✅ Decision history in PostgreSQL
- ✅ Structured logs
- ✅ Processing metrics

## 📋 Prerequisites

- Go 1.25+
- AWS Account (SQS and SNS configured)
- PostgreSQL 14+
- Docker (optional)

## 🛠️ Installation

```bash
# Clone the repository
git clone <repo-url>
cd purchase-decision-system

# Install dependencies
go mod download

# Configure environment variables
cp .env.development .env
# Edit the .env file with your credentials

# Run database migrations
go run cmd/api/main.go migrate

# Start the API
go run cmd/api/main.go

# In another terminal, start the worker
go run cmd/worker/main.go
```

## 🏗️ Architecture

```
purchase-decision-system/
├── cmd/
│   ├── api/          # REST API
│   └── worker/       # SQS Worker
├── internal/
│   ├── domain/       # Entities and business rules
│   ├── services/     # Business logic
│   ├── handlers/     # HTTP handlers
│   ├── repository/   # Data access
│   └── infrastructure/
│       ├── aws/      # AWS integration (SQS/SNS)
│       ├── config/   # Configuration
│       └── database/ # Database setup
├── scripts/          # Helper scripts
└── terraform/        # Infrastructure as code
```

## 📡 API Endpoints

### Purchase Decisions

```bash
# Approve purchase
POST /api/v1/purchases/:id/approve
{
  "reason": "Valid purchase",
  "approved_by": "admin@example.com"
}

# Reject purchase
POST /api/v1/purchases/:id/reject
{
  "reason": "Suspicious amount",
  "rejected_by": "admin@example.com"
}

# Ignore purchase
POST /api/v1/purchases/:id/ignore
{
  "reason": "Duplicate",
  "ignored_by": "admin@example.com"
}

# List pending purchases
GET /api/v1/purchases?status=pending

# Find purchase by ID
GET /api/v1/purchases/:id

# Decision history
GET /api/v1/purchases/:id/history
```

### Health Check

```bash
GET /health
GET /metrics
```

## 🔄 Processing Flow

1. **Event arrives in SQS** with purchase data
2. **Worker consumes** the event
3. **Automatic validations** are applied
4. **Business rules** determine the action:
   - Auto-approve if < R$ 500
   - Auto-reject if > R$ 9,500
   - Pending for manual review
5. **Decision is saved** to the database
6. **Event is published** to SNS
7. **Message is deleted** from SQS

## 🧪 Tests

```bash
# Run all tests
go test ./...

# Tests with coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...
```

## 🐳 Docker

```bash
# Build
docker build -t purchase-decision-system .

# Run API
docker run -p 8080:8080 --env-file .env purchase-decision-system api

# Run Worker
docker run --env-file .env purchase-decision-system worker
```

## 📊 SQS Event Example

```json
{
  "purchase_id": "550e8400-e29b-41d4-a716-446655440000",
  "customer_id": "customer-123",
  "amount": 1500.50,
  "currency": "BRL",
  "merchant": "Store XYZ",
  "payment_method": "credit_card",
  "timestamp": "2026-01-25T10:30:00Z",
  "metadata": {
    "ip_address": "192.168.1.1",
    "device": "mobile"
  }
}
```

## 📈 Metrics

- Total purchases processed
- Approval/rejection rate
- Average processing time
- Processing errors

## 🔐 Security

- Input validation
- Rate limiting
- JWT authentication (optional)
- Audit logs
