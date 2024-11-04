#!/bin/bash

grpcurl -plaintext localhost:8080 list

grpcurl -plaintext -d '{
  "user_id": "user123",
  "amount": 100.0,
  "currency": "USD",
  "method": "credit_card"
}' localhost:8080 payment.PaymentService.ProcessPayment

