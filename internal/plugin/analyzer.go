// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-default Authors

package plugin

import (
	"fmt"
	"strings"
)

type BumpLevel string

const (
	BumpNone  BumpLevel = "none"
	BumpPatch BumpLevel = "patch"
	BumpMinor BumpLevel = "minor"
	BumpMajor BumpLevel = "major"
)

type AnalysisResult struct {
	Bump   BumpLevel `json:"bump"`
	Reason string    `json:"reason"`
}

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
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
		lower := strings.ToLower(normalized)

		if strings.Contains(strings.ToUpper(normalized), "BREAKING CHANGE") {
			majorCount++
			continue
		}

		switch {
		case hasPrefix(lower, "feat"):
			minorCount++
		case hasPrefix(lower, "fix"), hasPrefix(lower, "perf"):
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

func hasPrefix(message string, prefix string) bool {
	return strings.HasPrefix(message, prefix+":") || strings.HasPrefix(message, prefix+"(")
}
