#!/bin/bash

grpcurl -plaintext localhost:8080 list

grpcurl -plaintext -d '{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 100.0,
  "currency": "USD",
  "method": "credit"
}' localhost:8080 payment.PaymentService.ProcessPayment
