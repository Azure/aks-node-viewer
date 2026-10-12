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
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/client-go/util/homedir"

	"github.com/Azure/aks-node-viewer/pkg/xdg"
)

const appName = "aks-node-viewer"

var (
	homeDir          string
	configPath       string
	legacyConfigPath string
	version          = "dev"
	commit           = ""
	date             = ""
	builtBy          = ""
)

func init() {
	homeDir = homedir.HomeDir()
	// Prefer the XDG config directory, falling back to the legacy ~/.aks-node-viewer
	// location when no XDG config has been created.
	configPath = filepath.Join(xdg.ConfigDir(appName), "config")
	legacyConfigPath = filepath.Join(homeDir, ".aks-node-viewer")
}

type Flags struct {
	Context         string
	NodeSelector    string
	ExtraLabels     string
	NodeSort        string
	Style           string
	Region          string
	Kubeconfig      string
	Resources       string
	DisablePricing  bool
	ShowAttribution bool
	Version         bool
}

func ParseFlags() (Flags, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var flags Flags

	cfg, err := loadConfigFile()
	if err != nil {
		return Flags{}, fmt.Errorf("load config file: %w", err)
	}

	flagSet.BoolVar(&flags.Version, "v", false, "Display aks-node-viewer version")
	flagSet.BoolVar(&flags.Version, "version", false, "Display aks-node-viewer version")

	contextDefault := cfg.getValue("context", "")
	flagSet.StringVar(&flags.Context, "context", contextDefault, "Name of the kubernetes context to use")

	nodeSelectorDefault := cfg.getValue("node-selector", "")
	flagSet.StringVar(&flags.NodeSelector, "node-selector", nodeSelectorDefault, "Node label selector used to filter nodes, if empty all nodes are selected ")

	extraLabelsDefault := cfg.getValue("extra-labels", "")
	flagSet.StringVar(&flags.ExtraLabels, "extra-labels", extraLabelsDefault, "A comma separated set of extra node labels to display")

	nodeSortDefault := cfg.getValue("node-sort", "creation=dsc")
	flagSet.StringVar(&flags.NodeSort, "node-sort", nodeSortDefault, "Sort order for the nodes, either 'creation' or a label name. The sort order defaults to ascending and can be controlled by appending =asc or =dsc to the value.")

	styleDefault := cfg.getValue("style", "#04B575,#FFFF00,#FF0000")
	flagSet.StringVar(&flags.Style, "style", styleDefault, "Three comma separated colors to use for styling 'good', 'ok' and 'bad' values. The colors are also used in the gradients displayed from bad -> good.")

	regionDefault := cfg.getValue("region", "westus2")
	flagSet.StringVar(&flags.Region, "region", regionDefault, "Azure region to use for pricing lookups")

	// flag overrides env. var. and env. var. overrides config file
	kubeconfigDefault := getStringEnv("KUBECONFIG", cfg.getValue("kubeconfig", filepath.Join(homeDir, ".kube", "config")))
	flagSet.StringVar(&flags.Kubeconfig, "kubeconfig", kubeconfigDefault, "Absolute path to the kubeconfig file")

	resourcesDefault := cfg.getValue("resources", "cpu")
	flagSet.StringVar(&flags.Resources, "resources", resourcesDefault, "List of comma separated resources to monitor")

	disablePricingDefault := cfg.getBoolValue("disable-pricing", false)
	flagSet.BoolVar(&flags.DisablePricing, "disable-pricing", disablePricingDefault, "Disable pricing lookups")

	flagSet.BoolVar(&flags.ShowAttribution, "attribution", false, "Show the Open Source Attribution")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return Flags{}, err
	}
	return flags, nil
}

// --- env vars ---

func getStringEnv(envName string, defaultValue string) string {
	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}
	return env
}

// --- config file ---

type configFile map[string]string

func (c configFile) getValue(key string, defaultValue string) string {
	if val, ok := c[key]; ok {
		return val
	}
	return defaultValue
}

func (c configFile) getBoolValue(key string, defaultValue bool) bool {
	if val, ok := c[key]; ok {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func loadConfigFile() (configFile, error) {
	fileContent := make(map[string]string)

	if err := migrateConfig(); err != nil {
		// migration is best-effort; a failure shouldn't stop the viewer
		log.Printf("warning: could not migrate legacy config: %v", err)
	}

	path := configPath
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// fall back to the legacy ~/.aks-node-viewer location if there is no
		// XDG config (e.g. the user removed it or migration was skipped)
		if _, legacyErr := os.Stat(legacyConfigPath); legacyErr == nil {
			path = legacyConfigPath
		} else {
			return fileContent, nil
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		lineKV := strings.SplitN(line, "=", 2)
		if len(lineKV) == 2 {
			key := strings.TrimSpace(lineKV[0])
			value := strings.TrimSpace(lineKV[1])
			fileContent[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return fileContent, nil
}

// migrateConfig copies the legacy ~/.aks-node-viewer config into the XDG config
// directory. It is a no-op when the XDG config already exists or there is no
// legacy file. The legacy file is left in place and is still read as a fallback
// if the XDG config is later removed.
func migrateConfig() error {
	if configPath == "" || legacyConfigPath == "" || legacyConfigPath == configPath {
		return nil
	}
	if _, err := os.Stat(configPath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data, err := os.ReadFile(legacyConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return err
	}
	log.Printf("migrated config from %s to %s", legacyConfigPath, configPath)
	return nil
}
