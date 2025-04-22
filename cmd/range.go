package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Structures for range commands.
type DeployRangeRequest struct {
	Name        string `json:"name"`
	BlueprintID int    `json:"blueprint_id"`
	Region      string `json:"region"`
	Description string `json:"description,omitempty"`
}

type DeployedRangeHeader struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	BlueprintID int       `json:"blueprint_id"`
	State       string    `json:"state"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Range Commands.
var rangeCmd = &cobra.Command{
	Use:   "range",
	Short: "Deploy and manage ranges",
	Long:  "This command will let you deploy, power on/off, and manage your ranges.",
}

var listRangesCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployed ranges",
	Long:  "This command will list all your deployed ranges.",
	Run: func(cmd *cobra.Command, args []string) {
		err := listRanges()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getRangeCmd = &cobra.Command{
	Use:   "get [range-id]",
	Short: "Get a deployed range",
	Long:  "This command will get details of a deployed range.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: range ID must be a number")
			return
		}
		err = getRange(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deployRangeCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a range",
	Long:  "This command will deploy a range from a blueprint to the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		blueprintID, _ := cmd.Flags().GetInt("blueprint-id")
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		description, _ := cmd.Flags().GetString("description")

		if blueprintID == 0 || name == "" || region == "" {
			fmt.Println("Error: --blueprint-id, --name, and --region are required")
			return
		}

		err := deployRange(blueprintID, name, region, description)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteRangeCmd = &cobra.Command{
	Use:   "delete [range-id]",
	Short: "Delete a deployed range",
	Long:  "This command will delete a deployed range.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: range ID must be a number")
			return
		}
		err = deleteRange(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Ranges Implementation.
func listRanges() error {
	client := NewClient()
	resp, err := client.DoRequest("GET", "/api/v1/ranges", nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var ranges []DeployedRangeHeader
	if err := ParseResponse(resp, &ranges); err != nil {
		return err
	}

	if len(ranges) == 0 {
		fmt.Println("No deployed ranges found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "State", "Created At"})

	for _, r := range ranges {
		table.Append([]string{
			strconv.Itoa(r.ID),
			r.Name,
			r.Description,
			r.State,
			r.CreatedAt.Format(time.RFC3339),
		})
	}

	table.Render()
	return nil
}

func getRange(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/ranges/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var deployedRange map[string]interface{}
	if err := ParseResponse(resp, &deployedRange); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(deployedRange)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func deployRange(blueprintID int, name, region, description string) error {
	request := DeployRangeRequest{
		BlueprintID: blueprintID,
		Name:        name,
		Region:      region,
		Description: description,
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/ranges/deploy", request)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

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

func deleteRange(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/ranges/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result bool
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	if result {
		fmt.Println("Range deleted successfully")
	} else {
		fmt.Println("Failed to delete range")
	}
	return nil
}

func init() {
	// Deploy command flags
	deployRangeCmd.Flags().Int("blueprint-id", 0, "ID of the blueprint to deploy")
	deployRangeCmd.Flags().String("name", "", "Name for the deployed range")
	deployRangeCmd.Flags().String("region", "", "Region to deploy the range in (e.g., us_east_1)")
	deployRangeCmd.Flags().String("description", "", "Optional description for the range")

	// Add subcommands to range command
	rangeCmd.AddCommand(listRangesCmd)
	rangeCmd.AddCommand(getRangeCmd)
	rangeCmd.AddCommand(deployRangeCmd)
	rangeCmd.AddCommand(deleteRangeCmd)
	
	// Add range command to root
	rootCmd.AddCommand(rangeCmd)
}
