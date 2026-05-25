// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The analyzer-default Authors

package plugin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzerAnalyze(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		commits    []string
		wantBump   BumpLevel
		wantReason string
	}{
		{
			name:       "breaking change wins",
			commits:    []string{"fix: patch", "refactor: update parser\n\nBREAKING CHANGE: parser output changed"},
			wantBump:   BumpMajor,
			wantReason: "breaking change commit(s)",
		},
		{
			name:       "feature bump",
			commits:    []string{"feat(api): add pagination"},
			wantBump:   BumpMinor,
			wantReason: "feature commit(s)",
		},
		{
			name:       "patch bump",
			commits:    []string{"perf(core): reduce allocations"},
			wantBump:   BumpPatch,
			wantReason: "fix/perf commit(s)",
		},
		{
			name:       "none",
			commits:    []string{"docs: update readme", "chore: refresh deps"},
			wantBump:   BumpNone,
			wantReason: "no commits require",
		},
		{
			name:       "minor beats patch",
			commits:    []string{"fix: patch", "feat: add command"},
			wantBump:   BumpMinor,
			wantReason: "feature commit(s)",
		},
	}

	analyzer := New()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := analyzer.Analyze(tt.commits)
			require.Equal(t, tt.wantBump, result.Bump)
			require.Contains(t, result.Reason, tt.wantReason)
		})
	}
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	require.True(t, hasPrefix("feat: add command", "feat"))
	require.True(t, hasPrefix("fix(parser): avoid panic", "fix"))
	require.False(t, hasPrefix("prefix: no-op", "fix"))
}
