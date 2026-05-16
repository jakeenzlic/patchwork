// Package patch provides compress/bundle utilities for patchwork.
//
// Bundle and Unbundle allow a set of patches to be packed into a single
// gzip-compressed JSON archive (*.patch.gz).  This is useful for shipping
// a complete migration set as a single artefact, e.g. as part of a CI
// release pipeline.
//
// Typical usage:
//
//	// Pack all patches in a directory into a release bundle.
//	patches, _ := patch.LoadDir("./migrations")
//	patch.Bundle(patches, "dist/v2.patch.gz")
//
//	// Later, on the target host, unpack and apply.
//	patches, _ := patch.Unbundle("dist/v2.patch.gz")
//	for _, p := range patches {
//		patch.Run(p, "config.json", "config.json", false, false)
//	}
package patch
