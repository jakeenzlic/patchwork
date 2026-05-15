package patch

import (
	"fmt"
	"regexp"
	"strings"
)

// envVarPattern matches ${VAR_NAME} style placeholders.
var envVarPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// EnvResolver resolves environment variable placeholders in patch operation values.
type EnvResolver struct {
	lookup func(string) (string, bool)
}

// NewEnvResolver creates an EnvResolver using the provided lookup function.
// Pass os.LookupEnv for production use.
func NewEnvResolver(lookup func(string) (string, bool)) *EnvResolver {
	return &EnvResolver{lookup: lookup}
}

// ResolvePatches returns a copy of the patches with all ${VAR} placeholders
// in string values substituted. Returns an error if any variable is unset.
func (r *EnvResolver) ResolvePatches(patches []PatchEntry) ([]PatchEntry, error) {
	resolved := make([]PatchEntry, len(patches))
	for i, p := range patches {
		cp := p
		if s, ok := p.Value.(string); ok {
			val, err := r.resolveString(s)
			if err != nil {
				return nil, fmt.Errorf("patch %d path %q: %w", i, p.Path, err)
			}
			cp.Value = val
		}
		resolved[i] = cp
	}
	return resolved, nil
}

// resolveString substitutes all ${VAR} occurrences in s.
func (r *EnvResolver) resolveString(s string) (string, error) {
	var firstErr error
	result := envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		if firstErr != nil {
			return match
		}
		name := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		val, ok := r.lookup(name)
		if !ok {
			firstErr = fmt.Errorf("environment variable %q is not set", name)
			return match
		}
		return val
	})
	if firstErr != nil {
		return "", firstErr
	}
	return result, nil
}
