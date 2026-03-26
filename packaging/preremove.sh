#!/bin/sh
set -e

# Stop and disable the service before removal
systemctl stop hostpeek 2>/dev/null || true
systemctl disable hostpeek 2>/dev/null || true
systemctl daemon-reload
