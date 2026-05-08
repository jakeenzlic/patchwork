# patchwork

A lightweight diff-based config migration tool that tracks and applies incremental changes to JSON/YAML configs across deployments.

---

## Installation

```bash
go install github.com/yourname/patchwork@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/patchwork.git && cd patchwork && go build ./...
```

---

## Usage

Create a migration patch file:

```yaml
# migrations/001_add_feature_flags.yaml
version: "001"
description: "Add feature flags block"
ops:
  - op: add
    path: /feature_flags
    value:
      dark_mode: false
      beta_ui: true
```

Apply migrations to your config:

```bash
patchwork apply --config config.yaml --migrations ./migrations/
```

Check migration status:

```bash
patchwork status --config config.yaml --migrations ./migrations/
```

Roll back the last applied migration:

```bash
patchwork rollback --config config.yaml --migrations ./migrations/
```

Patchwork tracks applied migrations in a `.patchwork_state` file alongside your config, making it safe to run in CI/CD pipelines and across multiple environments.

---

## Supported Formats

- JSON (`.json`)
- YAML (`.yaml`, `.yml`)

---

## License

MIT © 2024 yourname