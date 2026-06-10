# vcmp-go-server (private)

Project Safari gamemode library for VC:MP 0.4. Gameplay logic lives here; the native plugin binary is built from **vcmp-go-plugin**.

## Layout

| Path | Role |
|---|---|
| `safari/` | Gamemode logic (teams, Hydra, rounds, commands) |
| `safari/vcmpapi.go` | Thin adapter: `safari.API` → `vcmp.API` |
| `safari.json`, `safari_maps/` | Runtime config (copy to VC:MP server root on deploy) |
| `server.cfg.example`, `goserver.json` | Deploy templates |

Plugin entry and event wiring: [`../plugin/examples/safari/`](../plugin/examples/safari/).

CGO, callbacks, and `plugin.h` live in **vcmp-go-plugin** — not duplicated here.

## Dependencies

```bash
# server/go.mod
replace github.com/masteroz/vcmp-go-plugin => ../plugin
```

## Build (library tests)

```bash
make test    # go test ./...
make tidy
```

## Build Safari plugin (.so / .dll)

From the plugin repo:

```bash
cd ../plugin
make deps
make build-safari          # → plugins/goserver04rel64.so
make build-linux-safari    # Linux rel64
make build-windows-safari  # Windows DLL for server64.exe
```

## Deploy

Use the [Blank Server 64bit (August 2024)](http://files.thijn.ovh/download/9cedd88d75c4d0d76369b772342b4ba9/Blank%20Server%20-%209th%20August,%202024.zip) as the base. Copy into the blank server folder:

- `plugin/plugins/goserver04rel64.dll` (or `.so` on Linux)
- Keep bundled plugins (`xmlconf04rel64`, `announce04rel64`, …)
- `server.cfg` — see `server.cfg.example`
- `goserver.json`, `safari.json`, `safari_maps/patrol_default.json`

```
plugins xmlconf04rel64 announce04rel64 goserver04rel64
```

On load:

```
[plugin] loaded Safari v1.0.0 (API 2.0)
[safari] gamemode initialized (map=... db=...)
[safari] server ready — Project Safari: Hydra Warfare
```
