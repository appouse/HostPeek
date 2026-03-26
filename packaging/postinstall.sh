#!/bin/sh
set -e

# Create system user if not exists
if ! id -u hostpeek >/dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin hostpeek
fi

# Enable and start the service
systemctl daemon-reload
systemctl enable hostpeek
systemctl start hostpeek

echo "HostPeek installed and started. Listening on :8080 by default."
echo "Config: /etc/hostpeek/hostpeek.yaml"
