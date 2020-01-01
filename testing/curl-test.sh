#!/bin/bash
curl -ivL -H "Content-Type: application/json" --insecure https://localhost:8443/validate -X POST -d "@example-request.json"
