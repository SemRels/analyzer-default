# analyzer-default

Default analyzer plugin for Semantic Release.

Provides the default release analysis strategy for Semantic Release.

## Documentation

- Docs (coming soon): <https://github.com/SemRels/semrel/tree/main/docs/plugins/analyzer-default>
- Template source: <https://github.com/SemRels/plugin-template>

## Repository Layout

`	ext
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
`

## Development

`ash
go build ./cmd/plugin
go test ./...
`

## Configuration Example

`yaml
plugins:
  - name: analyzer-default
    type: analyzer
    config:
      major_keywords: [breaking change, breaking]
      minor_keywords: [feat]
      patch_keywords: [fix, perf]
`

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.