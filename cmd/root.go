package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information will be populated by the build process.
var (
	version   = ""
	buildTime = "unknown"
)

// Configuration variables - exported for use across package files
var (
	APIURL    string
	AuthToken string
	EncKey    string
	Debug     bool
)

var rootCmd = &cobra.Command{
	Use:   "openlabs",
	Short: "A command line interface for managing OpenLabs",
	Long:  "OpenLabs CLI is a command line interface for managing OpenLabs and its associated blueprints, ranges, workspaces, and plugins.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				fmt.Println("Error displaying help:", err)
			}
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print detailed version and build information for the OpenLabs CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OpenLabs CLI v%s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&APIURL, "api-url", "http://localhost:8000", "URL of the OpenLabs API server")
	rootCmd.PersistentFlags().StringVar(&AuthToken, "token", "", "Authentication token for OpenLabs API")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable debug mode to see detailed request/response information")

	rootCmd.AddCommand(versionCmd)
}
