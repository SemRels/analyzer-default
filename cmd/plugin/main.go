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
	os.Exit(run(os.Stdout, os.Stderr, os.Getenv))
}

func run(stdout, stderr io.Writer, getenv func(string) string) int {
	raw := getenv("SEMREL_COMMITS")

	var commits []string
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &commits); err != nil {
			fmt.Fprintln(stderr, "analyzer-default: invalid SEMREL_COMMITS JSON:", err)
			return 1
		}
	}

	if err := json.NewEncoder(stdout).Encode(plugin.New().Analyze(commits)); err != nil {
		fmt.Fprintln(stderr, "analyzer-default:", err)
		return 1
	}

	return 0
}
