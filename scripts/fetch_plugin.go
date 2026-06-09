//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const pluginHURL = "https://raw.githubusercontent.com/habi498/NPC-VCMP/master/plugin/plugin.h"

func main() {
	root, err := filepath.Abs(filepath.Join(".."))
	if err != nil {
		fatal(err)
	}
	out := filepath.Join(root, "include", "plugin.h")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fatal(err)
	}
	resp, err := http.Get(pluginHURL)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fatal(fmt.Errorf("download failed: %s", resp.Status))
	}
	f, err := os.Create(out)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		fatal(err)
	}
	fmt.Println("wrote", out)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
