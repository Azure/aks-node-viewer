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

// Package xdg resolves the XDG base directories as described by the
// XDG Base Directory Specification:
// https://specifications.freedesktop.org/basedir-spec/latest/
package xdg

import (
	"os"
	"path/filepath"
)

// home returns the user's home directory, preferring os.UserHomeDir and falling
// back to $HOME.
func home() string {
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return h
	}
	return os.Getenv("HOME")
}

// baseDir returns the value of the given XDG environment variable if it is set
// to an absolute path, otherwise home/<fallback>. Relative values are ignored
// as required by the specification.
func baseDir(env, fallback string) string {
	if dir := os.Getenv(env); filepath.IsAbs(dir) {
		return dir
	}
	return filepath.Join(home(), fallback)
}

// ConfigHome returns the base directory for user-specific configuration files,
// e.g. $XDG_CONFIG_HOME or ~/.config.
func ConfigHome() string {
	return baseDir("XDG_CONFIG_HOME", ".config")
}

// CacheHome returns the base directory for user-specific non-essential cached
// files, e.g. $XDG_CACHE_HOME or ~/.cache.
func CacheHome() string {
	return baseDir("XDG_CACHE_HOME", ".cache")
}

// ConfigDir returns the configuration directory for the named application.
func ConfigDir(app string) string {
	return filepath.Join(ConfigHome(), app)
}

// CacheDir returns the cache directory for the named application.
func CacheDir(app string) string {
	return filepath.Join(CacheHome(), app)
}
