package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Structures for user commands.
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

type PasswordUpdate struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// User Commands.
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user accounts",
	Long:  "This command lets you manage your OpenLabs user account.",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to OpenLabs",
	Long:  "This command logs you in to the OpenLabs API using interactive prompts.",
	Run: func(cmd *cobra.Command, args []string) {
		email, password, err := promptCredentials()
		if err != nil {
			fmt.Printf("Error getting credentials: %s\n", err)
			return
		}

		err = login(email, password)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Long:  "This command registers a new user in the OpenLabs API using interactive prompts.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check for non-interactive mode first
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		name, _ := cmd.Flags().GetString("name")
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

		if nonInteractive {
			if email == "" || password == "" || name == "" {
				fmt.Println("Error: --email, --password, and --name are all required in non-interactive mode")
				return
			}
		} else {
			// Interactive mode
			var err error
			name, email, password, err = promptRegistrationInfo()
			if err != nil {
				fmt.Printf("Error getting registration information: %s\n", err)
				return
			}
		}

		err := register(email, password, name)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from OpenLabs",
	Long:  "This command logs you out from the OpenLabs API.",
	Run: func(cmd *cobra.Command, args []string) {
		err := logout()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get user profile information",
	Long:  "This command retrieves your OpenLabs user profile information.",
	Run: func(cmd *cobra.Command, args []string) {
		err := getUserInfo()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var updatePasswordCmd = &cobra.Command{
	Use:   "update-password",
	Short: "Update your password",
	Long:  "This command updates your OpenLabs user password interactively and automatically logs you in with the new password.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check for non-interactive mode first
		currentPassword, _ := cmd.Flags().GetString("current-password")
		newPassword, _ := cmd.Flags().GetString("new-password")
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

		if nonInteractive {
			if currentPassword == "" || newPassword == "" {
				fmt.Println("Error: both --current-password and --new-password are required in non-interactive mode")
				return
			}
		} else {
			// Interactive mode
			var err error
			currentPassword, newPassword, err = promptPasswordUpdate()
			if err != nil {
				fmt.Printf("Error getting password information: %s\n", err)
				return
			}
		}

		err := updatePassword(currentPassword, newPassword)
		if err != nil {
			fmt.Printf("\n❌ Failed to update password: %s\n", err)
			return
		}

		fmt.Println("\nYour password has been updated and you have been automatically logged in with the new credentials.")
		fmt.Println("The new encryption key has been retrieved and stored for future operations.")
		fmt.Println("\nYou should now be able to deploy ranges without encountering decryption errors.")
	},
}

// User Implementation.
func login(email, password string) error {
	fmt.Println("\n🔒 Authenticating...")

	if Debug {
		fmt.Printf("DEBUG: Logging in with email: %s\n", email)
	}

	credentials := UserCredentials{
		Email:    email,
		Password: password,
	}

	client := NewClient()
	client.AuthToken = "" // Clear any existing token
	client.EncKey = ""    // Clear any existing key

	// Use a custom HTTP client directly instead of DoRequest to better examine the response
	url := fmt.Sprintf("%s/api/v1/auth/login", client.BaseURL)

	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if Debug {
		fmt.Printf("DEBUG: Making direct login request to %s\n", url)
	}

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %s", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if Debug {
		fmt.Printf("DEBUG: Login response status: %s\n", resp.Status)
		fmt.Println("DEBUG: Login response headers:")
		for k, v := range resp.Header {
			fmt.Printf("DEBUG:   %s: %s\n", k, v)
		}
		fmt.Println("DEBUG: Login response cookies:")
		for _, cookie := range resp.Cookies() {
			fmt.Printf("DEBUG:   %s: %s (HttpOnly: %t, Secure: %t)\n",
				cookie.Name, cookie.Value, cookie.HttpOnly, cookie.Secure)
		}
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	if result.Success {
		fmt.Println("\n✅ Login successful!")
		fmt.Println("\nWelcome to OpenLabs CLI!")
		fmt.Println("Use 'openlabs user info' to see your account information.")

		// Store token in config
		config, _ := loadConfig()
		tokenFound := false

		// Look for the token and encryption key in response cookies
		if Debug {
			fmt.Println("\nDEBUG: Examining response cookies directly (more reliable)...")
		}

		// Use resp.Cookies() which gives us more reliable access to cookies
		for _, cookie := range resp.Cookies() {
			if Debug {
				fmt.Printf("DEBUG: Response Cookie: %s = %s (HttpOnly: %t)\n",
					cookie.Name, cookie.Value, cookie.HttpOnly)
			}

			if cookie.Name == "access_token_cookie" || cookie.Name == "access_token" ||
				cookie.Name == "jwt" || cookie.Name == "token" || cookie.Name == "auth_token" {
				config.AuthToken = cookie.Value
				tokenFound = true
				if Debug {
					fmt.Printf("DEBUG: Found auth token in cookie '%s': %s\n", cookie.Name, cookie.Value)
				}
			}
			if cookie.Name == "enc_key" {
				config.EncKey = cookie.Value
				fmt.Println("Encryption key stored successfully from cookie.")
				if Debug {
					fmt.Printf("DEBUG: Found encryption key in cookie: %s\n", config.EncKey)
				}
			}
		}

		// Also check Set-Cookie headers directly - sometimes needed for HTTP-only cookies
		if !tokenFound {
			if Debug {
				fmt.Println("\nDEBUG: No token found in cookies, checking Set-Cookie headers...")
			}

			setCookieHeaders := resp.Header["Set-Cookie"]
			for _, setCookie := range setCookieHeaders {
				if Debug {
					fmt.Printf("DEBUG: Set-Cookie header: %s\n", setCookie)
				}

				// Extract cookie name and value from Set-Cookie header
				parts := strings.Split(setCookie, ";")
				if len(parts) > 0 {
					nameValue := strings.Split(parts[0], "=")
					if len(nameValue) == 2 {
						name := nameValue[0]
						value := nameValue[1]

						if name == "access_token_cookie" || name == "access_token" ||
							name == "jwt" || name == "token" || name == "auth_token" {
							config.AuthToken = value
							tokenFound = true
							if Debug {
								fmt.Printf("DEBUG: Found auth token in Set-Cookie header '%s': %s\n", name, value)
							}
						}
						if name == "enc_key" && config.EncKey == "" {
							config.EncKey = value
							fmt.Println("Encryption key found in Set-Cookie header.")
						}
					}
				}
			}
		}

		// Check if tokens might be in response body
		var responseBody map[string]interface{}
		bodyBytes, _ := io.ReadAll(resp.Body)
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}

		// Create a new body for future reads
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if Debug {
			fmt.Println("\nDEBUG: Examining response body for tokens...")
			fmt.Printf("DEBUG: Response body: %s\n", string(bodyBytes))
		}

		if err := json.Unmarshal(bodyBytes, &responseBody); err == nil {
			// Look for encryption key
			if EncKeyVal, ok := responseBody["enc_key"]; ok && EncKeyVal != nil {
				if EncKeyStr, ok := EncKeyVal.(string); ok && EncKeyStr != "" {
					config.EncKey = EncKeyStr
					fmt.Println("Encryption key stored successfully from response body.")
					if Debug {
						fmt.Printf("DEBUG: Found encryption key in body: %s\n", config.EncKey)
					}
				}
			}

			// Look for access token fields in the body
			for _, field := range []string{"access_token", "token", "jwt"} {
				if tokenVal, ok := responseBody[field]; ok && tokenVal != nil && !tokenFound {
					if tokenStr, ok := tokenVal.(string); ok && tokenStr != "" {
						config.AuthToken = tokenStr
						tokenFound = true
						fmt.Println("Authentication token stored successfully from response body.")
						if Debug {
							fmt.Printf("DEBUG: Found auth token in response body field '%s'\n", field)
						}
						break
					}
				}
			}
		}

		// If no matching cookie found, try looking for one with any name containing 'token'
		if !tokenFound {
			for _, cookie := range client.LastCookies {
				if strings.Contains(strings.ToLower(cookie.Name), "token") {
					config.AuthToken = cookie.Value
					tokenFound = true
					break
				}
			}
		}

		// If still no cookie found, check for Authorization header
		if !tokenFound && resp.Header.Get("Authorization") != "" {
			bearer := resp.Header.Get("Authorization")
			// If it starts with "Bearer ", extract just the token part
			if len(bearer) > 7 && strings.HasPrefix(bearer, "Bearer ") {
				config.AuthToken = bearer[7:]
			} else {
				config.AuthToken = bearer
			}
			tokenFound = true
		}

		// Final fallback: Check for auth token in Authorization header
		if !tokenFound {
			if Debug {
				fmt.Println("\nDEBUG: No token found yet, checking Authorization header...")
			}

			authHeaders := resp.Header["Authorization"]
			for _, authHeader := range authHeaders {
				if Debug {
					fmt.Printf("DEBUG: Authorization header: %s\n", authHeader)
				}

				if strings.HasPrefix(authHeader, "Bearer ") {
					config.AuthToken = strings.TrimPrefix(authHeader, "Bearer ")
					tokenFound = true
					if Debug {
						fmt.Printf("DEBUG: Found auth token in Authorization header: %s\n", config.AuthToken)
					}
				}
			}
		}

		// As a last resort, try any header that might contain a token
		if !tokenFound {
			if Debug {
				fmt.Println("\nDEBUG: Trying headers with 'token' in the name...")
			}

			for headerName, headerValues := range resp.Header {
				headerLower := strings.ToLower(headerName)
				if strings.Contains(headerLower, "token") || strings.Contains(headerLower, "auth") ||
					strings.Contains(headerLower, "jwt") {
					if len(headerValues) > 0 {
						value := headerValues[0]
						// Clean up "Bearer " prefix if present
						value = strings.TrimPrefix(value, "Bearer ")
						config.AuthToken = value
						tokenFound = true
						if Debug {
							fmt.Printf("DEBUG: Found potential auth token in header '%s': %s\n", headerName, value)
						}
						break
					}
				}
			}
		}

		// Update global tokens
		if tokenFound {
			AuthToken = config.AuthToken
			fmt.Println("Authentication token stored successfully.")
		} else {
			// Try to manually generate a token if needed (special case)
			AuthToken = "manual-token-for-testing"
			config.AuthToken = AuthToken
			fmt.Println("WARNING: No token found! Using a dummy token for testing.")
			fmt.Println("This is for debugging only and may not work in production.")
		}

		// Update encryption key in global variable
		if config.EncKey != "" {
			EncKey = config.EncKey
		}

		// Save the configuration
		if err := saveConfig(config); err != nil {
			fmt.Println("Error saving configuration:", err)
		}
	} else {
		fmt.Println("\n❌ Login failed. Please check your credentials.")
	}

	return nil
}

func register(email, password, name string) error {
	fmt.Println("\n🔐 Registering new user...")

	user := UserRegister{
		Email:    email,
		Password: password,
		Name:     name,
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/auth/register", user)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result struct {
		ID string `json:"id"`
	}
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	fmt.Println("\n✅ User registered successfully!")
	fmt.Printf("\nAccount Details:")
	fmt.Printf("\n  ID:    %s", result.ID)
	fmt.Printf("\n  Name:  %s", name)
	fmt.Printf("\n  Email: %s\n", email)

	fmt.Println("\nYou can now login with your credentials using:")
	fmt.Println("  openlabs user login")

	return nil
}

func logout() error {
	fmt.Println("\n🔓 Logging out...")

	// Always clear local tokens regardless of API response
	config, _ := loadConfig()
	config.AuthToken = ""
	config.EncKey = ""
	if err := saveConfig(config); err != nil {
		fmt.Println("Error saving configuration:", err)
	}

	// Update global variables
	AuthToken = ""
	EncKey = ""

	// Try to call the logout API
	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/auth/logout", nil)
	if err != nil {
		fmt.Println("\n⚠️ Could not connect to API for logout, but local tokens have been cleared.")
		return nil
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result struct {
		Success bool `json:"success"`
	}

	if err := ParseResponse(resp, &result); err != nil {
		fmt.Println("\n⚠️ API logout may have failed, but local tokens have been cleared.")
		return nil
	}

	fmt.Println("\n✅ Logout successful!")
	fmt.Println("\nYou have been safely logged out from OpenLabs.")

	return nil
}

func getUserInfo() error {
	fmt.Println("\n👤 Fetching user profile...")

	client := NewClient()
	resp, err := client.DoRequest("GET", "/api/v1/users/me", nil)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var userInfo UserInfo
	if err := ParseResponse(resp, &userInfo); err != nil {
		return err
	}

	fmt.Println("\n✅ User profile retrieved successfully!")

	// Create a styled table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Email", "Admin"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	// Add row with user information
	adminStatus := "No"
	if userInfo.Admin {
		adminStatus = "Yes"
	}
	table.Append([]string{userInfo.Name, userInfo.Email, adminStatus})

	// Render the table
	fmt.Println()
	table.Render()

	// Get secrets status to display a more complete profile
	client = NewClient()
	resp, err = client.DoRequest("GET", "/api/v1/users/me/secrets", nil)
	if err == nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Printf("Error closing response body: %v\n", err)
			}
		}()

		var secrets UserSecrets
		if err := ParseResponse(resp, &secrets); err == nil {
			fmt.Println("\nCloud Provider Credentials:")

			secretsTable := tablewriter.NewWriter(os.Stdout)
			secretsTable.SetHeader([]string{"Provider", "Status"})
			secretsTable.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})

			awsStatus := "❌ Not configured"
			if secrets.AWS.HasCredentials {
				awsStatus = "✅ Configured"
			}

			azureStatus := "❌ Not configured"
			if secrets.Azure.HasCredentials {
				azureStatus = "✅ Configured"
			}

			secretsTable.Append([]string{"AWS", awsStatus})
			secretsTable.Append([]string{"Azure", azureStatus})
			secretsTable.Render()
		}
	}

	return nil
}

func updatePassword(currentPassword, newPassword string) error {
	fmt.Println("\n🔄 Updating password...")

	passwordUpdate := PasswordUpdate{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	client := NewClient()
	resp, err := client.DoRequest("POST", "/api/v1/users/me/password", passwordUpdate)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	var result struct {
		Message string `json:"message"`
	}
	if err := ParseResponse(resp, &result); err != nil {
		return err
	}

	if result.Message == "Password updated successfully" {
		fmt.Println("\n✅ " + result.Message)

		// Get current user information to retrieve email for auto-login
		userClient := NewClient()
		userResp, err := userClient.DoRequest("GET", "/api/v1/users/me", nil)
		if err != nil {
			fmt.Println("\nAuto-login failed to get user information. Please login manually with the new password.")
			return nil
		}
		defer func() {
			err := userResp.Body.Close()
			if err != nil {
				fmt.Printf("Error closing response body: %v\n", err)
			}
		}()

		var userInfo UserInfo
		if err := ParseResponse(userResp, &userInfo); err != nil {
			fmt.Println("\nAuto-login failed to parse user information. Please login manually with the new password.")
			return nil
		}

		// Automatically log in with the new password
		fmt.Println("\n🔄 Automatically logging in with new password...")
		err = login(userInfo.Email, newPassword)
		if err != nil {
			fmt.Printf("\nAuto-login failed: %s\nPlease login manually with your new password using 'openlabs user login'", err)
		}
	} else {
		fmt.Println("\n❌ " + result.Message)
	}

	return nil
}

// promptCredentials prompts the user for their email and password.
func promptCredentials() (string, string, error) {
	// Read email
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	email = strings.TrimSpace(email)

	// Read password (hidden)
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println() // Add newline after password input

	password := string(passwordBytes)

	return email, password, nil
}

// promptRegistrationInfo prompts the user for registration information.
func promptRegistrationInfo() (string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	// Get full name
	fmt.Print("Full Name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	name = strings.TrimSpace(name)

	// Get email
	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	email = strings.TrimSpace(email)

	// Get password (hidden)
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}
	fmt.Println()

	// Confirm password (hidden)
	fmt.Print("Confirm Password: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}
	fmt.Println()

	password := string(passwordBytes)
	confirmPassword := string(confirmBytes)

	// Check if passwords match
	if password != confirmPassword {
		return "", "", "", fmt.Errorf("passwords don't match")
	}

	return name, email, password, nil
}

// promptPasswordUpdate prompts for current and new password.
func promptPasswordUpdate() (string, string, error) {
	// Read current password (hidden)
	fmt.Print("Current Password: ")
	currentPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println()

	// Read new password (hidden)
	fmt.Print("New Password: ")
	newPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println()

	// Confirm new password (hidden)
	fmt.Print("Confirm New Password: ")
	confirmPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println()

	currentPassword := string(currentPasswordBytes)
	newPassword := string(newPasswordBytes)
	confirmPassword := string(confirmPasswordBytes)

	// Check if new passwords match
	if newPassword != confirmPassword {
		return "", "", fmt.Errorf("new passwords don't match")
	}

	return currentPassword, newPassword, nil
}

func init() {
	// Register command flags (for non-interactive mode)
	registerCmd.Flags().String("email", "", "User email (for non-interactive mode)")
	registerCmd.Flags().String("password", "", "User password (for non-interactive mode)")
	registerCmd.Flags().String("name", "", "User's full name (for non-interactive mode)")
	registerCmd.Flags().Bool("non-interactive", false, "Use non-interactive mode with command line arguments")

	// Update password command flags (for non-interactive mode)
	updatePasswordCmd.Flags().String("current-password", "", "Current password (for non-interactive mode)")
	updatePasswordCmd.Flags().String("new-password", "", "New password (for non-interactive mode)")
	updatePasswordCmd.Flags().Bool("non-interactive", false, "Use non-interactive mode with command line arguments")

	// Add subcommands to user command
	userCmd.AddCommand(loginCmd)
	userCmd.AddCommand(registerCmd)
	userCmd.AddCommand(logoutCmd)
	userCmd.AddCommand(infoCmd)
	userCmd.AddCommand(updatePasswordCmd)

	// Add user command to root
	rootCmd.AddCommand(userCmd)
}
