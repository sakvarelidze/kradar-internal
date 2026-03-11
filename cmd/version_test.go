package cmd

import "testing"

func TestWithBuildInfoDefaults_UsesBuildSettingsForDefaultValues(t *testing.T) {
	version, commit, date := withBuildInfoDefaults("dev", "none", "unknown", map[string]string{
		"vcs.revision": "abc123",
		"vcs.time":     "2026-03-10T01:02:03Z",
	})

	if version != "abc123" {
		t.Fatalf("expected version from vcs.revision, got %q", version)
	}
	if commit != "abc123" {
		t.Fatalf("expected commit from vcs.revision, got %q", commit)
	}
	if date != "2026-03-10T01:02:03Z" {
		t.Fatalf("expected date from vcs.time, got %q", date)
	}
}

func TestWithBuildInfoDefaults_AppendsDirtySuffixForModifiedTree(t *testing.T) {
	version, _, _ := withBuildInfoDefaults("dev", "none", "unknown", map[string]string{
		"vcs.revision": "abc123",
		"vcs.modified": "true",
	})

	if version != "abc123-dirty" {
		t.Fatalf("expected dirty version suffix, got %q", version)
	}
}

func TestWithBuildInfoDefaults_DoesNotOverrideExplicitValues(t *testing.T) {
	version, commit, date := withBuildInfoDefaults("v1.2.3", "deadbeef", "2026-03-10", map[string]string{
		"vcs.revision": "abc123",
		"vcs.time":     "2026-03-10T01:02:03Z",
	})

	if version != "v1.2.3" || commit != "deadbeef" || date != "2026-03-10" {
		t.Fatalf("expected explicit values to be preserved, got version=%q commit=%q date=%q", version, commit, date)
	}
}
