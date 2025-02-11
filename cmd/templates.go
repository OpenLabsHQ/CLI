package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
    Use:   "templates",
    Short: "Upload and manage templates",
    Long:  "This command will let you upload, view, and delete templates for ranges, subnets, and hosts.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Placeholder")
    },
}

func init() {
    rootCmd.AddCommand(templatesCmd)
}
