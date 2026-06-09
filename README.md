# vcmp-go-server (private)

Project Safari gamemode for VC:MP 0.4. Builds on the public SDK:

**https://github.com/masteroz/vcmp-go-plugin**

## Layout

| Path | Role |
|---|---|
| `main.go` | Sets `vcmp.MetaProvider` / `vcmp.OnLoad` |
| `plugin.go` | Safari engine, SQLite store, event wiring entry |
| `safari_wiring.go` | Maps `vcmp.Events` → Safari engine |
| `safari/` | Gamemode logic (teams, Hydra, rounds, commands) |
| `safari/vcmpapi.go` | Thin adapter: `safari.API` → `vcmp.API` |

CGO, callbacks, and `plugin.h` live in **vcmp-go-plugin** — not duplicated here.

## Dependencies

```bash
# Local development (sibling checkout)
replace github.com/masteroz/vcmp-go-plugin => ../vcmp-go-plugin
```

After publishing the SDK, remove `replace` and pin a version tag in `go.mod`.

## Build

Requires Go 1.25.0, CGO, and a C toolchain matching your VC:MP server OS/arch.

```bash
make              # goserver04rel64.so in plugins/
make build-linux  # Linux rel64 (typical production target)
make build-windows  # goserver04rel64.dll for server64.exe
```

## Deploy

Use the [Blank Server 64bit (August 2024)](http://files.thijn.ovh/download/9cedd88d75c4d0d76369b772342b4ba9/Blank%20Server%20-%209th%20August,%202024.zip) as the base on Windows (`server64.exe`). VC:MP loads native plugins from the `plugins/` folder next to the server binary.

```bash
make build-windows   # -> plugins/goserver04rel64.dll
make build-linux     # -> plugins/goserver04rel64.so
```

Copy into the blank server's **Blank Server 64bit** folder:

- `plugins/goserver04rel64.dll` (or `.so` on Linux) — built by this repo
- Keep the blank server's bundled `plugins/` entries (`xmlconf04rel64`, `announce04rel64`, …)
- `server.cfg` — see `server.cfg.example`
- `goserver.json`
- `safari.json`
- `safari_maps/patrol_default.json`

`server.cfg` plugin line (order matters — xmlconf first):

```
plugins xmlconf04rel64 announce04rel64 goserver04rel64
```

On load you should see console lines like:

```
[plugin] loaded Safari v1.0.0 (API 2.0)
[safari] gamemode initialized (map=... db=...)
[safari] server ready — Project Safari: Hydra Warfare
```

No Squirrel or Java required.
