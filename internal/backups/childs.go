package backups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups/model"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/deserializers"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

type BackupChildResponse struct {
	Processed []string `json:"processed"`
	Took      float32  `json:"took"`
}

func RestoreBackupInChilds(client *http.Client, backup *model.BackupResponse) *BackupChildResponse {
	basePath := os.Getenv("CHILD_URL")
	pwd := os.Getenv("CHILD_PASSWORD")
	apiUrl := shared.ConcatBaseUrlAndUri(basePath, "/api/teleporter")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "export.zip")
	if err != nil {
		panic(err)
	}
	_, err = part.Write(backup.Content)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	part, err = writer.CreateFormField("import")
	if err != nil {
		panic(err)
	}
	_, err = part.Write(importJson)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		panic(err)
	}
	authorizationCode := AuthorizeWithPihole(basePath, pwd, client)
	fmt.Printf("Trying to restore the backup in %s...\n", basePath)
	req.Header.Add("accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-FTL-SID", authorizationCode)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error trying to sync the adlists... %s \n", content)
		panic(content)
	}

	resBackupChild := &BackupChildResponse{}
	if err := deserializers.JsonDeserialize(resp.Body, resBackupChild); err != nil {
		panic(err)
	}
	DeauthorizeToken(client, authorizationCode, basePath)
	return resBackupChild
}
