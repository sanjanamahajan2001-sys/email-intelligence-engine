#!/bin/bash
curl -s -X POST http://localhost:8080/v1/validate \
     -H "Content-Type: application/json" \
     -d '{"email": "test@new-burner-sync-2.com"}'
