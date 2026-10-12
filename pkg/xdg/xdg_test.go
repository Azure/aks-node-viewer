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

package xdg

import (
	"path/filepath"
	"testing"
)

func TestConfigDirDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/home/tester")
	if got, want := ConfigDir("app"), filepath.Join("/home/tester", ".config", "app"); got != want {
		t.Errorf("expected ConfigDir to be %q, got %q", want, got)
	}
}

func TestConfigDirOverride(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/xdg/config")
	if got, want := ConfigDir("app"), filepath.Join("/xdg/config", "app"); got != want {
		t.Errorf("expected ConfigDir to be %q, got %q", want, got)
	}
}

func TestCacheDirDefaults(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "/home/tester")
	if got, want := CacheDir("app"), filepath.Join("/home/tester", ".cache", "app"); got != want {
		t.Errorf("expected CacheDir to be %q, got %q", want, got)
	}
}

func TestCacheDirOverride(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "/xdg/cache")
	if got, want := CacheDir("app"), filepath.Join("/xdg/cache", "app"); got != want {
		t.Errorf("expected CacheDir to be %q, got %q", want, got)
	}
}

func TestRelativeEnvIsIgnored(t *testing.T) {
	t.Setenv("HOME", "/home/tester")
	t.Setenv("XDG_CONFIG_HOME", "relative/config")
	t.Setenv("XDG_CACHE_HOME", "relative/cache")
	if got, want := ConfigDir("app"), filepath.Join("/home/tester", ".config", "app"); got != want {
		t.Errorf("expected relative XDG_CONFIG_HOME to be ignored, got %q", got)
	}
	if got, want := CacheDir("app"), filepath.Join("/home/tester", ".cache", "app"); got != want {
		t.Errorf("expected relative XDG_CACHE_HOME to be ignored, got %q", got)
	}
}
