// Package patch — redact module
//
// # Redact
//
// Redact produces a sanitised copy of a config map by replacing the values at
// specified paths with a configurable mask string (default "***").
//
// Usage:
//
//	rules := []patch.RedactRule{
//		{Path: "database/password"},
//		{Path: "api/secret_key", MaskWith: "<hidden>"},
//	}
//	result := patch.Redact(config, rules)
//	fmt.Print(patch.FormatRedactResult(result))
//
// The original config map is never mutated; Redact operates on a deep copy.
//
// Paths follow the same slash-separated convention used throughout patchwork
// (e.g. "app/database/host"). Paths that do not exist in the config are
// silently skipped and will not appear in RedactResult.Masked.
//
// The redacted config can be exported to JSON or YAML via patch.Export.
package patch
