package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Configuration represents the CLI configuration.
type Configuration struct {
	APIURL    string `json:"api_url"`
	AuthToken string `json:"auth_token"`
	EncKey    string `json:"enc_key"`
}

// Config Commands.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  "This command lets you view and update CLI configuration settings.",
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current configuration",
	Long:  "This command displays the current CLI configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		fmt.Printf("API URL: %s\n", config.APIURL)
		fmt.Printf("Auth Token: %s\n", config.AuthToken)
		fmt.Printf("Encryption Key: %s\n", config.EncKey)
	},
}

var configSetAPIURLCmd = &cobra.Command{
	Use:   "set-api-url [url]",
	Short: "Set the API URL",
	Long:  "This command sets the API URL for the CLI.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setAPIURL(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var configSetTokenCmd = &cobra.Command{
	Use:   "set-token [token]",
	Short: "Set the auth token",
	Long:  "This command sets the authentication token for the CLI.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setAuthToken(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

var configSetEncKeyCmd = &cobra.Command{
	Use:   "set-enckey [key]",
	Short: "Set the encryption key",
	Long:  "This command sets the encryption key for the CLI.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setEncryptionKey(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Configuration Implementation.
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".openlabs")

	// Create directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			return "", err
		}
	}

	return configDir, nil
}

func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

func loadConfig() (Configuration, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Configuration{}, err
	}

	config := Configuration{
		APIURL: "http://localhost:8000",
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		if err := saveConfig(config); err != nil {
			return Configuration{}, err
		}
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Configuration{}, err
	}

	// Parse config
	if err := json.Unmarshal(data, &config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}

func saveConfig(config Configuration) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Serialize config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write config file
	return os.WriteFile(configPath, data, 0600)
}

func setAPIURL(url string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.APIURL = url

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Println("API URL updated successfully")
	return nil
}

func setAuthToken(token string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.AuthToken = token
	AuthToken = token

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Println("Auth token updated successfully")
	return nil
}

func setEncryptionKey(key string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.EncKey = key
	EncKey = key

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Println("Encryption key updated successfully")
	return nil
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetAPIURLCmd)
	configCmd.AddCommand(configSetTokenCmd)
	configCmd.AddCommand(configSetEncKeyCmd)

	rootCmd.AddCommand(configCmd)

	// Load config at startup
	config, err := loadConfig()
	if err == nil {
		APIURL = config.APIURL
		AuthToken = config.AuthToken
		EncKey = config.EncKey
	}
}
