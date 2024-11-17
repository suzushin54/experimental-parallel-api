#!/bin/bash

grpcurl -plaintext localhost:8080 list

grpcurl -plaintext -d '{
  "user_data": {
    "email": "user@example.com",
    "password": "securepassword"
  },
  "payment_data": {
    "amount": 100.0,
    "currency": "USD",
    "method": "credit"
  }
}' localhost:8080 proto.payment.v1.PaymentService.ProcessPayment