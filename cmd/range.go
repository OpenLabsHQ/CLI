package cmd

import (
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"

    "github.com/spf13/cobra"
)

type RequestData struct {
	ID string `json:"id"`
}

var rangeCmd = &cobra.Command{
    Use:   "range",
    Short: "Deploy, and manage range",
    Long:  "This command will let you deploy, power on/off, and manage your range.",
}

var deployRangeCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy a range",
    Long:  "This command will deploy a range to the OpenLabs API.",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("No range ID provided")
            return
        }

        err := deployRange(args)
        if err != nil {
            fmt.Println(err)
        }
    },
}

func deployRange(templateIDs []string) error {
    url := "http://localhost:8000/api/v1/ranges/deploy"

    var requestData []RequestData

    for _, templateID := range templateIDs {
        requestData = append(requestData, RequestData{ID: templateID})
    }

    jsonData, err := json.Marshal(requestData)    
    if err != nil {
        return fmt.Errorf("Failed to marshal request body: %s", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("Failed to create request: %s", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Failed to send request: %s", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        fmt.Println("Range deployed successfully")
    } else {
        return fmt.Errorf("Failed to deploy range: %s", resp.Status)
    }
    
    return nil
}

func init() {
    rangeCmd.AddCommand(deployRangeCmd)
    rootCmd.AddCommand(rangeCmd)
}

