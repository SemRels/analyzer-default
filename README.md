# analyzer-default

Determines the semantic version bump by matching commit messages with configurable regular expressions.

This plugin is distributed as the standalone Go binary `semrel-plugin-analyzer-default`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

```bash
go install github.com/SemRels/analyzer-default/cmd/plugin@latest
```

## Configuration

```yaml
plugins:
  - name: analyzer-default
    path: ~/.semrel/plugins/semrel-plugin-analyzer-default
    env:
      SEMREL_PLUGIN_MINOR_PATTERN: "^feat(\\(.+\\))?:"
      SEMREL_PLUGIN_PATCH_PATTERN: "^(fix|perf|refactor)(\\(.+\\))?:"
      SEMREL_PLUGIN_MAJOR_PATTERN: "BREAKING CHANGE"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_MINOR_PATTERN` | Optional | Regular expression that marks a commit as minor. | None |
| `SEMREL_PLUGIN_PATCH_PATTERN` | Optional | Regular expression that marks a commit as patch. | None |
| `SEMREL_PLUGIN_MAJOR_PATTERN` | Optional | Regular expression that marks a commit as major. | None |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_BUMP` | Calculated bump level such as major, minor, or patch. |

## Example behavior

The plugin evaluates commit messages against the configured patterns and returns the highest required bump level.

## License

Apache-2.0
