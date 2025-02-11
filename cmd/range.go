package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var rangeCmd = &cobra.Command{
    Use:   "range",
    Short: "Deploy, and manage range",
    Long:  "This command will let you deploy, power on/off, and manage your range.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Placeholder")
        //client := api.Client{BaseURL: "https://localhost/api/v1/templates"}
        //data, err := client.
    },
}

func init() {
    rootCmd.AddCommand(rangeCmd)
}

