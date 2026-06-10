// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-default Authors

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	plugin "github.com/SemRels/analyzer-default/internal/plugin"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.LookupEnv))
}

func run(stdout, stderr io.Writer, lookupEnv func(string) (string, bool)) int {
	raw, _ := lookupEnv("SEMREL_COMMITS")

	var commits []string
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &commits); err != nil {
			fmt.Fprintln(stderr, "analyzer-default: invalid SEMREL_COMMITS JSON:", err)
			return 1
		}
	}

	analyzer, err := plugin.NewFromEnv(lookupEnv)
	if err != nil {
		fmt.Fprintln(stderr, "analyzer-default:", err)
		return 1
	}

	result := analyzer.Analyze(commits)
	result.PluginSchemaVersion = plugin.PluginSchemaVersion

	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintln(stderr, "analyzer-default:", err)
		return 1
	}

	return 0
}
