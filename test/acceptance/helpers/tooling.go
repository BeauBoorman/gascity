package acceptancehelpers

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	// EnvRequireTooling makes a missing bd or dolt fatal instead of a skip.
	EnvRequireTooling = "GC_REQUIRE_ACCEPTANCE_TOOLING"
	// EnvRequireLegacyGC makes a missing pre-journal gc fatal instead of a skip.
	EnvRequireLegacyGC = "GC_REQUIRE_ACCEPTANCE_LEGACY_GC"
)

// requireSwitchOn reports whether a GC_REQUIRE_* switch is on. Unset, empty and
// "0" are off; anything else is on.
func requireSwitchOn(name string) bool {
	v := strings.TrimSpace(os.Getenv(name))
	return v != "" && v != "0"
}

// MissingTooling reports a dependency the beads acceptance shapes cannot run
// without. It skips by default so a developer without bd or dolt on PATH still
// gets a useful local run, and fails under EnvRequireTooling.
//
// The switch exists because the skip is indistinguishable from a pass in a job
// summary. Every one of these tests skipped in every CI job for want of a bd,
// so a suite that had never executed read as green for as long as it took
// someone to check. A job whose whole purpose is to run them sets the switch,
// and a runner that loses its bd fails loudly instead of quietly running
// nothing.
func MissingTooling(t *testing.T, format string, args ...any) {
	t.Helper()
	skipOrFail(t, EnvRequireTooling, fmt.Sprintf(format, args...))
}

// MissingLegacyGC is the same contract for the pre-journal gc binary, on its
// own switch. The migration tests and the M5 shape need a gc from before the
// ownership journal, which no job builds yet; a separate switch lets CI demand
// bd and dolt without also demanding a binary nothing produces.
func MissingLegacyGC(t *testing.T, format string, args ...any) {
	t.Helper()
	skipOrFail(t, EnvRequireLegacyGC, fmt.Sprintf(format, args...))
}

func skipOrFail(t *testing.T, env, reason string) {
	t.Helper()
	if requireSwitchOn(env) {
		t.Fatalf("%s is set, so this test must run, but %s", env, reason)
	}
	t.Skip(reason)
}
