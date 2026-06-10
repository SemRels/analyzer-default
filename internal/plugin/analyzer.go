// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-default Authors

package plugin

import (
	"fmt"
	"regexp"
	"strings"
)

type BumpLevel string

const (
	BumpNone  BumpLevel = "none"
	BumpPatch BumpLevel = "patch"
	BumpMinor BumpLevel = "minor"
	BumpMajor BumpLevel = "major"

	PluginSchemaVersion = 1
)

const (
	majorPatternEnv  = "SEMREL_PLUGIN_MAJOR_PATTERN"
	minorPatternEnv  = "SEMREL_PLUGIN_MINOR_PATTERN"
	patchPatternEnv  = "SEMREL_PLUGIN_PATCH_PATTERN"
	majorPatternsEnv = "SEMREL_PLUGIN_MAJOR_PATTERNS"
	minorPatternsEnv = "SEMREL_PLUGIN_MINOR_PATTERNS"
	patchPatternsEnv = "SEMREL_PLUGIN_PATCH_PATTERNS"
)

var (
	defaultMajorPatterns = []string{`(?i)BREAKING CHANGE`}
	defaultMinorPatterns = []string{`(?i)^feat(?:\(|:)`}
	defaultPatchPatterns = []string{`(?i)^fix(?:\(|:)`, `(?i)^perf(?:\(|:)`}
)

type AnalysisResult struct {
	Bump                BumpLevel `json:"bump"`
	Reason              string    `json:"reason"`
	PluginSchemaVersion int       `json:"plugin_schema_version,omitempty"`
}

type Analyzer struct {
	majorPatterns []*regexp.Regexp
	minorPatterns []*regexp.Regexp
	patchPatterns []*regexp.Regexp
}

func New() *Analyzer {
	analyzer, err := newAnalyzer(defaultMajorPatterns, defaultMinorPatterns, defaultPatchPatterns)
	if err != nil {
		panic(err)
	}

	return analyzer
}

func NewFromEnv(lookupEnv func(string) (string, bool)) (*Analyzer, error) {
	majorPatterns, err := compileConfiguredPatterns(lookupEnv, majorPatternsEnv, majorPatternEnv, defaultMajorPatterns)
	if err != nil {
		return nil, err
	}

	minorPatterns, err := compileConfiguredPatterns(lookupEnv, minorPatternsEnv, minorPatternEnv, defaultMinorPatterns)
	if err != nil {
		return nil, err
	}

	patchPatterns, err := compileConfiguredPatterns(lookupEnv, patchPatternsEnv, patchPatternEnv, defaultPatchPatterns)
	if err != nil {
		return nil, err
	}

	return &Analyzer{
		majorPatterns: majorPatterns,
		minorPatterns: minorPatterns,
		patchPatterns: patchPatterns,
	}, nil
}

func newAnalyzer(major, minor, patch []string) (*Analyzer, error) {
	majorPatterns, err := compilePatterns(major)
	if err != nil {
		return nil, err
	}

	minorPatterns, err := compilePatterns(minor)
	if err != nil {
		return nil, err
	}

	patchPatterns, err := compilePatterns(patch)
	if err != nil {
		return nil, err
	}

	return &Analyzer{
		majorPatterns: majorPatterns,
		minorPatterns: minorPatterns,
		patchPatterns: patchPatterns,
	}, nil
}

func compileConfiguredPatterns(lookupEnv func(string) (string, bool), pluralEnv string, singularEnv string, defaults []string) ([]*regexp.Regexp, error) {
	if lookupEnv == nil {
		return compilePatterns(defaults)
	}

	envName, rawPatterns, ok := configuredPatternValue(lookupEnv, pluralEnv, singularEnv)
	if !ok {
		return compilePatterns(defaults)
	}

	patterns, err := parsePatterns(rawPatterns)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", envName, err)
	}

	compiledPatterns, err := compilePatterns(patterns)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", envName, err)
	}

	return compiledPatterns, nil
}

func configuredPatternValue(lookupEnv func(string) (string, bool), pluralEnv string, singularEnv string) (string, string, bool) {
	if value, ok := lookupEnv(pluralEnv); ok {
		return pluralEnv, value, true
	}

	if value, ok := lookupEnv(singularEnv); ok {
		return singularEnv, value, true
	}

	return "", "", false
}

func parsePatterns(raw string) ([]string, error) {
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	patterns := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		patterns = append(patterns, trimmed)
	}

	return patterns, nil
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regexp %q: %w", pattern, err)
		}

		compiled = append(compiled, re)
	}

	return compiled, nil
}

func (a *Analyzer) Analyze(commits []string) AnalysisResult {
	result := AnalysisResult{
		Bump:   BumpNone,
		Reason: "no commits require a version bump",
	}

	majorCount := 0
	minorCount := 0
	patchCount := 0

	for _, commit := range commits {
		normalized := strings.TrimSpace(commit)

		switch {
		case matchesAny(a.majorPatterns, normalized):
			majorCount++
		case matchesAny(a.minorPatterns, normalized):
			minorCount++
		case matchesAny(a.patchPatterns, normalized):
			patchCount++
		}
	}

	if majorCount > 0 {
		return AnalysisResult{
			Bump:   BumpMajor,
			Reason: fmt.Sprintf("%d breaking change commit(s) detected", majorCount),
		}
	}

	if minorCount > 0 {
		return AnalysisResult{
			Bump:   BumpMinor,
			Reason: fmt.Sprintf("%d feature commit(s) detected", minorCount),
		}
	}

	if patchCount > 0 {
		return AnalysisResult{
			Bump:   BumpPatch,
			Reason: fmt.Sprintf("%d fix/perf commit(s) detected", patchCount),
		}
	}

	return result
}

func matchesAny(patterns []*regexp.Regexp, message string) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(message) {
			return true
		}
	}

	return false
}
