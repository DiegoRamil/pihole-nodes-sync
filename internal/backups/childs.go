package backups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups/model"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/deserializers"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

type BackupChildResponse struct {
	Processed []string `json:"processed"`
	Took      float32  `json:"took"`
}

type ChildNode struct {
	URL      string
	Password string
}

// GetChildNodes retrieves the list of child nodes from environment variables
func GetChildNodes() []ChildNode {
	childUrls := strings.Split(shared.RetrieveEnvVar("CHILD_URLS"), ",")
	childPasswords := strings.Split(shared.RetrieveEnvVar("CHILD_PASSWORDS"), ",")

	// Ensure we have the same number of URLs and passwords
	if len(childUrls) != len(childPasswords) {
		fmt.Printf("Warning: Number of CHILD_URLS (%d) doesn't match CHILD_PASSWORDS (%d)\n", len(childUrls), len(childPasswords))
		// Use the shorter length to avoid index out of range
		minLen := min(len(childPasswords), len(childUrls))
		childUrls = childUrls[:minLen]
		childPasswords = childPasswords[:minLen]
	}

	var nodes []ChildNode
	for i := range childUrls {
		url := strings.TrimSpace(childUrls[i])
		password := strings.TrimSpace(childPasswords[i])
		if url != "" && password != "" {
			nodes = append(nodes, ChildNode{
				URL:      url,
				Password: password,
			})
		}
	}
	return nodes
}

func RestoreBackupInChild(client *http.Client, backup *model.BackupResponse, node ChildNode) *BackupChildResponse {
	if node.URL == "" {
		return nil
	}

	apiUrl := shared.ConcatBaseUrlAndUri(node.URL, "/api/teleporter")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "export.zip")
	if err != nil {
		fmt.Printf("Error creating form file for %s: %s\n", node.URL, err)
		return nil
	}
	_, err = part.Write(backup.Content)
	if err != nil {
		fmt.Printf("Error writing the file for %s: %s\n", node.URL, err)
		return nil
	}

	importData := map[string]any{
		"config":      false,
		"dhcp_leases": false,
		"gravity": map[string]any{
			"group":               false,
			"adlist":              true,
			"adlist_by_group":     true,
			"domainlist":          false,
			"domainlist_by_group": false,
			"client":              false,
			"client_by_group":     false,
		},
	}

	importJson, err := json.Marshal(importData)
	if err != nil {
		fmt.Printf("Error marshalling the import data for %s: %s\n", node.URL, err)
		return nil
	}

	part, err = writer.CreateFormField("import")
	if err != nil {
		fmt.Printf("Error creating form field for %s: %s\n", node.URL, err)
		return nil
	}
	_, err = part.Write(importJson)
	if err != nil {
		fmt.Printf("Error writing the import data for %s: %s\n", node.URL, err)
		return nil
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing the writer for %s: %s\n", node.URL, err)
		return nil
	}

	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		fmt.Printf("Error creating the request for %s: %s\n", node.URL, err)
		return nil
	}

	authorizationCode := AuthorizeWithPihole(node.URL, node.Password, client)
	fmt.Printf("Trying to restore the backup in %s...\n", node.URL)
	req.Header.Add("accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-FTL-SID", authorizationCode)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error trying to restore the backup for %s: %s\n", node.URL, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error trying to sync the adlists for %s: %s\n", node.URL, content)
		return nil
	}

	resBackupChild := &BackupChildResponse{}
	if err := deserializers.JsonDeserialize(resp.Body, resBackupChild); err != nil {
		fmt.Printf("Error deserializing the body for %s: %s\n", node.URL, err)
		return nil
	}

	update_gravity, err := strconv.ParseBool(shared.RetrieveEnvVar("UPDATE_GRAVITY"))
	if err != nil {
		fmt.Printf("Error parsing the update gravity env var for %s: %s\n", node.URL, err)
	}

	if update_gravity {
		UpdateGravity(client, node.URL, authorizationCode)
	}
	DeauthorizeToken(client, authorizationCode, node.URL)

	return &BackupChildResponse{
		Processed: []string{node.URL},
		Took:      resBackupChild.Took,
	}
}

// RestoreBackupInChilds is kept for backward compatibility
func RestoreBackupInChilds(client *http.Client, backup *model.BackupResponse) *BackupChildResponse {
	nodes := GetChildNodes()
	var processed []string
	var totalTook float32

	for _, node := range nodes {
		result := RestoreBackupInChild(client, backup, node)
		if result != nil {
			processed = append(processed, result.Processed...)
			totalTook += result.Took
		}
	}

	return &BackupChildResponse{
		Processed: processed,
		Took:      totalTook,
	}
}
