package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "OpenLabs",
    Short: "OpenLabs CLI",
    Long: "openlabs",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0{
            cmd.Help()
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
