// Package safari is the Project Safari gamemode library.
//
// It is imported by vcmp-go-plugin/plugin and linked into the native
// plugin binary. This package is not compiled with -buildmode=c-shared.
//
// Layout:
//
//   - apidef/      shared types and the API interface
//   - gameplay/    round, teams, marking, loadouts, hydra patrol
//   - persist/     SQLite store and background worker
//   - stream/      client-script binary codec
//   - clientscript/ outbound UI payload builders
//   - admin/       allowlist and hot-reload helpers
//
// Root files: engine, commands, config, script, hydra, admin, messages, api.
package safari
