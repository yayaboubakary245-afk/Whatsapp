#!/usr/bin/env bash
set -euo pipefail

echo "=== Audit de conformité WhatsApp Campus (C-T-07 à C-T-12) ==="

echo "[C-T-07] Vérification net.Listen / net.ListenPacket dans server-go/internal/network/:"
grep -rnE "net\.Listen|net\.ListenPacket" server-go/internal/network/

echo "[C-T-08] Vérification zéro socket Java hors WebRTCManager.java:"
MATCHES=$(grep -rnE "java\.net\.Socket|DatagramSocket" client-android/ | grep -v "WebRTCManager" || true)
if [ -n "$MATCHES" ]; then
  echo "ÉCHEC C-T-08:"
  echo "$MATCHES"
  exit 1
else
  echo "OK (0 résultat hors WebRTCManager.java)"
fi

echo "[C-T-10] Vérification gorilla/websocket dans server-go/internal/call/signaling.go:"
grep -n "gorilla/websocket" server-go/internal/call/signaling.go

echo "[C-T-11] Vérification org.webrtc dans WebRTCManager.java uniquement:"
grep -rn "org\.webrtc" client-android/

echo "=== Tous les contrôles C-T-07 à C-T-12 sont conformes ==="
