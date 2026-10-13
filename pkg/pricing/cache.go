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

package pricing

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Azure/aks-node-viewer/pkg/xdg"
)

const appName = "aks-node-viewer"

// cacheData is the on-disk representation of the pricing data, persisted so
// that subsequent runs start from the most recently fetched prices instead of
// the embedded static snapshot.
type cacheData struct {
	Region             string             `json:"region"`
	OnDemandUpdateTime time.Time          `json:"onDemandUpdateTime"`
	OnDemandPrices     map[string]float64 `json:"onDemandPrices"`
	SpotUpdateTime     time.Time          `json:"spotUpdateTime"`
	SpotPrices         map[string]float64 `json:"spotPrices"`
}

// unsafeFilenameChars matches characters that are not safe to use in a file
// name; they are replaced so an arbitrary region string cannot escape the cache
// directory.
var unsafeFilenameChars = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// cachePath returns the location of the pricing cache for the given region
// below the XDG cache directory.
func cachePath(region string) string {
	name := "pricing-" + unsafeFilenameChars.ReplaceAllString(region, "_") + ".json"
	return filepath.Join(xdg.CacheDir(appName), name)
}

// readCache loads the cached pricing for the given region. It returns false if
// there is no usable cache (missing, unreadable, corrupt, for another region or
// without any prices).
func readCache(region string) (*cacheData, bool) {
	data, err := os.ReadFile(cachePath(region))
	if err != nil {
		return nil, false
	}
	var cached cacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, false
	}
	if cached.Region != region {
		return nil, false
	}
	if len(cached.OnDemandPrices) == 0 && len(cached.SpotPrices) == 0 {
		return nil, false
	}
	return &cached, true
}

// writeCache persists the pricing data atomically (write to a temporary file in
// the cache directory and rename it into place).
func writeCache(data *cacheData) error {
	path := cachePath(data.Region)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	contents, err := json.Marshal(data)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".pricing-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(contents); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}
