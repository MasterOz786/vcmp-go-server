// Package safari is the Project Safari gamemode library.
//
// It is imported by vcmp-go-plugin/plugin and linked into the native
// plugin binary. This package is not compiled with -buildmode=c-shared.
//
// Architecture: vcmp.Events handlers in plugin/wiring.go call Engine methods
// directly on the VC:MP callback thread. There is no gameplay event queue.
// Only SQLite I/O is deferred to DBWorker.
package safari
