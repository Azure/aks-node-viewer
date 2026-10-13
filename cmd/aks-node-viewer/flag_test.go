/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// withConfigPaths points the config paths at temporary files for the duration of
// the test and restores them afterwards.
func withConfigPaths(t *testing.T) (xdgPath, legacyPath string) {
	t.Helper()
	dir := t.TempDir()
	xdgPath = filepath.Join(dir, "xdg", "aks-node-viewer", "config")
	legacyPath = filepath.Join(dir, ".aks-node-viewer")

	oldConfig, oldLegacy := configPath, legacyConfigPath
	configPath, legacyConfigPath = xdgPath, legacyPath
	t.Cleanup(func() {
		configPath, legacyConfigPath = oldConfig, oldLegacy
	})
	return xdgPath, legacyPath
}

func TestMigrateConfigCopiesLegacy(t *testing.T) {
	xdgPath, legacyPath := withConfigPaths(t)

	// no legacy file: nothing to do
	if err := migrateConfig(); err != nil {
		t.Fatalf("migrateConfig: %v", err)
	}
	if _, err := os.Stat(xdgPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no XDG config to be created without a legacy file")
	}

	// a legacy file is copied to the XDG location
	if err := os.WriteFile(legacyPath, []byte("resources=cpu,memory\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := migrateConfig(); err != nil {
		t.Fatalf("migrateConfig: %v", err)
	}
	if got := readFile(t, xdgPath); got != "resources=cpu,memory\n" {
		t.Errorf("expected the legacy config to be copied, got %q", got)
	}

	// an existing XDG config is never overwritten
	if err := os.WriteFile(legacyPath, []byte("resources=cpu\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := migrateConfig(); err != nil {
		t.Fatalf("migrateConfig: %v", err)
	}
	if got := readFile(t, xdgPath); got != "resources=cpu,memory\n" {
		t.Errorf("expected the existing XDG config to be preserved, got %q", got)
	}
}

func TestLoadConfigFileReadsLegacyWhenXDGMissing(t *testing.T) {
	xdgPath, legacyPath := withConfigPaths(t)

	if err := os.WriteFile(legacyPath, []byte("region=eastus\nresources=cpu,memory\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := loadConfigFile()
	if err != nil {
		t.Fatalf("loadConfigFile: %v", err)
	}
	if got := cfg.getValue("region", ""); got != "eastus" {
		t.Errorf("expected region=eastus, got %q", got)
	}
	if got := cfg.getValue("resources", ""); got != "cpu,memory" {
		t.Errorf("expected resources=cpu,memory, got %q", got)
	}

	// loading migrated the legacy file into the XDG location
	if got := readFile(t, xdgPath); got == "" {
		t.Errorf("expected the legacy config to have been migrated")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	return string(data)
}
