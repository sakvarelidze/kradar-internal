package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

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

	settings := make(map[string]string, len(info.Settings))
	for _, setting := range info.Settings {
		settings[setting.Key] = setting.Value
	}

	return settings
}

func withBuildInfoDefaults(version, commit, date string, settings map[string]string) (string, string, string) {
	if settings == nil {
		return version, commit, date
	}

	if version == "dev" {
		if vcsModified := settings["vcs.modified"]; vcsModified == "true" {
			if vcsRevision := settings["vcs.revision"]; vcsRevision != "" {
				version = vcsRevision + "-dirty"
			}
		} else if vcsRevision := settings["vcs.revision"]; vcsRevision != "" {
			version = vcsRevision
		}
	}

	if commit == "none" {
		if vcsRevision := settings["vcs.revision"]; vcsRevision != "" {
			commit = vcsRevision
		}
	}

	if date == "unknown" {
		if vcsTime := settings["vcs.time"]; vcsTime != "" {
			date = vcsTime
		}
	}

	return version, commit, date
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Run: func(_ *cobra.Command, _ []string) {
		_, _ = fmt.Fprintf(os.Stdout, "kradar version=%s commit=%s date=%s\n", Version, Commit, Date)
	},
}
