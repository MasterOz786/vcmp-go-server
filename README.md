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
make              # goserver04rel64.so
make build-linux  # Linux rel64 (typical production target)
```

## Deploy

Copy next to the official VC:MP server:

- `goserver04rel64.so` (or `.dll` on Windows)
- `server.cfg` — see `server.cfg.example`
- `goserver.json`
- `safari.json`
- `safari_maps/patrol_default.json`

`server.cfg` plugin line:

```
plugins goserver04rel64
```

No Squirrel or Java required.
