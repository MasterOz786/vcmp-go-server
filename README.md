# vcmp-go-server

Project Safari **gamemode library** and VC:MP **server deploy folder** (configs, maps, `server64.exe`).

This repository does **not** compile a native plugin. There is no `main` package and no `go build -buildmode=c-shared` here.

## What lives here

| Path | Role |
|------|------|
| `safari/` | Gamemode logic (teams, Hydra, rounds, commands) — Go **library** |
| `safari/vcmpapi.go` | Adapter: `safari.API` → `vcmp.API` |
| `safari.json`, `safari_maps/` | Runtime config (keep next to `server64.exe`) |
| `server.cfg`, `goserver.json` | VC:MP server config (local deploy) |
| `plugins/` | **Deploy only** — copy built `goserver04rel64.dll` here (not built in this repo) |

## What does NOT live here

- CGO / `plugin.h` / `//export` callbacks → [**vcmp-go-plugin**](https://github.com/masteroz/vcmp-go-plugin)
- Plugin `main`, event wiring → `vcmp-go-plugin/plugin/`

## Build & test (this repo)

```powershell
cd D:\vcmp-go-server
make test
make tidy
```

## Build the plugin (after changing `safari/` or plugin wiring)

```powershell
cd D:\vcmp-go-server
.\build.ps1              # test -> stop server -> build -> deploy
.\build.ps1 -StartServer # same + launch server64.exe
```

Steps run automatically: `go test` in vcmp-go-server, stop `server64`, build `goserver04rel64.dll`, copy to `plugins/`.

## Run the server

```powershell
cd D:\vcmp-go-server
.\server64.exe
```

`server.cfg` must load the Safari plugin:

```
plugins xmlconf04rel64 announce04rel64 goserver04rel64
```

On load you should see:

```
[plugin] loaded Safari v1.0.0 (API 2.0)
[safari] gamemode initialized (map=... db=...)
[safari] server ready — Project Safari: Hydra Warfare
```

## Local development

```go
// vcmp-go-server/go.mod
replace github.com/masteroz/vcmp-go-plugin => ../vcmp-go-plugin
```

```go
// vcmp-go-plugin/plugin/go.mod
replace github.com/masteroz/vcmp-go-server => ../../vcmp-go-server
```
