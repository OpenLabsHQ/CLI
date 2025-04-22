package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Workspace model structures.
type Workspace struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DefaultTimeLimit int       `json:"default_time_limit"`
	OwnerID         int       `json:"owner_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WorkspaceCreate struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	DefaultTimeLimit int    `json:"default_time_limit,omitempty"`
}

type WorkspaceUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	TimeLimit int    `json:"time_limit"`
}

type WorkspaceUserCreate struct {
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
	TimeLimit int    `json:"time_limit,omitempty"`
}

type WorkspaceUserUpdate struct {
	Role      string `json:"role,omitempty"`
	TimeLimit int    `json:"time_limit,omitempty"`
}

type WorkspaceBlueprint struct {
	BlueprintID   int    `json:"blueprint_id"`
	BlueprintType string `json:"blueprint_type"`
	Permission    string `json:"permission"`
}

// Workspace Commands.
var workspacesCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Long:  "This command will let you create, list, and manage workspaces.",
}

var listWorkspacesCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspaces",
	Long:  "This command will list all workspaces you have access to.",
	Run: func(cmd *cobra.Command, args []string) {
		err := listWorkspaces()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var getWorkspaceCmd = &cobra.Command{
	Use:   "get [workspace-id]",
	Short: "Get a workspace",
	Long:  "This command will get details of a workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}
		err = getWorkspace(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var createWorkspaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workspace",
	Long:  "This command will create a new workspace.",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		timeLimit, _ := cmd.Flags().GetInt("time-limit")

		if name == "" {
			fmt.Println("Error: --name is required")
			return
		}

		err := createWorkspace(name, description, timeLimit)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var deleteWorkspaceCmd = &cobra.Command{
	Use:   "delete [workspace-id]",
	Short: "Delete a workspace",
	Long:  "This command will delete a workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}
		err = deleteWorkspace(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Workspace Users Commands
var listWorkspaceUsersCmd = &cobra.Command{
	Use:   "list-users [workspace-id]",
	Short: "List workspace users",
	Long:  "This command will list all users in a workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}
		err = listWorkspaceUsers(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var addWorkspaceUserCmd = &cobra.Command{
	Use:   "add-user [workspace-id]",
	Short: "Add a user to a workspace",
	Long:  "This command will add a user to a workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}

		userID, _ := cmd.Flags().GetInt("user-id")
		role, _ := cmd.Flags().GetString("role")
		timeLimit, _ := cmd.Flags().GetInt("time-limit")

		if userID == 0 || role == "" {
			fmt.Println("Error: --user-id and --role are required")
			return
		}

		err = addWorkspaceUser(workspaceID, userID, role, timeLimit)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var updateWorkspaceUserCmd = &cobra.Command{
	Use:   "update-user [workspace-id] [user-id]",
	Short: "Update a user in a workspace",
	Long:  "This command will update a user's role or time limit in a workspace.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: user ID must be a number")
			return
		}

		role, _ := cmd.Flags().GetString("role")
		timeLimit, _ := cmd.Flags().GetInt("time-limit")

		if role == "" && timeLimit == 0 {
			fmt.Println("Error: at least one of --role or --time-limit must be specified")
			return
		}

		err = updateWorkspaceUser(workspaceID, userID, role, timeLimit)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var removeWorkspaceUserCmd = &cobra.Command{
	Use:   "remove-user [workspace-id] [user-id]",
	Short: "Remove a user from a workspace",
	Long:  "This command will remove a user from a workspace.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: user ID must be a number")
			return
		}

		err = removeWorkspaceUser(workspaceID, userID)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Workspace Blueprints Commands
var listWorkspaceBlueprintsCmd = &cobra.Command{
	Use:   "list-blueprints [workspace-id]",
	Short: "List workspace blueprints",
	Long:  "This command will list all blueprints shared with a workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}
		err = listWorkspaceBlueprints(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var addWorkspaceBlueprintCmd = &cobra.Command{
	Use:   "add-blueprint [workspace-id]",
	Short: "Share a blueprint with a workspace",
	Long:  "This command will share a blueprint with a workspace.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}

		blueprintID, _ := cmd.Flags().GetInt("blueprint-id")
		blueprintType, _ := cmd.Flags().GetString("blueprint-type")
		permission, _ := cmd.Flags().GetString("permission")

		if blueprintID == 0 || blueprintType == "" || permission == "" {
			fmt.Println("Error: --blueprint-id, --blueprint-type, and --permission are required")
			return
		}

		err = addWorkspaceBlueprint(workspaceID, blueprintID, blueprintType, permission)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var removeWorkspaceBlueprintCmd = &cobra.Command{
	Use:   "remove-blueprint [workspace-id] [blueprint-id]",
	Short: "Remove a blueprint from a workspace",
	Long:  "This command will remove a blueprint from a workspace.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: workspace ID must be a number")
			return
		}

		blueprintID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: blueprint ID must be a number")
			return
		}

		blueprintType, _ := cmd.Flags().GetString("blueprint-type")
		if blueprintType == "" {
			fmt.Println("Error: --blueprint-type is required")
			return
		}

		err = removeWorkspaceBlueprint(workspaceID, blueprintID, blueprintType)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Workspace Implementation.
func listWorkspaces() error {
	client := NewClient()
	resp, err := client.DoRequest("GET", "/api/v1/workspaces", nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var workspaces []Workspace
	if err := ParseResponse(resp, &workspaces); err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Default Time Limit", "Created At"})

	for _, w := range workspaces {
		table.Append([]string{
			strconv.Itoa(w.ID),
			w.Name,
			w.Description,
			fmt.Sprintf("%d seconds", w.DefaultTimeLimit),
			w.CreatedAt.Format(time.RFC3339),
		})
	}

	table.Render()
	return nil
}

func getWorkspace(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/workspaces/%d", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var workspace Workspace
	if err := ParseResponse(resp, &workspace); err != nil {
		return err
	}

	prettyJSON, err := FormatResponse(workspace)
	if err != nil {
		return err
	}

	fmt.Println(prettyJSON)
	return nil
}

func createWorkspace(name, description string, timeLimit int) error {
	request := WorkspaceCreate{
		Name:            name,
		Description:     description,
	}

	if timeLimit > 0 {
		request.DefaultTimeLimit = timeLimit
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/workspaces", request)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var workspace Workspace
	if err := ParseResponse(resp, &workspace); err != nil {
		return err
	}

	fmt.Printf("Workspace created successfully!\n  ID: %d\n  Name: %s\n", workspace.ID, workspace.Name)
	return nil
}

func deleteWorkspace(id int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/workspaces/%d", id), nil)
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
		fmt.Println("Workspace deleted successfully")
	} else {
		fmt.Println("Failed to delete workspace")
	}
	return nil
}

// Workspace Users Implementation.
func listWorkspaceUsers(workspaceID int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var users []WorkspaceUser
	if err := ParseResponse(resp, &users); err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users found in this workspace")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Email", "Role", "Time Limit"})

	for _, u := range users {
		table.Append([]string{
			strconv.Itoa(u.ID),
			u.Name,
			u.Email,
			u.Role,
			fmt.Sprintf("%d seconds", u.TimeLimit),
		})
	}

	table.Render()
	return nil
}

func addWorkspaceUser(workspaceID, userID int, role string, timeLimit int) error {
	request := WorkspaceUserCreate{
		UserID: userID,
		Role:   role,
	}

	if timeLimit > 0 {
		request.TimeLimit = timeLimit
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), request)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var user WorkspaceUser
	if err := ParseResponse(resp, &user); err != nil {
		return err
	}

	fmt.Printf("User added to workspace successfully!\n  User ID: %d\n  Name: %s\n  Role: %s\n", user.ID, user.Name, user.Role)
	return nil
}

func updateWorkspaceUser(workspaceID, userID int, role string, timeLimit int) error {
	request := WorkspaceUserUpdate{}

	if role != "" {
		request.Role = role
	}

	if timeLimit > 0 {
		request.TimeLimit = timeLimit
	}

	client := NewClient()
	resp, err := client.DoRequest("PUT", fmt.Sprintf("/api/v1/workspaces/%d/users/%d", workspaceID, userID), request)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var user WorkspaceUser
	if err := ParseResponse(resp, &user); err != nil {
		return err
	}

	fmt.Printf("User updated in workspace successfully!\n  User ID: %d\n  Name: %s\n  Role: %s\n  Time Limit: %d seconds\n", user.ID, user.Name, user.Role, user.TimeLimit)
	return nil
}

func removeWorkspaceUser(workspaceID, userID int) error {
	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/workspaces/%d/users/%d", workspaceID, userID), nil)
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
		fmt.Println("User removed from workspace successfully")
	} else {
		fmt.Println("Failed to remove user from workspace")
	}
	return nil
}

// Workspace Blueprints Implementation.
func listWorkspaceBlueprints(workspaceID int) error {
	client := NewClient()
	resp, err := client.DoRequest("GET", fmt.Sprintf("/api/v1/workspaces/%d/blueprints", workspaceID), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var blueprints []map[string]interface{}
	if err := ParseResponse(resp, &blueprints); err != nil {
		return err
	}

	if len(blueprints) == 0 {
		fmt.Println("No blueprints found shared with this workspace")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Blueprint ID", "Blueprint Type", "Permission", "Name"})

	for _, b := range blueprints {
		blueprintID := fmt.Sprintf("%v", b["blueprint_id"])
		blueprintType := fmt.Sprintf("%v", b["blueprint_type"])
		permission := fmt.Sprintf("%v", b["permission"])
		name := ""
		if n, ok := b["name"]; ok {
			name = fmt.Sprintf("%v", n)
		}

		table.Append([]string{
			blueprintID,
			blueprintType,
			permission,
			name,
		})
	}

	table.Render()
	return nil
}

func addWorkspaceBlueprint(workspaceID, blueprintID int, blueprintType, permission string) error {
	request := WorkspaceBlueprint{
		BlueprintID:   blueprintID,
		BlueprintType: blueprintType,
		Permission:    permission,
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", fmt.Sprintf("/api/v1/workspaces/%d/blueprints", workspaceID), request)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result map[string]interface{}
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Printf("Blueprint shared with workspace successfully!\n  Blueprint ID: %d\n  Type: %s\n  Permission: %s\n", blueprintID, blueprintType, permission)
	return nil
}

func removeWorkspaceBlueprint(workspaceID, blueprintID int, blueprintType string) error {
	request := map[string]interface{}{
		"blueprint_id":   blueprintID,
		"blueprint_type": blueprintType,
	}

	client := NewClient()
	resp, err := client.DoRequest("DELETE", fmt.Sprintf("/api/v1/workspaces/%d/blueprints/%d", workspaceID, blueprintID), request)
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
		fmt.Println("Blueprint removed from workspace successfully")
	} else {
		fmt.Println("Failed to remove blueprint from workspace")
	}
	return nil
}

func init() {
	// Create workspace command flags
	createWorkspaceCmd.Flags().String("name", "", "Name for the workspace")
	createWorkspaceCmd.Flags().String("description", "", "Optional description for the workspace")
	createWorkspaceCmd.Flags().Int("time-limit", 0, "Default time limit for users in the workspace (in seconds)")

	// Workspace user command flags
	addWorkspaceUserCmd.Flags().Int("user-id", 0, "ID of the user to add")
	addWorkspaceUserCmd.Flags().String("role", "", "Role for the user (owner, manager, or member)")
	addWorkspaceUserCmd.Flags().Int("time-limit", 0, "Time limit for the user in the workspace (in seconds)")

	updateWorkspaceUserCmd.Flags().String("role", "", "New role for the user (owner, manager, or member)")
	updateWorkspaceUserCmd.Flags().Int("time-limit", 0, "New time limit for the user in the workspace (in seconds)")

	// Workspace blueprint command flags
	addWorkspaceBlueprintCmd.Flags().Int("blueprint-id", 0, "ID of the blueprint to share")
	addWorkspaceBlueprintCmd.Flags().String("blueprint-type", "", "Type of the blueprint (range, vpc, subnet, or host)")
	addWorkspaceBlueprintCmd.Flags().String("permission", "", "Permission level (view, deploy, or edit)")

	removeWorkspaceBlueprintCmd.Flags().String("blueprint-type", "", "Type of the blueprint (range, vpc, subnet, or host)")

	// Add workspace subcommands
	workspacesCmd.AddCommand(listWorkspacesCmd)
	workspacesCmd.AddCommand(getWorkspaceCmd)
	workspacesCmd.AddCommand(createWorkspaceCmd)
	workspacesCmd.AddCommand(deleteWorkspaceCmd)

	// Add workspace user subcommands
	workspacesCmd.AddCommand(listWorkspaceUsersCmd)
	workspacesCmd.AddCommand(addWorkspaceUserCmd)
	workspacesCmd.AddCommand(updateWorkspaceUserCmd)
	workspacesCmd.AddCommand(removeWorkspaceUserCmd)

	// Add workspace blueprint subcommands
	workspacesCmd.AddCommand(listWorkspaceBlueprintsCmd)
	workspacesCmd.AddCommand(addWorkspaceBlueprintCmd)
	workspacesCmd.AddCommand(removeWorkspaceBlueprintCmd)

	// Add workspaces command to root
	rootCmd.AddCommand(workspacesCmd)
}