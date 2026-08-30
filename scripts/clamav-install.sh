#!/usr/bin/env bash
# Install ClamAV and ensure virus signatures are present (CI / release).
set -euo pipefail

sudo apt-get update
sudo apt-get install -y clamav

# Package install starts freshclam daemon; stop it so manual freshclam can run.
sudo systemctl stop clamav-freshclam 2>/dev/null || true
sudo killall freshclam 2>/dev/null || true

if [ ! -f /var/lib/clamav/main.cvd ] && [ ! -f /var/lib/clamav/main.cld ]; then
  echo "ClamAV database missing; running freshclam..."
  sudo freshclam --stdout
fi

if [ ! -f /var/lib/clamav/main.cvd ] && [ ! -f /var/lib/clamav/main.cld ]; then
  echo "ClamAV database still missing after freshclam" >&2
  exit 1
fi

echo "ClamAV ready: $(ls -lh /var/lib/clamav/*.cvd /var/lib/clamav/*.cld 2>/dev/null | wc -l) signature file(s)"
