package patch

import (
	"fmt"
	"regexp"
	"strings"
)

// templateVarRe matches {{VAR_NAME}} style placeholders.
var templateVarRe = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)

// TemplateContext holds named values used to render patch templates.
type TemplateContext map[string]string

// RenderTemplate replaces all {{VAR}} placeholders in s using ctx.
// Returns an error if any placeholder is not found in ctx.
func RenderTemplate(s string, ctx TemplateContext) (string, error) {
	var missing []string
	result := templateVarRe.ReplaceAllStringFunc(s, func(match string) string {
		key := match[2 : len(match)-2] // strip {{ and }}
		if val, ok := ctx[key]; ok {
			return val
		}
		missing = append(missing, key)
		return match
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("template: unresolved variables: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// RenderPatches applies template substitution to all string-valued fields
// (Path and Value) in each operation within the patch.
func RenderPatches(p Patch, ctx TemplateContext) (Patch, error) {
	out := p
	out.Ops = make([]Op, len(p.Ops))
	for i, op := range p.Ops {
		renderedPath, err := RenderTemplate(op.Path, ctx)
		if err != nil {
			return Patch{}, fmt.Errorf("op[%d] path: %w", i, err)
		}
		op.Path = renderedPath

		if strVal, ok := op.Value.(string); ok {
			renderedVal, err := RenderTemplate(strVal, ctx)
			if err != nil {
				return Patch{}, fmt.Errorf("op[%d] value: %w", i, err)
			}
			op.Value = renderedVal
		}
		out.Ops[i] = op
	}
	return out, nil
}

// ExtractTemplateVars returns the set of unique variable names referenced
// across all ops in the patch (paths and string values).
func ExtractTemplateVars(p Patch) []string {
	seen := map[string]struct{}{}
	for _, op := range p.Ops {
		for _, m := range templateVarRe.FindAllStringSubmatch(op.Path, -1) {
			seen[m[1]] = struct{}{}
		}
		if strVal, ok := op.Value.(string); ok {
			for _, m := range templateVarRe.FindAllStringSubmatch(strVal, -1) {
				seen[m[1]] = struct{}{}
			}
		}
	}
	vars := make([]string, 0, len(seen))
	for k := range seen {
		vars = append(vars, k)
	}
	return vars
}
