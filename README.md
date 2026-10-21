# analyzer-default

[![Latest Release](https://img.shields.io/github/v/release/SemRels/analyzer-default?label=version&color=blue)](https://github.com/SemRels/analyzer-default/releases/latest)
[![CI](https://github.com/SemRels/analyzer-default/actions/workflows/ci.yml/badge.svg)](https://github.com/SemRels/analyzer-default/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/SemRels/analyzer-default)](LICENSE)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/SemRels/analyzer-default/badge)](https://scorecard.dev/viewer/?uri=github.com/SemRels/analyzer-default)
[![REUSE status](https://api.reuse.software/badge/github.com/SemRels/analyzer-default)](https://api.reuse.software/info/github.com/SemRels/analyzer-default)

Determines the semantic version bump by matching commit messages with configurable regular expressions.

This plugin is distributed as the standalone Go binary `semrel-plugin-analyzer-default`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

### Binary

```bash
go install github.com/SemRels/analyzer-default/cmd/plugin@latest
```

### Docker

Pre-built, multi-platform images (linux/amd64, linux/arm64) are published to the GitHub Container Registry on every release:

```bash
docker pull ghcr.io/semrels/analyzer-default:latest
```

Images are signed with [cosign](https://github.com/sigstore/cosign) and include a full SBOM attestation. Verify the signature:

```bash
cosign verify ghcr.io/semrels/analyzer-default:latest \
  --certificate-identity-regexp 'https://github.com/SemRels/analyzer-default/.github/workflows/container.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - uses: "@semrel/analyzer-default"
    args:
      major_patterns: "BREAKING CHANGE,feat!"
      minor_patterns: "^feat(\\(.+\\))?:,^feature(\\(.+\\))?:,^\\[feature\\]"
      patch_patterns: "^(fix|perf)(\\(.+\\))?:,^bugfix(\\(.+\\))?:,^hotfix(\\(.+\\))?:"
```

Use a comma-separated list to configure multiple regexes for the same bump level. Whitespace around each pattern is trimmed and the bump is triggered when any configured pattern matches. The singular `SEMREL_PLUGIN_*_PATTERN` variables remain supported for backward compatibility, and also accept comma-separated values. If both singular and plural forms are set, the plural form takes precedence.

## Options

Every option can be set as an entry under `args` in `.semrel.yaml`. semrel passes each entry to the plugin as the environment variable shown in the table, so you can also set that variable directly (for example in your CI environment) instead of using `args`.

| Argument | Environment variable | Required | Description | Default |
| --- | --- | --- | --- | --- |
| `minor_patterns` | `SEMREL_PLUGIN_MINOR_PATTERNS` | Optional | Comma-separated regular expressions that mark a commit as minor. | Built-in `feat` matching |
| `patch_patterns` | `SEMREL_PLUGIN_PATCH_PATTERNS` | Optional | Comma-separated regular expressions that mark a commit as patch. | Built-in `fix` / `perf` matching |
| `major_patterns` | `SEMREL_PLUGIN_MAJOR_PATTERNS` | Optional | Comma-separated regular expressions that mark a commit as major. | Built-in `BREAKING CHANGE` matching |
| `minor_pattern` | `SEMREL_PLUGIN_MINOR_PATTERN` | Optional | Backward-compatible singular form for minor matching. | Inherited from defaults unless set |
| `patch_pattern` | `SEMREL_PLUGIN_PATCH_PATTERN` | Optional | Backward-compatible singular form for patch matching. | Inherited from defaults unless set |
| `major_pattern` | `SEMREL_PLUGIN_MAJOR_PATTERN` | Optional | Backward-compatible singular form for major matching. | Inherited from defaults unless set |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_BUMP` | Calculated bump level such as major, minor, or patch. |

## Example behavior

The plugin evaluates commit messages against the configured patterns and returns the highest required bump level. For example, `SEMREL_PLUGIN_PATCH_PATTERNS="^fix,^bugfix,^hotfix"` matches any commit starting with `fix`, `bugfix`, or `hotfix`.

## License

Apache-2.0
