package backups

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/backups/model"
	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func CreateBackup(client *http.Client) *model.BackupResponse {
	basePath := os.Getenv("BASE_URL")
	pwd := os.Getenv("PASSWORD")
	apiUrl := shared.ConcatBaseUrlAndUri(basePath, "/api/teleporter")
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		panic(err)
	}
	authorizationCode := AuthorizeWithPihole(basePath, pwd, client)
	fmt.Printf("Creating backup in %s...\n", basePath)
	req.Header.Add("X-FTL-SID", authorizationCode)
	req.Header.Set("accept", "application/zip")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to fetch: %s", resp.Status))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	DeauthorizeToken(client, authorizationCode, basePath)
	return &model.BackupResponse{Content: content}
}
