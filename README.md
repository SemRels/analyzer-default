# analyzer-default

[![Latest Release](https://img.shields.io/github/v/release/SemRels/analyzer-default?label=version\&color=blue)](https://github.com/SemRels/analyzer-default/releases/latest)

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
  --certificate-identity-regexp 'https://github.com/SemRels/analyzer-default/.github/workflows/release.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - name: analyzer-default
    path: ~/.semrel/plugins/semrel-plugin-analyzer-default
    env:
      SEMREL_PLUGIN_MAJOR_PATTERNS: "BREAKING CHANGE,feat!"
      SEMREL_PLUGIN_MINOR_PATTERNS: "^feat(\\(.+\\))?:,^feature(\\(.+\\))?:,^\\[feature\\]"
      SEMREL_PLUGIN_PATCH_PATTERNS: "^(fix|perf)(\\(.+\\))?:,^bugfix(\\(.+\\))?:,^hotfix(\\(.+\\))?:"
```

Use a comma-separated list to configure multiple regexes for the same bump level. Whitespace around each pattern is trimmed and the bump is triggered when any configured pattern matches. The singular `SEMREL_PLUGIN_*_PATTERN` variables remain supported for backward compatibility, and also accept comma-separated values. If both singular and plural forms are set, the plural form takes precedence.

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_MINOR_PATTERNS` | Optional | Comma-separated regular expressions that mark a commit as minor. | Built-in `feat` matching |
| `SEMREL_PLUGIN_PATCH_PATTERNS` | Optional | Comma-separated regular expressions that mark a commit as patch. | Built-in `fix` / `perf` matching |
| `SEMREL_PLUGIN_MAJOR_PATTERNS` | Optional | Comma-separated regular expressions that mark a commit as major. | Built-in `BREAKING CHANGE` matching |
| `SEMREL_PLUGIN_MINOR_PATTERN` | Optional | Backward-compatible singular form for minor matching. | Inherited from defaults unless set |
| `SEMREL_PLUGIN_PATCH_PATTERN` | Optional | Backward-compatible singular form for patch matching. | Inherited from defaults unless set |
| `SEMREL_PLUGIN_MAJOR_PATTERN` | Optional | Backward-compatible singular form for major matching. | Inherited from defaults unless set |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_BUMP` | Calculated bump level such as major, minor, or patch. |

## Example behavior

The plugin evaluates commit messages against the configured patterns and returns the highest required bump level. For example, `SEMREL_PLUGIN_PATCH_PATTERNS="^fix,^bugfix,^hotfix"` matches any commit starting with `fix`, `bugfix`, or `hotfix`.

## License

Apache-2.0
