package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Structures for range commands.
type RangeDeployRequest struct {
	ID string `json:"id"`
}

// Range Commands.
var rangeCmd = &cobra.Command{
	Use:   "range",
	Short: "Deploy and manage ranges",
	Long:  "This command will let you deploy, power on/off, and manage your ranges.",
}

var deployRangeCmd = &cobra.Command{
	Use:   "deploy [range-id...]",
	Short: "Deploy a range",
	Long:  "This command will deploy one or more ranges to the OpenLabs API.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deployRange(args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Ranges Implementation.
func deployRange(templateIDs []string) error {
	// Convert each templateID to a RangeDeployRequest
	var requests []RangeDeployRequest
	for _, id := range templateIDs {
		requests = append(requests, RangeDeployRequest{ID: id})
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/ranges/deploy", requests)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Response is a deployment status object
	var result map[string]interface{}
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(result)
	if err != nil {
		return err
	}

	fmt.Println("Range deployment initiated successfully")
	fmt.Println("Deployment status:")
	fmt.Println(prettyJSON)

	return nil
}

func init() {
	rangeCmd.AddCommand(deployRangeCmd)
	rootCmd.AddCommand(rangeCmd)
}
