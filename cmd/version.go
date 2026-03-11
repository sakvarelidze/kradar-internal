package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func init() {
	Version, Commit, Date = withBuildInfoDefaults(Version, Commit, Date, readBuildSettings())
	rootCmd.AddCommand(versionCmd)
}

func readBuildSettings() map[string]string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	settings := make(map[string]string, len(info.Settings)+1)
	for _, setting := range info.Settings {
		settings[setting.Key] = setting.Value
	}
	settings["main.version"] = info.Main.Version

	return settings
}

func withBuildInfoDefaults(version, commit, date string, settings map[string]string) (string, string, string) {
	if settings == nil {
		return version, commit, date
	}

	if version == "dev" {
		version = semverLikeVersion(settings)
	}

	if commit == "none" {
		if vcsRevision := settings["vcs.revision"]; vcsRevision != "" {
			commit = shortSHA(vcsRevision)
		}
	}

	if date == "unknown" {
		if vcsTime := settings["vcs.time"]; vcsTime != "" {
			date = vcsTime
		}
	}

	return version, commit, date
}

func semverLikeVersion(settings map[string]string) string {
	if mv := strings.TrimSpace(settings["main.version"]); mv != "" && mv != "(devel)" {
		if v, err := semver.NewVersion(strings.TrimPrefix(mv, "v")); err == nil {
			return v.String()
		}
	}

	version := "0.0.0-dev"
	if settings["vcs.modified"] == "true" {
		version += "+dirty"
	}
	return version
}

func shortSHA(revision string) string {
	revision = strings.TrimSpace(revision)
	if len(revision) <= 7 {
		return revision
	}
	return revision[:7]
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Run: func(_ *cobra.Command, _ []string) {
		_, _ = fmt.Fprintf(os.Stdout, "kradar version=%s commit=%s date=%s\n", Version, Commit, Date)
	},
}
