package patch

import (
	"os/user"
	"time"
)

// AuditOptions controls audit recording behaviour during Run.
type AuditOptions struct {
	Enabled    bool
	ConfigPath string
}

// currentUser attempts to resolve the OS username, falling back to "unknown".
func currentUser() string {
	u, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return u.Username
}

// recordSuccess appends a successful application entry to the audit log.
func recordSuccess(opts AuditOptions, patchFile, version string, ops int) {
	if !opts.Enabled {
		return
	}
	_ = RecordAudit(opts.ConfigPath, AuditEntry{
		Timestamp: time.Now().UTC(),
		PatchFile: patchFile,
		Version:   version,
		Ops:       ops,
		User:      currentUser(),
		Status:    "applied",
	})
}

// recordFailure appends a failed application entry to the audit log.
func recordFailure(opts AuditOptions, patchFile, version string, ops int, msg string) {
	if !opts.Enabled {
		return
	}
	_ = RecordAudit(opts.ConfigPath, AuditEntry{
		Timestamp: time.Now().UTC(),
		PatchFile: patchFile,
		Version:   version,
		Ops:       ops,
		User:      currentUser(),
		Status:    "failed",
		Message:   msg,
	})
}

// recordRollback appends a rollback entry to the audit log.
func recordRollback(opts AuditOptions, patchFile, version string, ops int) {
	if !opts.Enabled {
		return
	}
	_ = RecordAudit(opts.ConfigPath, AuditEntry{
		Timestamp: time.Now().UTC(),
		PatchFile: patchFile,
		Version:   version,
		Ops:       ops,
		User:      currentUser(),
		Status:    "rolled_back",
	})
}
