# Purchase Decision System - Curl Examples

## Health Check
curl -X GET http://localhost:8080/health

## Metrics
curl -X GET http://localhost:8080/metrics

## List pending purchases
curl -X GET "http://localhost:8080/api/v1/purchases?status=pending&limit=20&offset=0"

## Get purchase by ID
curl -X GET http://localhost:8080/api/v1/purchases/{purchase_id}

## Approve purchase
curl -X POST http://localhost:8080/api/v1/purchases/{purchase_id}/approve \
-H "Content-Type: application/json" \
-d '{
"reason": "Valid purchase within standards",
"approved_by": "admin@example.com"
}'

## Reject purchase
curl -X POST http://localhost:8080/api/v1/purchases/{purchase_id}/reject \
-H "Content-Type: application/json" \
-d '{
"reason": "Suspicious amount detected",
"rejected_by": "fraud@example.com"
}'

## Ignore purchase
curl -X POST http://localhost:8080/api/v1/purchases/{purchase_id}/ignore \
-H "Content-Type: application/json" \
-d '{
"reason": "Duplicate purchase",
"ignored_by": "system@example.com"
}'

## Decision history
curl -X GET http://localhost:8080/api/v1/purchases/{purchase_id}/history

## List approved purchases
curl -X GET "http://localhost:8080/api/v1/purchases?status=approved&limit=10"

## List rejected purchases
curl -X GET "http://localhost:8080/api/v1/purchases?status=rejected&limit=10"