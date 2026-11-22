# WhatsApp Campus — SUP'PTIC (Groupe 12)

Application mobile de communication et de collaboration universitaire conforme au cahier des charges (Yaya, Ndue, Mengue).

## Architecture

- **`database/`** : 6 migrations PostgreSQL (22 tables) + `seed.sql` (rôles, départements, filières, niveaux, serveurs ICE STUN).
- **`server-go/`** : Serveur centralisé 100 % Golang :
  - API REST (`net/http` + `gorilla/mux`) sur `:8080`
  - Signalisation WebRTC (`gorilla/websocket`) sur `/api/ws`
  - Sockets TCP (`net.Listen`, framing `[4 octets big-endian][JSON]`) sur `:5000`
  - Sockets UDP (`net.ListenPacket`, `PING -> PONG`, présence) sur `:5001`
  - Accès exclusif à PostgreSQL via `pgx/v5`
- **`client-android/`** : Projet Android Studio (Java) avec les 16+ écrans, `ApiClient.java` (Retrofit/OkHttp), `WebSocketManager.java`, et `WebRTCManager.java` (unique classe utilisant `org.webrtc`).

## Démarrage rapide

### 1. Base de données PostgreSQL
```bash
createdb whatsapp_campus
for f in database/migrations/*.sql; do psql -d whatsapp_campus -f "$f"; done
psql -d whatsapp_campus -f database/seed.sql
```

### 2. Serveur Go
```bash
cd server-go
cp .env.example .env
go mod tidy
go run ./cmd/server
```

### 3. Client Android Studio
Ouvrir le dossier `client-android/` dans Android Studio, lancer Gradle Sync et exécuter sur l'émulateur (`10.0.2.2:8080` pointe vers votre machine hôte).

### 4. Vérification de conformité (C-T-07 à C-T-12)
```bash
bash scripts/verify_compliance.sh
```
