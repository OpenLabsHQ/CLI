package cmd

import (
    "fmt"
    "io"
    "os"
    "net/http"
    "encoding/json"

    "github.com/spf13/cobra"
)


/*
/templates/ranges/{range_id}
/templates/ranges
/templates/vpcs/{range_id}
/templates/vpcs
/templates/subnets/{subnet_id}
/templates/subnets
/templates/hosts/{host_id}
/templates/hosts
*/

var templatesCmd = &cobra.Command{
    Use:   "templates",
    Short: "Upload and manage templates",
    Long:  "This command will let you upload, view, and delete templates for ranges, subnets, and hosts.",
}

var uploadTemplateCmd = &cobra.Command{
    Use:   "upload",
    Short: "Upload a range template",
    Long:  "This command will upload a range template to the OpenLabs API.",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("No file provided")
            return
        }

        err := uploadTemplate(args[0])
        if err != nil {
            fmt.Println(err)
        }
    },
}

var getTemplateCmd = &cobra.Command{
    Use:   "get",
    Short: "Get a range template",
    Long:  "This command will get a range template from the OpenLabs API.",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Println("No range ID provided")
            return
        }

        err := getTemplate(args[0])
        if err != nil {
            fmt.Println(err)
        }
    },
}


func getTemplate(rangeID string) error {
    url := fmt.Sprintf("http://localhost:8000/api/v1/templates/ranges/%s", rangeID)

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to get template: %s", resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("Failed to read response body: %s", err)
    }

    var templateData map[string]interface{}
    if err := json.Unmarshal(body, &templateData); err != nil {
        return fmt.Errorf("Failed to unmarshal response: %s", err)
    }

    prettyJSON, err := json.MarshalIndent(templateData, "", "  ")
    if err != nil {
        return fmt.Errorf("Failed to marshal response: %s", err)
    }

    fmt.Println(string(prettyJSON))

    return nil
}

func uploadTemplate(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    req, err := http.NewRequest("POST", "http://localhost:8000/api/v1/templates/ranges", file)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    var templateUploadResponseData struct {
        ID string `json:"id"`
    }

    if err := json.Unmarshal(body, &templateUploadResponseData); err != nil {
        return fmt.Errorf("Failed to unmarshal response: %s", err)
    }

    fmt.Printf("Template uploaded successfully!\n  ID: %s\n", templateUploadResponseData.ID)


    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to upload template: %s", resp.Status)
    }
    
    return nil
}

func init() {
    templatesCmd.AddCommand(uploadTemplateCmd)
    templatesCmd.AddCommand(getTemplateCmd)
    rootCmd.AddCommand(templatesCmd)
}
