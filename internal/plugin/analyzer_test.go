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

func TestNewFromEnvPatterns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		env     map[string]string
		commits []string
		want    BumpLevel
	}{
		{
			name: "single pattern unchanged",
			env: map[string]string{
				minorPatternEnv: `^feat(\(.+\))?:`,
			},
			commits: []string{"feat(api): add pagination"},
			want:    BumpMinor,
		},
		{
			name: "multi pattern first matches",
			env: map[string]string{
				patchPatternEnv: `^fix,^bugfix`,
			},
			commits: []string{"fix: resolve panic"},
			want:    BumpPatch,
		},
		{
			name: "multi pattern second matches",
			env: map[string]string{
				patchPatternEnv: `^fix,^bugfix`,
			},
			commits: []string{"bugfix: resolve panic"},
			want:    BumpPatch,
		},
		{
			name: "multi pattern none match",
			env: map[string]string{
				patchPatternEnv: `^fix,^bugfix`,
			},
			commits: []string{"docs: update readme"},
			want:    BumpNone,
		},
		{
			name: "whitespace trimmed",
			env: map[string]string{
				patchPatternEnv: `  ^fix  ,  ^hotfix  `,
			},
			commits: []string{"hotfix: resolve panic"},
			want:    BumpPatch,
		},
		{
			name: "empty string pattern does not match",
			env: map[string]string{
				patchPatternEnv: ``,
			},
			commits: []string{"fix: resolve panic"},
			want:    BumpNone,
		},
		{
			name: "plural env takes precedence",
			env: map[string]string{
				minorPatternEnv:  `^feat`,
				minorPatternsEnv: `^feature`,
			},
			commits: []string{"feat: add pagination"},
			want:    BumpNone,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			analyzer, err := NewFromEnv(lookupEnv(tt.env))
			require.NoError(t, err)

			result := analyzer.Analyze(tt.commits)
			require.Equal(t, tt.want, result.Bump)
		})
	}
}

func TestNewFromEnvRejectsInvalidRegexp(t *testing.T) {
	t.Parallel()

	_, err := NewFromEnv(lookupEnv(map[string]string{minorPatternEnv: `[`}))
	require.Error(t, err)
	require.Contains(t, err.Error(), minorPatternEnv)
}

func lookupEnv(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
