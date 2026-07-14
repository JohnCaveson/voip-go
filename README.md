How to run
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# From the voip-app directory
cd voip-app

# Build the frontend
cd frontend && npm install && npm run build && cd ..

# Quick reference
# Development
wails dev -tags webkit2_41

# Production build
wails build -tags webkit2_41

# Start MySQL
docker compose up -d

# Run the app with MySQL backend
VOIP_STORAGE_TYPE=mysql wails dev -tags webkit2_41

# Run the WAN signaling server separately
go run ./server/cmd/server --port 9321
