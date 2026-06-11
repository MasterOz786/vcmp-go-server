.PHONY: test tidy build build-start hotreload dev-watch

# Gamemode library tests (this repo only).
test:
	go test ./...

tidy:
	go mod tidy

# Build Safari plugin + deploy to plugins/ (runs ../vcmp-go-plugin/build.ps1).
build:
	powershell -NoProfile -ExecutionPolicy Bypass -File build.ps1

# Same as build, then start server64.exe.
build-start:
	powershell -NoProfile -ExecutionPolicy Bypass -File build.ps1 -StartServer

# Rebuild plugin and restart server64 (same as /reload server, from shell).
hotreload:
	powershell -NoProfile -ExecutionPolicy Bypass -File tools/hotreload.ps1 -WaitSeconds 0

# Auto-rebuild on Go source changes (requires: brew install fswatch).
dev-watch:
	bash tools/dev-watch.sh
