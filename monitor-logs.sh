#!/bin/bash

# Monitor Throome Docker container logs in real-time
echo "📊 Monitoring Throome Logs (Press Ctrl+C to stop)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Follow logs with timestamps
docker logs -f throome-test --timestamps

