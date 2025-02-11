package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
    Use:   "plugins",
    Short: "View and deploy plugins",
    Long:  "This command will let view plugins and deploy them to your range.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Placeholder")
    },
}

func init() {
    rootCmd.AddCommand(pluginsCmd)
}
