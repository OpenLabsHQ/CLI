package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Template model structures.
type Template struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Name     string `json:"name"`
	VPN      bool   `json:"vpn"`
	VNC      bool   `json:"vnc"`
}

type TemplateID struct {
	ID string `json:"id"`
}

type VPCTemplate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CIDR string `json:"cidr"`
}

type SubnetTemplate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CIDR string `json:"cidr"`
}

type HostTemplate struct {
	ID       string   `json:"id"`
	Hostname string   `json:"hostname"`
	OS       string   `json:"os"`
	Spec     string   `json:"spec"`
	Size     int      `json:"size"`
	Tags     []string `json:"tags,omitempty"`
}

// Commands.
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Upload and manage templates",
	Long:  "This command will let you upload, view, and delete templates for ranges, VPCs, subnets, and hosts.",
}

// Range Template Commands.
var rangeTemplatesCmd = &cobra.Command{
	Use:   "range",
	Short: "Manage range templates",
	Long:  "Create, list, get, and delete range templates",
}

var listRangeTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all range templates",
	Long:  "This command will list all range templates from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		err := listRangeTemplates()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getRangeTemplateCmd = &cobra.Command{
	Use:   "get [range-id]",
	Short: "Get a range template",
	Long:  "This command will get a range template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := getRangeTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadRangeTemplateCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a range template",
	Long:  "This command will upload a range template to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadRangeTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteRangeTemplateCmd = &cobra.Command{
	Use:   "delete [range-id]",
	Short: "Delete a range template",
	Long:  "This command will delete a range template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteRangeTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// VPC Template Commands.
var vpcTemplatesCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Manage VPC templates",
	Long:  "Create, list, get, and delete VPC templates",
}

var listVPCTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPC templates",
	Long:  "This command will list all VPC templates from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")
		err := listVPCTemplates(standaloneOnly)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getVPCTemplateCmd = &cobra.Command{
	Use:   "get [vpc-id]",
	Short: "Get a VPC template",
	Long:  "This command will get a VPC template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := getVPCTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadVPCTemplateCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a VPC template",
	Long:  "This command will upload a VPC template to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadVPCTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteVPCTemplateCmd = &cobra.Command{
	Use:   "delete [vpc-id]",
	Short: "Delete a VPC template",
	Long:  "This command will delete a VPC template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteVPCTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Subnet Template Commands.
var subnetTemplatesCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Manage subnet templates",
	Long:  "Create, list, get, and delete subnet templates",
}

var listSubnetTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subnet templates",
	Long:  "This command will list all subnet templates from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")
		err := listSubnetTemplates(standaloneOnly)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getSubnetTemplateCmd = &cobra.Command{
	Use:   "get [subnet-id]",
	Short: "Get a subnet template",
	Long:  "This command will get a subnet template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := getSubnetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadSubnetTemplateCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a subnet template",
	Long:  "This command will upload a subnet template to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadSubnetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteSubnetTemplateCmd = &cobra.Command{
	Use:   "delete [subnet-id]",
	Short: "Delete a subnet template",
	Long:  "This command will delete a subnet template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteSubnetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Host Template Commands.
var hostTemplatesCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage host templates",
	Long:  "Create, list, get, and delete host templates",
}

var listHostTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all host templates",
	Long:  "This command will list all host templates from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		standaloneOnly, _ := cmd.Flags().GetBool("standalone")
		err := listHostTemplates(standaloneOnly)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getHostTemplateCmd = &cobra.Command{
	Use:   "get [host-id]",
	Short: "Get a host template",
	Long:  "This command will get a host template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := getHostTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var uploadHostTemplateCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a host template",
	Long:  "This command will upload a host template to the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadHostTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteHostTemplateCmd = &cobra.Command{
	Use:   "delete [host-id]",
	Short: "Delete a host template",
	Long:  "This command will delete a host template from the OpenLabs API.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteHostTemplate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Range Templates Implementation.
func listRangeTemplates() error {
	client := NewClient()
	resp, err := client.DoRequest("GET", "/api/v1/templates/ranges", nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var templates []Template
	if err := ParseResponse(resp, &templates); err != nil {
		return err
	}

	if len(templates) == 0 {
		fmt.Println("No range templates found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "Provider", "VNC", "VPN"})

	for _, t := range templates {
		table.Append([]string{
			t.Name,
			t.ID,
			t.Provider,
			fmt.Sprintf("%t", t.VNC),
			fmt.Sprintf("%t", t.VPN),
		})
	}

	table.Render()
	return nil
}

func getRangeTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/templates/ranges/%s", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var template map[string]interface{}
	if err := ParseResponse(resp, &template); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(template)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadRangeTemplate(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %s", err)
	}

	var templateData interface{}
	if err := json.Unmarshal(data, &templateData); err != nil {
		return fmt.Errorf("failed to parse template JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/templates/ranges", templateData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result TemplateID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Range template uploaded successfully!\n  ID: %s\n", result.ID)
	return nil
}

func deleteRangeTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/templates/ranges/%s", id), nil)
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
		fmt.Println("Range template deleted successfully")
	} else {
		fmt.Println("Failed to delete range template")
	}
	return nil
}

// VPC Templates Implementation.
func listVPCTemplates(standaloneOnly bool) error {
	client := NewClient()
	path := "/api/v1/templates/vpcs"
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

	var templates []VPCTemplate
	if err := ParseResponse(resp, &templates); err != nil {
		return err
	}

	if len(templates) == 0 {
		fmt.Println("No VPC templates found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "CIDR"})

	for _, t := range templates {
		table.Append([]string{
			t.Name,
			t.ID,
			t.CIDR,
		})
	}

	table.Render()
	return nil
}

func getVPCTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/templates/vpcs/%s", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var template map[string]interface{}
	if err := ParseResponse(resp, &template); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(template)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadVPCTemplate(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %s", err)
	}

	var templateData interface{}
	if err := json.Unmarshal(data, &templateData); err != nil {
		return fmt.Errorf("failed to parse template JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/templates/vpcs", templateData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result TemplateID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("VPC template uploaded successfully!\n  ID: %s\n", result.ID)
	return nil
}

func deleteVPCTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/templates/vpcs/%s", id), nil)
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
		fmt.Println("VPC template deleted successfully")
	} else {
		fmt.Println("Failed to delete VPC template")
	}
	return nil
}

// Subnet Templates Implementation.
func listSubnetTemplates(standaloneOnly bool) error {
	client := NewClient()
	path := "/api/v1/templates/subnets"
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

	var templates []SubnetTemplate
	if err := ParseResponse(resp, &templates); err != nil {
		return err
	}

	if len(templates) == 0 {
		fmt.Println("No subnet templates found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "CIDR"})

	for _, t := range templates {
		table.Append([]string{
			t.Name,
			t.ID,
			t.CIDR,
		})
	}

	table.Render()
	return nil
}

func getSubnetTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/templates/subnets/%s", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var template map[string]interface{}
	if err := ParseResponse(resp, &template); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(template)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadSubnetTemplate(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %s", err)
	}

	var templateData interface{}
	if err := json.Unmarshal(data, &templateData); err != nil {
		return fmt.Errorf("failed to parse template JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/templates/subnets", templateData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result TemplateID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Subnet template uploaded successfully!\n  ID: %s\n", result.ID)
	return nil
}

func deleteSubnetTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/templates/subnets/%s", id), nil)
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
		fmt.Println("Subnet template deleted successfully")
	} else {
		fmt.Println("Failed to delete subnet template")
	}
	return nil
}

// Host Templates Implementation.
func listHostTemplates(standaloneOnly bool) error {
	client := NewClient()
	path := "/api/v1/templates/hosts"
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

	var templates []HostTemplate
	if err := ParseResponse(resp, &templates); err != nil {
		return err
	}

	if len(templates) == 0 {
		fmt.Println("No host templates found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Hostname", "ID", "OS", "Spec", "Size", "Tags"})

	for _, t := range templates {
		table.Append([]string{
			t.Hostname,
			t.ID,
			t.OS,
			t.Spec,
			fmt.Sprintf("%d", t.Size),
			strings.Join(t.Tags, ", "),
		})
	}

	table.Render()
	return nil
}

func getHostTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/templates/hosts/%s", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var template map[string]interface{}
	if err := ParseResponse(resp, &template); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(template)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func uploadHostTemplate(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %s", err)
	}

	var templateData interface{}
	if err := json.Unmarshal(data, &templateData); err != nil {
		return fmt.Errorf("failed to parse template JSON: %s", err)
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/templates/hosts", templateData)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result TemplateID
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Host template uploaded successfully!\n  ID: %s\n", result.ID)
	return nil
}

func deleteHostTemplate(id string) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/templates/hosts/%s", id), nil)
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
		fmt.Println("Host template deleted successfully")
	} else {
		fmt.Println("Failed to delete host template")
	}
	return nil
}

func init() {
	// Setup range template subcommands
	listVPCTemplatesCmd.Flags().Bool("standalone", true, "List only standalone templates (not part of a range template)")
	listSubnetTemplatesCmd.Flags().Bool("standalone", true, "List only standalone templates (not part of a range/vpc template)")
	listHostTemplatesCmd.Flags().Bool("standalone", true, "List only standalone templates (not part of a range/vpc/subnet template)")

	// Range template commands
	rangeTemplatesCmd.AddCommand(listRangeTemplatesCmd)
	rangeTemplatesCmd.AddCommand(getRangeTemplateCmd)
	rangeTemplatesCmd.AddCommand(uploadRangeTemplateCmd)
	rangeTemplatesCmd.AddCommand(deleteRangeTemplateCmd)

	// VPC template commands
	vpcTemplatesCmd.AddCommand(listVPCTemplatesCmd)
	vpcTemplatesCmd.AddCommand(getVPCTemplateCmd)
	vpcTemplatesCmd.AddCommand(uploadVPCTemplateCmd)
	vpcTemplatesCmd.AddCommand(deleteVPCTemplateCmd)

	// Subnet template commands
	subnetTemplatesCmd.AddCommand(listSubnetTemplatesCmd)
	subnetTemplatesCmd.AddCommand(getSubnetTemplateCmd)
	subnetTemplatesCmd.AddCommand(uploadSubnetTemplateCmd)
	subnetTemplatesCmd.AddCommand(deleteSubnetTemplateCmd)

	// Host template commands
	hostTemplatesCmd.AddCommand(listHostTemplatesCmd)
	hostTemplatesCmd.AddCommand(getHostTemplateCmd)
	hostTemplatesCmd.AddCommand(uploadHostTemplateCmd)
	hostTemplatesCmd.AddCommand(deleteHostTemplateCmd)

	// Add all template subcommands to the templates command
	templatesCmd.AddCommand(rangeTemplatesCmd)
	templatesCmd.AddCommand(vpcTemplatesCmd)
	templatesCmd.AddCommand(subnetTemplatesCmd)
	templatesCmd.AddCommand(hostTemplatesCmd)

	// Add the templates command to the root command
	rootCmd.AddCommand(templatesCmd)
}
