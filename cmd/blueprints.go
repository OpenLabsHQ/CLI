package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Blueprint model structures.
type BlueprintHeader struct {
	ID          int    `json:"id"`
	Provider    string `json:"provider"`
	Name        string `json:"name"`
	VPN         bool   `json:"vpn"`
	VNC         bool   `json:"vnc"`
	Description string `json:"description,omitempty"`
}

type BlueprintID struct {
	ID int `json:"id"`
}

type VPCBlueprint struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	CIDR string `json:"cidr"`
}

type SubnetBlueprint struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	CIDR string `json:"cidr"`
}

type HostBlueprint struct {
	ID       int      `json:"id"`
	Hostname string   `json:"hostname"`
	OS       string   `json:"os"`
	Spec     string   `json:"spec"`
	Size     int      `json:"size"`
	Tags     []string `json:"tags,omitempty"`
}

// Commands.
var blueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Upload and manage blueprints",
	Long:  "This command will let you upload, view, and delete blueprints for ranges, VPCs, subnets, and hosts.",
}

// Range Blueprint Commands.
var rangeBlueprintsCmd = &cobra.Command{
	Use:   "range",
	Short: "Manage range blueprints",
	Long:  "Create, list, get, and delete range blueprints",
}

var listRangeBlueprintsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all range blueprints",
	Long:  "This command will list all range blueprints from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		err := listRangeBlueprints()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getRangeBlueprintCmd = &cobra.Command{
	Use:   "get [blueprint-id]",
	Short: "Get a range blueprint",
	Long:  "This command will get a range blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = getRangeBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadRangeBlueprintCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a range blueprint",
	Long:  "This command will upload a range blueprint to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadRangeBlueprint(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteRangeBlueprintCmd = &cobra.Command{
	Use:   "delete [blueprint-id]",
	Short: "Delete a range blueprint",
	Long:  "This command will delete a range blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = deleteRangeBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// VPC Blueprint Commands.
var vpcBlueprintsCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Manage VPC blueprints",
	Long:  "Create, list, get, and delete VPC blueprints",
}

var listVPCBlueprintsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPC blueprints",
	Long:  "This command will list all VPC blueprints from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")
		err := listVPCBlueprints(standaloneOnly)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getVPCBlueprintCmd = &cobra.Command{
	Use:   "get [blueprint-id]",
	Short: "Get a VPC blueprint",
	Long:  "This command will get a VPC blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = getVPCBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadVPCBlueprintCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a VPC blueprint",
	Long:  "This command will upload a VPC blueprint to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadVPCBlueprint(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteVPCBlueprintCmd = &cobra.Command{
	Use:   "delete [blueprint-id]",
	Short: "Delete a VPC blueprint",
	Long:  "This command will delete a VPC blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = deleteVPCBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Subnet Blueprint Commands.
var subnetBlueprintsCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Manage subnet blueprints",
	Long:  "Create, list, get, and delete subnet blueprints",
}

var listSubnetBlueprintsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subnet blueprints",
	Long:  "This command will list all subnet blueprints from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")
		err := listSubnetBlueprints(standaloneOnly)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getSubnetBlueprintCmd = &cobra.Command{
	Use:   "get [blueprint-id]",
	Short: "Get a subnet blueprint",
	Long:  "This command will get a subnet blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = getSubnetBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadSubnetBlueprintCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a subnet blueprint",
	Long:  "This command will upload a subnet blueprint to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadSubnetBlueprint(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteSubnetBlueprintCmd = &cobra.Command{
	Use:   "delete [blueprint-id]",
	Short: "Delete a subnet blueprint",
	Long:  "This command will delete a subnet blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = deleteSubnetBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Host Blueprint Commands.
var hostBlueprintsCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage host blueprints",
	Long:  "Create, list, get, and delete host blueprints",
}

var listHostBlueprintsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all host blueprints",
	Long:  "This command will list all host blueprints from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")
		err := listHostBlueprints(standaloneOnly)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getHostBlueprintCmd = &cobra.Command{
	Use:   "get [blueprint-id]",
	Short: "Get a host blueprint",
	Long:  "This command will get a host blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = getHostBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadHostBlueprintCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a host blueprint",
	Long:  "This command will upload a host blueprint to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadHostBlueprint(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteHostBlueprintCmd = &cobra.Command{
	Use:   "delete [blueprint-id]",
	Short: "Delete a host blueprint",
	Long:  "This command will delete a host blueprint from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}
		err = deleteHostBlueprint(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Range Blueprints Implementation.
func listRangeBlueprints() error {
	client := NewClient()
	resp, err := client.DoRequest("GET", "/api/v1/blueprints/ranges", nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprints []BlueprintHeader
	if err := ParseResponse(resp, &blueprints); err != nil {
		return err
	}

	if len(blueprints) == 0 {
		fmt.Println("No range blueprints found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "Provider", "VNC", "VPN", "Description"})

	for _, t := range blueprints {
		table.Append([]string{
			t.Name,
			strconv.Itoa(t.ID),
			t.Provider,
			fmt.Sprintf("%t", t.VNC),
			fmt.Sprintf("%t", t.VPN),
			t.Description,
		})
	}

	table.Render()
	return nil
}

func getRangeBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/blueprints/ranges/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprint map[string]interface{}
	if err := ParseResponse(resp, &blueprint); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(blueprint)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadRangeBlueprint(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read blueprint file: %s", err)
	}

	var blueprintData interface{}
	if err := json.Unmarshal(data, &blueprintData); err != nil {
		return fmt.Errorf("failed to parse blueprint JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/blueprints/ranges", blueprintData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result BlueprintHeader
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Range blueprint uploaded successfully!\n  ID: %d\n", result.ID)
	return nil
}

func deleteRangeBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/blueprints/ranges/%d", id), nil)
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
		fmt.Println("Range blueprint deleted successfully")
	} else {
		fmt.Println("Failed to delete range blueprint")
	}
	return nil
}

// VPC Blueprints Implementation.
func listVPCBlueprints(standaloneOnly bool) error {
	client := NewClient()
	path := "/api/v1/blueprints/vpcs"
	if !standaloneOnly {
		path += "?standalone_only=false"
	}

	resp, err := client.DoRequest("GET", path, nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprints []VPCBlueprint
	if err := ParseResponse(resp, &blueprints); err != nil {
		return err
	}

	if len(blueprints) == 0 {
		fmt.Println("No VPC blueprints found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "CIDR"})

	for _, t := range blueprints {
		table.Append([]string{
			t.Name,
			strconv.Itoa(t.ID),
			t.CIDR,
		})
	}

	table.Render()
	return nil
}

func getVPCBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/blueprints/vpcs/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprint map[string]interface{}
	if err := ParseResponse(resp, &blueprint); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(blueprint)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadVPCBlueprint(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read blueprint file: %s", err)
	}

	var blueprintData interface{}
	if err := json.Unmarshal(data, &blueprintData); err != nil {
		return fmt.Errorf("failed to parse blueprint JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/blueprints/vpcs", blueprintData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result BlueprintID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("VPC blueprint uploaded successfully!\n  ID: %d\n", result.ID)
	return nil
}

func deleteVPCBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/blueprints/vpcs/%d", id), nil)
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
		fmt.Println("VPC blueprint deleted successfully")
	} else {
		fmt.Println("Failed to delete VPC blueprint")
	}
	return nil
}

// Subnet Blueprints Implementation.
func listSubnetBlueprints(standaloneOnly bool) error {
	client := NewClient()
	path := "/api/v1/blueprints/subnets"
	if !standaloneOnly {
		path += "?standalone_only=false"
	}

	resp, err := client.DoRequest("GET", path, nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprints []SubnetBlueprint
	if err := ParseResponse(resp, &blueprints); err != nil {
		return err
	}

	if len(blueprints) == 0 {
		fmt.Println("No subnet blueprints found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "CIDR"})

	for _, t := range blueprints {
		table.Append([]string{
			t.Name,
			strconv.Itoa(t.ID),
			t.CIDR,
		})
	}

	table.Render()
	return nil
}

func getSubnetBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/blueprints/subnets/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprint map[string]interface{}
	if err := ParseResponse(resp, &blueprint); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(blueprint)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadSubnetBlueprint(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read blueprint file: %s", err)
	}

	var blueprintData interface{}
	if err := json.Unmarshal(data, &blueprintData); err != nil {
		return fmt.Errorf("failed to parse blueprint JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/blueprints/subnets", blueprintData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result BlueprintID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Subnet blueprint uploaded successfully!\n  ID: %d\n", result.ID)
	return nil
}

func deleteSubnetBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/blueprints/subnets/%d", id), nil)
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
		fmt.Println("Subnet blueprint deleted successfully")
	} else {
		fmt.Println("Failed to delete subnet blueprint")
	}
	return nil
}

// Host Blueprints Implementation.
func listHostBlueprints(standaloneOnly bool) error {
	client := NewClient()
	path := "/api/v1/blueprints/hosts"
	if !standaloneOnly {
		path += "?standalone_only=false"
	}

	resp, err := client.DoRequest("GET", path, nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprints []HostBlueprint
	if err := ParseResponse(resp, &blueprints); err != nil {
		return err
	}

	if len(blueprints) == 0 {
		fmt.Println("No host blueprints found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Hostname", "ID", "OS", "Spec", "Size", "Tags"})

	for _, t := range blueprints {
		table.Append([]string{
			t.Hostname,
			strconv.Itoa(t.ID),
			t.OS,
			t.Spec,
			fmt.Sprintf("%d", t.Size),
			strings.Join(t.Tags, ", "),
		})
	}

	table.Render()
	return nil
}

func getHostBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/blueprints/hosts/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprint map[string]interface{}
	if err := ParseResponse(resp, &blueprint); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(blueprint)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadHostBlueprint(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read blueprint file: %s", err)
	}

	var blueprintData interface{}
	if err := json.Unmarshal(data, &blueprintData); err != nil {
		return fmt.Errorf("failed to parse blueprint JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/blueprints/hosts", blueprintData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result BlueprintID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Host blueprint uploaded successfully!\n  ID: %d\n", result.ID)
	return nil
}

func deleteHostBlueprint(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/blueprints/hosts/%d", id), nil)
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
		fmt.Println("Host blueprint deleted successfully")
	} else {
		fmt.Println("Failed to delete host blueprint")
	}
	return nil
}

func init() {
	// Setup range blueprint subcommands
	listVPCBlueprintsCmd.Flags().Bool("standalone", true, "List only standalone blueprints (not part of a range blueprint)")
	listSubnetBlueprintsCmd.Flags().Bool("standalone", true, "List only standalone blueprints (not part of a range/vpc blueprint)")
	listHostBlueprintsCmd.Flags().Bool("standalone", true, "List only standalone blueprints (not part of a range/vpc/subnet blueprint)")

	// Range blueprint commands
	rangeBlueprintsCmd.AddCommand(listRangeBlueprintsCmd)
	rangeBlueprintsCmd.AddCommand(getRangeBlueprintCmd)
	rangeBlueprintsCmd.AddCommand(uploadRangeBlueprintCmd)
	rangeBlueprintsCmd.AddCommand(deleteRangeBlueprintCmd)

	// VPC blueprint commands
	vpcBlueprintsCmd.AddCommand(listVPCBlueprintsCmd)
	vpcBlueprintsCmd.AddCommand(getVPCBlueprintCmd)
	vpcBlueprintsCmd.AddCommand(uploadVPCBlueprintCmd)
	vpcBlueprintsCmd.AddCommand(deleteVPCBlueprintCmd)

	// Subnet blueprint commands
	subnetBlueprintsCmd.AddCommand(listSubnetBlueprintsCmd)
	subnetBlueprintsCmd.AddCommand(getSubnetBlueprintCmd)
	subnetBlueprintsCmd.AddCommand(uploadSubnetBlueprintCmd)
	subnetBlueprintsCmd.AddCommand(deleteSubnetBlueprintCmd)

	// Host blueprint commands
	hostBlueprintsCmd.AddCommand(listHostBlueprintsCmd)
	hostBlueprintsCmd.AddCommand(getHostBlueprintCmd)
	hostBlueprintsCmd.AddCommand(uploadHostBlueprintCmd)
	hostBlueprintsCmd.AddCommand(deleteHostBlueprintCmd)

	// Add all blueprint subcommands to the blueprints command
	blueprintsCmd.AddCommand(rangeBlueprintsCmd)
	blueprintsCmd.AddCommand(vpcBlueprintsCmd)
	blueprintsCmd.AddCommand(subnetBlueprintsCmd)
	blueprintsCmd.AddCommand(hostBlueprintsCmd)

	// Add the blueprints command to the root command
	rootCmd.AddCommand(blueprintsCmd)
	
	// Also keep the templates command for backward compatibility
	var templatesCmd = &cobra.Command{
		Use:   "templates",
		Short: "Alias for blueprints (deprecated)",
		Long:  "This is an alias for the blueprints command. Please use 'blueprints' instead as this command will be removed in the future.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Warning: The 'templates' command is deprecated and will be removed in a future version.")
			fmt.Println("Please use the 'blueprints' command instead.")
			if err := blueprintsCmd.Help(); err != nil {
				fmt.Println("Error displaying help:", err)
			}
		},
	}
	rootCmd.AddCommand(templatesCmd)
}
