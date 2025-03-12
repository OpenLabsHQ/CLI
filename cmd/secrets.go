package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Structures for secrets commands.
type AWSSecrets struct {
	AWSAccessKey string `json:"aws_access_key"`
	AWSSecretKey string `json:"aws_secret_key"`
}

type AzureSecrets struct {
	ClientID       string `json:"azure_client_id"`
	ClientSecret   string `json:"azure_client_secret"`
	TenantID       string `json:"azure_tenant_id"`
	SubscriptionID string `json:"azure_subscription_id"`
}

type SecretStatus struct {
	HasCredentials bool       `json:"has_credentials"`
	CreatedAt      *time.Time `json:"created_at"`
}

// UserSecrets holds the credential status for cloud providers
type UserSecrets struct {
	AWS   SecretStatus `json:"aws"`
	Azure SecretStatus `json:"azure"`
}

// Secrets Commands.
var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Upload and manage secrets",
	Long:  "This command will upload and manage secrets for the OpenLabs providers.",
}

var getSecretsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of your secrets",
	Long:  "This command will retrieve the status of your cloud provider secrets.",
	Run: func(cmd *cobra.Command, args []string) {
		err := getSecretsStatus()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var updateAWSSecretsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Update AWS secrets",
	Long:  "This command will update your AWS secrets on the OpenLabs API interactively.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if non-interactive mode is requested
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

		// Get credentials from flags if in non-interactive mode
		if nonInteractive {
			accessKey, _ := cmd.Flags().GetString("access-key")
			secretKey, _ := cmd.Flags().GetString("secret-key")

			if accessKey == "" || secretKey == "" {
				fmt.Println("Error: both --access-key and --secret-key are required in non-interactive mode")
				return
			}

			err := updateAWSSecrets(accessKey, secretKey)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Interactive mode
		accessKey, secretKey, err := promptAWSCredentials()
		if err != nil {
			fmt.Printf("Error getting AWS credentials: %s\n", err)
			return
		}

		err = updateAWSSecrets(accessKey, secretKey)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var updateAzureSecretsCmd = &cobra.Command{
	Use:   "azure",
	Short: "Update Azure secrets",
	Long:  "This command will update your Azure secrets on the OpenLabs API interactively.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if non-interactive mode is requested
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

		// Get credentials from flags if in non-interactive mode
		if nonInteractive {
			clientID, _ := cmd.Flags().GetString("client-id")
			clientSecret, _ := cmd.Flags().GetString("client-secret")
			tenantID, _ := cmd.Flags().GetString("tenant-id")
			subscriptionID, _ := cmd.Flags().GetString("subscription-id")

			if clientID == "" || clientSecret == "" || tenantID == "" || subscriptionID == "" {
				fmt.Println("Error: all Azure credential parameters are required in non-interactive mode")
				return
			}

			err := updateAzureSecrets(clientID, clientSecret, tenantID, subscriptionID)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Interactive mode
		clientID, clientSecret, tenantID, subscriptionID, err := promptAzureCredentials()
		if err != nil {
			fmt.Printf("Error getting Azure credentials: %s\n", err)
			return
		}

		err = updateAzureSecrets(clientID, clientSecret, tenantID, subscriptionID)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Prompt functions for interactive credential input.
func promptAWSCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🔑 Enter your AWS credentials:")

	// Get AWS Access Key
	fmt.Print("\nAWS Access Key: ")
	accessKey, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	accessKey = strings.TrimSpace(accessKey)

	// Get AWS Secret Key (masked input)
	fmt.Print("AWS Secret Key: ")
	secretKeyBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println() // Add newline after hidden input

	secretKey := string(secretKeyBytes)
	secretKey = strings.TrimSpace(secretKey)

	if accessKey == "" || secretKey == "" {
		return "", "", fmt.Errorf("both AWS Access Key and Secret Key are required")
	}

	return accessKey, secretKey, nil
}

func promptAzureCredentials() (string, string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🔑 Enter your Azure credentials:")

	// Get Azure Client ID
	fmt.Print("\nAzure Client ID: ")
	clientID, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", "", err
	}
	clientID = strings.TrimSpace(clientID)

	// Get Azure Client Secret (masked input)
	fmt.Print("Azure Client Secret: ")
	clientSecretBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", "", err
	}
	fmt.Println() // Add newline after hidden input

	clientSecret := string(clientSecretBytes)
	clientSecret = strings.TrimSpace(clientSecret)

	// Get Azure Tenant ID
	fmt.Print("Azure Tenant ID: ")
	tenantID, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", "", err
	}
	tenantID = strings.TrimSpace(tenantID)

	// Get Azure Subscription ID
	fmt.Print("Azure Subscription ID: ")
	subscriptionID, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", "", err
	}
	subscriptionID = strings.TrimSpace(subscriptionID)

	if clientID == "" || clientSecret == "" || tenantID == "" || subscriptionID == "" {
		return "", "", "", "", fmt.Errorf("all Azure credential fields are required")
	}

	return clientID, clientSecret, tenantID, subscriptionID, nil
}

// Secrets Implementation.
func getSecretsStatus() error {
	fmt.Println("\n🔍 Fetching cloud provider credentials status...")

	client := NewClient()
	resp, err := client.DoRequest("GET", "/api/v1/users/me/secrets", nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var secrets UserSecrets
	if err := ParseResponse(resp, &secrets); err != nil {
		return err
	}

	fmt.Println("\n✅ Cloud provider credentials status retrieved successfully!")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Provider", "Status", "Created At"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	// Format the created date
	awsCreatedAt := "N/A"
	if secrets.AWS.CreatedAt != nil {
		awsCreatedAt = secrets.AWS.CreatedAt.Format("2006-01-02 15:04:05")
	}

	azureCreatedAt := "N/A"
	if secrets.Azure.CreatedAt != nil {
		azureCreatedAt = secrets.Azure.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// Format status with icons
	awsStatus := "❌ Not configured"
	if secrets.AWS.HasCredentials {
		awsStatus = "✅ Configured"
	}

	azureStatus := "❌ Not configured"
	if secrets.Azure.HasCredentials {
		azureStatus = "✅ Configured"
	}

	table.Append([]string{"AWS", awsStatus, awsCreatedAt})
	table.Append([]string{"Azure", azureStatus, azureCreatedAt})

	table.Render()
	return nil
}

func updateAWSSecrets(accessKey, secretKey string) error {
	fmt.Println("\n🔄 Updating AWS credentials...")

	secrets := AWSSecrets{
		AWSAccessKey: accessKey,
		AWSSecretKey: secretKey,
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/users/me/secrets/aws", secrets)
	if err != nil {
		return fmt.Errorf("failed to update AWS credentials: %s", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	// Check for response message
	var result struct {
		Message string `json:"message"`
	}
	if err := ParseResponse(resp, &result); err != nil {
		// Even if parsing fails, consider it a success if the request was successful
		fmt.Println("\n✅ AWS credentials updated successfully!")
		return nil
	}

	if result.Message != "" {
		fmt.Printf("\n✅ %s\n", result.Message)
	} else {
		fmt.Println("\n✅ AWS credentials updated successfully!")
	}

	return nil
}

func updateAzureSecrets(clientID, clientSecret, tenantID, subscriptionID string) error {
	fmt.Println("\n🔄 Updating Azure credentials...")

	secrets := AzureSecrets{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TenantID:       tenantID,
		SubscriptionID: subscriptionID,
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/users/me/secrets/azure", secrets)
	if err != nil {
		return fmt.Errorf("failed to update Azure credentials: %s", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	// Check for response message
	var result struct {
		Message string `json:"message"`
	}
	if err := ParseResponse(resp, &result); err != nil {
		// Even if parsing fails, consider it a success if the request was successful
		fmt.Println("\n✅ Azure credentials updated successfully!")
		return nil
	}

	if result.Message != "" {
		fmt.Printf("\n✅ %s\n", result.Message)
	} else {
		fmt.Println("\n✅ Azure credentials updated successfully!")
	}

	return nil
}

func init() {
	// AWS secrets flags
	updateAWSSecretsCmd.Flags().String("access-key", "", "AWS access key (for non-interactive mode)")
	updateAWSSecretsCmd.Flags().String("secret-key", "", "AWS secret key (for non-interactive mode)")
	updateAWSSecretsCmd.Flags().Bool("non-interactive", false, "Use non-interactive mode with command line arguments")

	// Azure secrets flags
	updateAzureSecretsCmd.Flags().String("client-id", "", "Azure client ID (for non-interactive mode)")
	updateAzureSecretsCmd.Flags().String("client-secret", "", "Azure client secret (for non-interactive mode)")
	updateAzureSecretsCmd.Flags().String("tenant-id", "", "Azure tenant ID (for non-interactive mode)")
	updateAzureSecretsCmd.Flags().String("subscription-id", "", "Azure subscription ID (for non-interactive mode)")
	updateAzureSecretsCmd.Flags().Bool("non-interactive", false, "Use non-interactive mode with command line arguments")

	// Add subcommands to secrets command
	secretsCmd.AddCommand(getSecretsStatusCmd)
	secretsCmd.AddCommand(updateAWSSecretsCmd)
	secretsCmd.AddCommand(updateAzureSecretsCmd)

	// Add secrets command to root
	rootCmd.AddCommand(secretsCmd)
}
