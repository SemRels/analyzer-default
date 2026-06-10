// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-default Authors

package main

import (
	"bytes"
	"encoding/json"
	"testing"

	plugin "github.com/SemRels/analyzer-default/internal/plugin"
	"github.com/stretchr/testify/require"
)

func TestRunWritesAnalysisResult(t *testing.T) {
	t.Parallel()

	lookupEnv := func(key string) (string, bool) {
		if key == "SEMREL_COMMITS" {
			return `["fix: patch issue","feat: add feature"]`, true
		}
		return "", false
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, lookupEnv)

	require.Equal(t, 0, code)
	require.Empty(t, stderr.String())

	var result plugin.AnalysisResult
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &result))
	require.Equal(t, plugin.BumpMinor, result.Bump)
	require.Equal(t, plugin.PluginSchemaVersion, result.PluginSchemaVersion)
}

func TestRunRejectsInvalidCommitJSON(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, func(string) (string, bool) { return `[`, true })

	require.Equal(t, 1, code)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "invalid SEMREL_COMMITS JSON")
}

func TestRunAllowsMissingCommitEnv(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, func(string) (string, bool) { return "", false })

	require.Equal(t, 0, code)
	require.Empty(t, stderr.String())
	require.Contains(t, stdout.String(), "\"bump\":\"none\"")
	require.Contains(t, stdout.String(), "\"plugin_schema_version\":1")
}

func TestRunRejectsInvalidPatternConfig(t *testing.T) {
	t.Parallel()

	lookupEnv := func(key string) (string, bool) {
		if key == "SEMREL_PLUGIN_MINOR_PATTERN" {
			return `[`, true
		}
		return "", false
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(&stdout, &stderr, lookupEnv)

	require.Equal(t, 1, code)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "SEMREL_PLUGIN_MINOR_PATTERN")
}
