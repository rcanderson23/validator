#!/bin/bash
curl -ivL -H "Content-Type: application/json" http://localhost:8443/validate -X POST -d "@example-request.json"
