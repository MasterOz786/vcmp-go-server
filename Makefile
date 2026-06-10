.PHONY: test tidy build build-start

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
